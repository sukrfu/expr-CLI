package expr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
)

// DelField: 删除 map or slice变量中指定的字段
func DelField(s interface{}, fieldName string) error {
	return delField(reflect.ValueOf(s), fieldName, fieldName)
}

func canNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return true
	}
	return false
}

func delField(v reflect.Value, fieldName, currName string) error {
	log.Debugf("First kind %s", v.Kind())
	if !canNil(v) {
		return fmt.Errorf("Struct must be a pointer")
	}

	if v.IsNil() {
		return fmt.Errorf("Obj %v cannot be nil", v)
	}

	// v 第一次指针转换
	v = reflect.Indirect(v)
	// todo: 去掉interface包装
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct, reflect.Ptr:
		currName, nextFieldName := getCurrAndNextFieldName(currName)
		currFieldName := getFieldName(currName)
		log.Debugf("Curr Name %s and next field Name %s", currFieldName, nextFieldName)

		if v.Kind() == reflect.Struct {
			v = v.FieldByName(currFieldName)
		} else {
			v = v.Elem().FieldByName(currFieldName)
		}
		if !v.IsValid() {
			return fmt.Errorf("No such field: %s in obj", currFieldName)
		}

		vl := v
		if vl.Kind() == reflect.Ptr {
			// vl 第二次指针转换
			vl = vl.Elem()
		}

		if vl.Kind() != reflect.Slice && vl.Kind() != reflect.Map {
			fieldName = nextFieldName
		}

		log.Debugf("Field type %s", v.Kind())
		if v.Kind() == reflect.Ptr {
			err := delField(v, fieldName, nextFieldName)
			if err != nil {
				return err
			}
		} else {
			err := delField(v.Addr(), fieldName, nextFieldName)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		av := v
		currName, nextFieldName := getCurrAndNextFieldName(fieldName)
		log.Debugf("Curr Name %s and next field Name %s", currName, nextFieldName)

		index, err := getFieldSliceIndex(currName)
		if err != nil {
			return err
		}

		var newslice reflect.Value
		log.Debugf("Sclice index %d", index)

		if index != -1 && index != -2 {
			arrayElement := av.Index(index)
			if !arrayElement.IsValid() {
				return fmt.Errorf("Slice index invalid")
			}

			if nextFieldName != "" {
				err = DelField(arrayElement.Addr().Interface(), nextFieldName)
				if err != nil {
					return err
				}
			} else {
				newslice = av
				len := newslice.Len()
				slice1 := newslice.Slice(0, index)
				slice2 := newslice.Slice(index+1, len)
				newslice = reflect.AppendSlice(slice1, slice2)
				v.Set(newslice)

			}
		} else {
			return fmt.Errorf("Slice index invalid")
		}

	case reflect.Map:
		av := v
		currName, nextFieldName := getCurrAndNextFieldName(fieldName)
		log.Debugf("Curr Name [%s] and next field Name [%s].", currName, nextFieldName)

		key, err := getFieldMapKey(currName)
		if err != nil {
			return err
		}

		log.Debugf("Map key %s", key)
		if key != "" {
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

				err = delField(curObj, nextFieldName, nextFieldName)
				if err != nil {
					return err
				}
				av.SetMapIndex(keyValue, obj)
			} else {
				av.SetMapIndex(keyValue, reflect.Value{})
			}

		} else {
			return fmt.Errorf("Map key is none")
		}
	default:
		return fmt.Errorf("Obj must map or slice")

	}
	return nil
}
