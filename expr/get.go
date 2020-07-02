package expr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

// Map: 字段名需要以[$fieldName]格式
func GetField(s interface{}, fieldName string) (v reflect.Value, err error) {
	field, err := getNestField(reflect.ValueOf(s), fieldName, fieldName)
	return field, err
}

// 获取v指向的Value对象
func getElem(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		return v.Elem()
	case reflect.Interface:
		v = v.Elem()
		if v.Kind() == reflect.Ptr {
			return v.Elem()
		}
	}
	return v
}

func nextPos(fieldName string) (pos int, slice bool) {
	i := strings.Index(fieldName, ".")
	j := strings.Index(fieldName, "[")
	if i < j {
		if i <= -1 {
			return j, true
		}
		return i, false
	} else {
		if j <= -1 {
			return i, false
		}
		return j, true
	}
}

func getNestField(s reflect.Value, fullName, fieldName string) (reflect.Value, error) {

	//Make sure types are strcut or ptr
	s = getElem(s)

	// todo: 去掉外层interface包装
	if s.Kind() == reflect.Interface {
		s = s.Elem()
	}
	t := s.Kind()
	fmt.Sprintf("Interface %v  kind %v ,fullName %v,fieldName %v.\n", s, t, fullName, fieldName)

	val := s
	if fieldName == "" {
		return val, nil
	}

	if t == reflect.Struct || t == reflect.Ptr || t == reflect.Slice || t == reflect.Map {
	} else {
		return reflect.Value{}, fmt.Errorf("Struct must be struct interface %v", t)
	}

	if t == reflect.Struct {
		if i, slice := nextPos(fieldName); i > -1 {
			currFieldName := fieldName[0:i]
			fname := getFieldName(currFieldName)
			field := val.FieldByName(fname)
			if !field.IsValid() {
				//We should ignore the error since there might be empty field.
				return reflect.Value{}, fmt.Errorf("No such field: %s in obj", fname)
			}
			nextStart := i + 1
			if slice {
				nextStart = i
			}
			nextFieldName := fieldName[nextStart:len(fieldName)]
			fieldValue := field
			switch fieldValue.Kind() {
			case reflect.Slice:
				index, err := getFieldSliceIndex(currFieldName)
				if err != nil {
					return reflect.Value{}, err
				}

				if index != -1 {
					field = fieldValue.Index(index)
				}

				return getNestField(field, fullName, nextFieldName)

			case reflect.Map:
				key, err := getFieldMapKey(currFieldName)

				if err != nil {
					return reflect.Value{}, err
				}
				if key != "" {

					keyType := field.Type().Key()
					keyValue := reflect.New(keyType).Elem()
					err = setStringValue(keyValue, key)
					if err != nil {
						return reflect.Value{}, err
					}
					field = field.MapIndex(keyValue)
				}
				return getNestField(field, fullName, nextFieldName)

			}
			return getNestField(field, fullName, nextFieldName)

		}
	}

	if !val.IsValid() {
		return reflect.Value{}, fmt.Errorf("Nil pointer: %s in obj", fullName)
	}

	field := val
	if t == reflect.Struct || t == reflect.Ptr {
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		field = field.FieldByName(getFieldName(fieldName))

		if !field.IsValid() {
			return field, fmt.Errorf("No such field: %s in obj", fieldName)
		}
	}

	switch field.Kind() {
	case reflect.Slice:
		index, err := getFieldSliceIndex(fieldName)
		if err != nil {
			return reflect.Value{}, err
		}
		log.Debugf("Field slice index %d", index)

		if index != -1 {
			field = field.Index(index)
		} else {
			return field, nil
		}

		i := strings.Index(fieldName, "]")

		nextFieldName := fieldName[i+1 : len(fieldName)]
		nextFieldName = strings.TrimLeft(nextFieldName, ".")
		return getNestField(field, fullName, nextFieldName)

	case reflect.Map:
		key, err := getFieldMapKey(fieldName)
		log.Debugf("Field map key %s", key)

		if err != nil {
			return reflect.Value{}, err
		}
		if key != "" {
			keyType := field.Type().Key()
			keyValue := reflect.New(keyType).Elem()
			err = setStringValue(keyValue, key)
			if err != nil {
				return reflect.Value{}, err
			}

			field = field.MapIndex(keyValue)

			i := strings.Index(fieldName, "]")
			nextFieldName := fieldName[i+1 : len(fieldName)]
			nextFieldName = strings.TrimLeft(nextFieldName, ".")
			return getNestField(field, fullName, nextFieldName)

		}
	case reflect.Struct:
		//TODO
	}

	return field, nil

}

// getFieldName 以'['为分隔符获取有效字段名，'['之后为无效字段
func getFieldName(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		return fieldName[0:strings.Index(fieldName, "[")]
	}

	return fieldName
}
// [2]
func getFieldSliceIndex(fieldName string) (int, error) {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		i, err := strconv.Atoi(index)
		if err != nil {
			return -2, nil
		}
		return i, nil
	}

	return -1, nil
}


func getFieldMapKey(fieldName string) (string, error) {
	// todo: 需要处理不含有']'的输入
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("field %s not found, err: %s", fieldName, err)
		}
	}()
	if strings.Index(fieldName, "[") >= 0 {
		key := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		return key, nil
	}

	return "", nil
}

/*
func getReflectValue(in interface{}) reflect.Value {
	var value reflect.Value
	if reflect.TypeOf(in).Kind() == reflect.Ptr {
		value = reflect.ValueOf(in).Elem()
	} else {
		value = reflect.ValueOf(in)
	}
	return value
}
*/
