package expr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"

	"github.com/apaxa-go/eval"
	"go.uber.org/zap"
)


func Set(s interface{}, fieldName string, value string, args eval.Args, fun func(reflect.Value) reflect.Value) error {
	log.Debugf("Set fieldName: %v, value: %v.", fieldName, value)
	expr, err := eval.ParseString(value, "")
	if err != nil {
		log.Error("Parse string error", zap.Any("err", err))
		return err
	}
	r, err := expr.EvalToInterface(args)
	fmt.Println(r)
	if err != nil {
		log.Error("Parse string error", zap.Any("err", err))
		return err
	}

	return SetField(s, fieldName, r, fun)
}

// deep获取结构体字段
func Get(s interface{}, fieldName string, value string) (string, error) {
	log.Debugf("Get fieldName: %v.", fieldName)
	deep := int16(1)
	if value != "" {
		num, err := strconv.Atoi(value)
		if err == nil && num <= 4 {
			deep = int16(num)
		}
	}
	if fieldName == "." {
		str := FormatObj(s, deep)
		return str, nil
	}
	obj, err := GetField(s, fieldName)
	if err != nil {
		return "", err
	}
	str := FormatValue(obj, deep)
	return str, nil
}

// todo: 删除slice元素
func Del(s interface{}, fieldName string) error {
	return DelField(s, fieldName)
}
