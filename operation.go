package main

import (
	"gitee.com/legou-lib/expr"
	"reflect"
)

type SimpleData interface {
	Get(fieldName string) (string, error)
	Set(fieldName ,valueEval string) error
	Delete(key string) error
	Print()
}

type SimpleDataOperation struct {
	GET string
	SET string
	DEL string
	PRT string
	USE string
	LIST string
}

func NewSimpleDataOperation()SimpleDataOperation{
	return SimpleDataOperation{
		"get",
		"set",
		"delete",
		"print",
		"use",
		"list",
	}
}

func (operation SimpleDataOperation)Get(data SimpleData, fieldName string) (string, error) {
	return data.Get(fieldName)
}

func (operation SimpleDataOperation)Set(data SimpleData, fieldName, valueEval string) error {
	return data.Set(fieldName, valueEval)
}

func (operation SimpleDataOperation)Delete(data SimpleData, fieldName string) error {
	return data.Delete(fieldName)
}

func (operation SimpleDataOperation)Print(data SimpleData){
	data.Print()
}

func (operation SimpleDataOperation) GetObjectRuntimeType(data SimpleData) reflect.Kind{
	return reflect.Indirect(reflect.ValueOf(data)).Kind()
}

// GetObjectFieldType 如果fieldName为空，则获取data的类型 (用于显示字段类型提示)
func (operation SimpleDataOperation)GetObjectFieldType(data SimpleData, fieldName string) string {
	if fieldName == "" {
		return reflect.TypeOf(data).String()
	}
	field, err := expr.GetField(data, fieldName)
	if err != nil {
		return ""
	}
	field = reflect.Indirect(field)
	return field.Type().String()
}

// GetField 获取data中的字段, fieldName 不存在或出错则返回nil
func (operation SimpleDataOperation) GetField(data SimpleData, fieldName string) interface{} {
	fieldValue, err := expr.GetField(data, fieldName)
	if err != nil {
		return nil
	}
	return fieldValue.Interface()
}