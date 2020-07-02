package expr

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	//"legou-lib/logger"
	"reflect"
	"strconv"
	"strings"
)

var (
	//log = logger.MustGetLogger("expr")
)

// 仅支持string类型的设置值
func SetField(s interface{}, fieldName string, value interface{}, fun func(data reflect.Value) reflect.Value) error {
	return setField(reflect.ValueOf(s), fieldName, fieldName, value, fun)
}

func setField(v reflect.Value, fieldName, currName string, value interface{},
	fun func(data reflect.Value) reflect.Value) error {
	log.Debugf("First kind %s", v.Kind())
	log.Debugf("SetField FieldName %s ,currName %s", fieldName, currName)
	if !canNil(v) {
		return fmt.Errorf("Struct must be a pointer")
	}

	if v.IsNil() {
		log.Debugf("Field %s is nil and set a new ", fieldName)
		v.Set(reflect.New(v.Type().Elem()))
	}

	v = reflect.Indirect(v)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	if fun != nil {
		if v.Kind() == reflect.Interface {
			v = fun(v)
		}
	}
	if fieldName == "" {
		switch inputv := value.(type) {
		case string:
			err := setStringValue(v, inputv)
			if err != nil {
				return err
			}
		case []byte:
			err := setStringValue(v, string(inputv))
			if err != nil {
				return err
			}
		default:
			valv := reflect.ValueOf(value)
			for valv.Kind() == reflect.Ptr {
				valv = valv.Elem()
			}
			if v.Type() != valv.Type() {
				return fmt.Errorf("Provided fieldName(%v) value type (%v) didn't match obj field type (%v)\n", fieldName, valv.Type(), v.Type())
			}
			v.Set(valv)
		}

		return nil
	}

	switch v.Kind() {
	case reflect.Struct, reflect.Ptr:
		currName, nextFieldName := getCurrAndNextFieldName(currName)
		currFieldName := getFieldName(currName)
		log.Debugf("Struct FieldName %s Curr Name %s and next field Name %s", currName, currFieldName, nextFieldName)

		if v.Kind() == reflect.Struct {
			v = v.FieldByName(currFieldName)
		} else {
			v = v.Elem().FieldByName(currFieldName)
		}
		if !v.IsValid() {
			return fmt.Errorf("No such field: %s in obj", currFieldName)
		}
		log.Debugf("Field type %s", v.Kind())
		vl := v
		if vl.Kind() == reflect.Ptr {
			vl = vl.Elem()
		}
		if vl.Kind() != reflect.Slice && vl.Kind() != reflect.Map {
			fieldName = nextFieldName
		}

		if v.Kind() == reflect.Ptr {
			err := setField(v, fieldName, nextFieldName, value, fun)
			if err != nil {
				return err
			}
		} else {
			err := setField(v.Addr(), fieldName, nextFieldName, value, fun)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		av := v
		currName, nextFieldName := getCurrAndNextFieldName(fieldName)
		log.Debugf("Slice FieldName %s Curr Name %s and next field Name %s", fieldName, currName, nextFieldName)

		index, err := getFieldSliceIndex(currName)
		if err != nil {
			return err
		}

		elementType := v.Type().Elem()
		var newslice reflect.Value
		log.Debugf("Slice index %d", index)
		if index != -1 && index != -2 {
			//Set to specific
			arrayElement := av.Index(index)
			if !arrayElement.IsValid() {
				arrayElement.Set(reflect.New(elementType).Elem())
			}
			log.Debugf("Slice value %v  set %v", arrayElement, arrayElement.CanSet())
			if nextFieldName != "" {
				err = SetField(arrayElement.Addr().Interface(), nextFieldName, value, fun)
				if err != nil {
					return err
				}
			} else {
				arrayElement.Set(reflect.ValueOf(value))
			}

		} else if index == -2 {
			valueOf := reflect.ValueOf(value)
			if v.Type() != valueOf.Type() {
				return fmt.Errorf("Slice index must be int")
			}
			v.Set(valueOf)
		} else {
			arrayyElement := reflect.New(elementType).Elem()
			arrayyElement.Set(reflect.ValueOf(value))
			newslice = av
			newslice = reflect.Append(newslice, arrayyElement)
			v.Set(newslice)

		}
	case reflect.Map:
		av := v
		currName, nextFieldName := getCurrAndNextFieldName(fieldName)
		log.Debugf("Map FieldName [%s] Curr Name [%s] and next field Name [%s].", fieldName, currName, nextFieldName)

		key, err := getFieldMapKey(currName)
		if err != nil {
			return err
		}

		log.Debugf("Map key %s value %v", key, value)
		if key != "" {
			//Set to specific
			// 获取mapkey类型
			keyType := av.Type().Key()
			keyValue := reflect.New(keyType).Elem()
			err = setStringValue(keyValue, key)
			if err != nil {
				return err
			}

			if nextFieldName != "" {
				obj := av.MapIndex(keyValue)
				curObj := obj
				if obj.Kind() == reflect.Ptr {
					curObj = obj.Elem()
				}
				if !curObj.CanSet() {
					return fmt.Errorf("cannot set")
				}

				if curObj.Kind() != reflect.Ptr {
					curObj = curObj.Addr()
				}

				err = setField(curObj, nextFieldName, nextFieldName, value, fun)
				if err != nil {
					return err
				}
				av.SetMapIndex(keyValue, obj)
			} else {
				av.SetMapIndex(keyValue, reflect.ValueOf(value))
			}

		} else if key == "" {
			valueOf := reflect.ValueOf(value)
			if v.Type() != valueOf.Type() {
				return fmt.Errorf("Provided key(%v) value type (%v) didn't match obj field type (%v)\n", key, valueOf.Type(), v.Type())
			}
			v.Set(valueOf)
		}
	default:
		valueOf := reflect.ValueOf(value)
		if v.Type() != valueOf.Type() {
			return fmt.Errorf("Provided default value type (%v) didn't match obj field type (%v)\n", valueOf.Type(), v.Type())
		}
		v.Set(valueOf)

	}
	return nil
}

// getCurrAndNextFieldName 以'.'为分隔符获取当前字段名和下一个字段名
func getCurrAndNextFieldName(name string) (string, string) {
	currName := name
	nextFieldName := ""
	if i := strings.Index(name, "."); i > -1 {
		currName = name[0:i]
		nextFieldName = name[i+1 : len(name)]
	}
	return currName, nextFieldName
}

func setStringValue(v reflect.Value, value string) (err error) {
	s := value
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
		v.SetBytes([]byte(s))
		return
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return
		}
		if v.OverflowInt(n) {
			err = fmt.Errorf("overflow int64 for %d.", n)
			return
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return
		}
		if v.OverflowUint(n) {
			err = fmt.Errorf("overflow uint64 for %d.", n)
			return
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return
		}
		if v.OverflowFloat(n) {
			err = fmt.Errorf("overflow float64 for %v.", n)
			return
		}
		v.SetFloat(n)
	case reflect.Bool:
		var n bool
		n, err = strconv.ParseBool(s)
		if err != nil {
			return
		}
		v.SetBool(n)
	default:
		err = fmt.Errorf("value %+v can only been set to primary type but was %+v", value, v)
	}

	return
}
