package expr

import (
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
)

var (
	_ATOMIC_VALUE_TYPE = reflect.TypeOf(atomic.Value{})
)

func formatMapKey(values []reflect.Value) string {
	report := ""
	v := values
	if len(values) > 64 {
		v = values[:64]
	}

	for _, v := range v {
		if v.CanInterface() {
			report += fmt.Sprintf("%v, ", v.Interface())
		} else if v.Kind() == reflect.Ptr {
			e := v.Elem()
			if e.CanInterface() {
				report += fmt.Sprintf("%v, ", e.Interface())
			} else {
				report += fmt.Sprintf("NO SUPPORT, ")
			}
		} else {
			report += fmt.Sprintf("%v, ", v)
		}
	}

	length := len(values)
	if length > 64 {
		report += fmt.Sprintf("...len %v", length)
	}

	return report
}

func formatStructByDeep(s reflect.Value, deep int16, maxDeep int16) string {
	var report string
	if deep > maxDeep {
		return report
	}
	if s.Kind() == reflect.Interface {
		s = s.Elem()
	}
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}

	prefix := ""
	for strdeep := deep; strdeep >= 0; strdeep-- {
		prefix += "\t"
	}

	if !s.IsValid() {
		return fmt.Sprintf("%s nil\n", prefix)
	}
	typeOfT := s.Type()
	if typeOfT == _ATOMIC_VALUE_TYPE {
		return fmt.Sprintf("%s atomicValue\n", prefix)
	}
	if s.Kind() == reflect.Struct {
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			if !f.IsValid() {
				continue
			}
			if f.Type() == _ATOMIC_VALUE_TYPE {
				report += fmt.Sprintf("%s%s atomicValue\n", prefix,
					typeOfT.Field(i).Name)
			}
			if f.Kind() == reflect.Map {
				report += fmt.Sprintf("%s%s keys: {%v}\n", prefix,
					typeOfT.Field(i).Name, formatMapKey(f.MapKeys()))
			} else if (f.Kind() == reflect.Slice) || (f.Kind() == reflect.Array) {
				report += fmt.Sprintf("%s%s len: %d\n", prefix,
					typeOfT.Field(i).Name, f.Len())
			} else if f.Kind() == reflect.Struct {
				report += formatSprintfValue(f, maxDeep, prefix, typeOfT, i, deep)
			} else if f.Kind() == reflect.Interface {
				report += formatSprintfValue(f, maxDeep, prefix, typeOfT, i, deep)
			} else if f.CanInterface() {
				report += formatSprintfValue(f, maxDeep, prefix, typeOfT, i, deep)
			} else if f.Kind() == reflect.Ptr {
				if f.IsNil() {
					report += fmt.Sprintf("%s%s=nil\n", prefix,
						typeOfT.Field(i).Name)
				} else {
					e := f.Elem()
					report += formatSprintfValue(e, maxDeep, prefix, typeOfT, i, deep)
				}
			} else {
				report += fmt.Sprintf("%s%s=%v\n", prefix,
					typeOfT.Field(i).Name, f)
			}
		}
	} else {
		if s.Kind() == reflect.Map || s.Kind() == reflect.Array || s.Kind() == reflect.Slice || s.Kind() == reflect.Chan {
			report += fmt.Sprintf("%s%s=%v len %v\n", prefix,
				typeOfT.Name(), s, s.Len())
		} else {
			report += fmt.Sprintf("%s%s=%v\n", prefix,
				typeOfT.Name(), s)
		}
	}
	return report
}

func isBaseType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		fallthrough
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		fallthrough
	case reflect.String:
		return true
	}
	return false
}

func formatSprintfValue(value reflect.Value, maxDeep int16, prefix string, typeOfT reflect.Type, i int, deep int16) string {
	if isBaseType(value.Kind()) {
		return fmt.Sprintf("%s%s=%v\n", prefix, typeOfT.Field(i).Name, value)
	}
	if deep > maxDeep {
		return fmt.Sprintf("%s%s=...\n", prefix,
			typeOfT.Field(i).Name)
	} else {
		report := fmt.Sprintf("%s%s:\n", prefix, typeOfT.Field(i).Name)
		report += formatStructByDeep(value, deep+1, maxDeep)
		return report
	}
}

func formatStruct(obj reflect.Value, maxDeep int16) string {
	return formatStructByDeep(obj, 0, maxDeep)
}

func GetSplitField(src string, delim string, index int) string {
	var ret string
	fields := strings.Split(src, delim)
	if index >= 0 && index < len(fields) {
		ret = fields[index]
	}

	return ret
}

func FormatObj(value interface{}, deep int16) string {
	v := reflect.ValueOf(value)
	return FormatValue(v, deep)
}

func FormatValue(v reflect.Value, deep int16) string {
	str := fmt.Sprintf("Get value: %v .", v)
	str = formatStruct(v, deep)
	return str
}
