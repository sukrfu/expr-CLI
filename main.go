package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"github.com/c-bata/go-prompt"
	"gitee.com/legou-lib/expr"
	log "github.com/sirupsen/logrus"
)

type Obj struct {
	Name   *string
	Man *Person
	Id    *int
}
type Person struct {
	Name *string
	Friend *Person
	Age *int
	Phone *string
}

var myName = "wuwj"
var nickName = "nick"
var tomName = "tom"
var tomAge = 19
var nickAge = 25
var myAge = 22
var nickPhone = "123456"
var tomPhone = "666666"
var obj = StructType{
	Name: &myName,
	Id: &myAge,
	Man: &Person{
		Name:  &nickName,
		Age:   &nickAge,
		Friend: &Person{
			Name:   &tomName,
			Age:    &tomAge,
			Phone:  &tomPhone,
		},
		Phone: &nickPhone,
	},
}

type MapType map[string]string

type SliceType []string

type StructType Obj

var mapObj = MapType{
	"hello": "123",
	"test": "456",
}

var sliceObj = SliceType{"1","2","3","4","5"}

var globalTestMap = map[string]interface{}{
	"map": mapObj,
	"struct": obj,
	"slice": sliceObj,
}

const (
	GET string = "get"
	SET string = "set"
	DEL string = "delete"
	PRT string = "print"
	USE string = "use"
)

var target  interface{}= obj

// 记录command光标前一个单词
var lastWord string

func main() {
	p := prompt.New(
		executorFunc,
		completer,
		prompt.OptionTitle("expr: interactive Expr CLI"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.DarkGreen),
	)
	p.Run()
}

func executorFunc(command string) {
	command = strings.TrimSpace(command)
	if command == "" {
		return
	} else if command == "quit" || command == "exit" {
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	ops, fieldName, value := getOpsAndFieldNameAndValue(command)
	ops = strings.ToLower(ops)
	switch ops {
	case GET:
		field, err := expr.Get(&target, fieldName, strconv.Itoa(len(strings.Split(fieldName, "."))))
		if err != nil {
			log.Debugf("field %s not found, err: %s\n", fieldName, err)
			return
		}
		// todo: 删除log信息
		fmt.Printf("field %s: %+v\n", fieldName, field)
	case SET:

		err := expr.Set(&target, fieldName, value,nil, nil)
		if err != nil {
			// todo: 删除log信息
			fmt.Printf("field %s set failed, err: %s\n", fieldName, err)
			return
		}
		fmt.Println("ok!")
	case DEL:
		err := expr.Del(&target, fieldName)
		if err != nil {
			fmt.Printf("field %s delete failed, err: %s\n", fieldName, err)
			return
		}
		// todo: 删除log信息
		fmt.Println("ok!")
		log.Debugf("field %s delete success", fieldName)
	case PRT:
		fmt.Printf("obj: %+v\n", target)
	case USE:
		switch fieldName {
		case "slice":
			target = globalTestMap[fieldName]
			fallthrough
		case "map":
			target = globalTestMap[fieldName]
			fallthrough
		case "struct":
			target = globalTestMap[fieldName]
			fallthrough
		default:
			// todo: 删除log信息
			fmt.Println("ok!")
		}
	default:
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			// todo: 删除log信息
			fmt.Printf("Got error: %s\n", err.Error())
		}
	}
	return
}

// 操作命令的提示
func operationCompleter(d prompt.Document) []prompt.Suggest{
	ops := []prompt.Suggest{
		{Text: "print", Description: "打印当前操作对象"},
		{Text: "get", Description: "获取操作对象某个字段"},
		{Text: "set", Description: "设置操作对象某个字段上的值"},
		{Text: "delete", Description: "删除操作对象某个字段(map或slice)"},
		{Text: "use", Description: "切换当前对象(slice, map, struct)"},
		{Text: "quit", Description: "退出"},
	}
	if reflect.TypeOf(target).Kind() == reflect.Struct {
		ops = append(ops[:3], ops[4:]...)
	}
	return prompt.FilterHasPrefix(ops, d.GetWordBeforeCursor(), true)
}

// map or struct 的字段提示
func fieldNameCompleter(d prompt.Document,object interface{}) []prompt.Suggest{
	wordBeforeCursor := d.GetWordBeforeCursor()
	nestFieldNames := strings.Split(wordBeforeCursor, ".")
	switch reflect.TypeOf(object).Kind() {
	case reflect.Struct:
		dummyObject := object
		if len(nestFieldNames) != 1 {
			for index, fieldName := range nestFieldNames{
				if index == len(nestFieldNames)-1{
					continue
				}
				realObject := reflect.Indirect(reflect.ValueOf(dummyObject))
				nestField := realObject.FieldByName(fieldName)
				if !nestField.IsValid() {
					return nilCompleter(d)
				}
				dummyObject = nestField.Interface()
			}
		}
		return prompt.FilterHasPrefix(
			getFieldSuggest(getStructFieldNames(dummyObject)),
			nestFieldNames[len(nestFieldNames)-1], true)
	case reflect.Map:
		return prompt.FilterHasPrefix(
			getFieldSuggest(getMapKeyNames(object)),
			wordBeforeCursor, true)
	}
	return nilCompleter(d)
}

func objectNameCompleter(d prompt.Document) []prompt.Suggest{
	objectSuggests := []prompt.Suggest{
		{Text: "slice"},
		{Text: "map"},
		{Text: "struct"},
	}
	return prompt.FilterHasPrefix(objectSuggests, d.GetWordBeforeCursor(), true)
}

// 空提示
func nilCompleter(d prompt.Document)[]prompt.Suggest{
	return prompt.FilterHasPrefix(nil, d.GetWordBeforeCursor(), true)
}

// command 提示(main)
func completer(d prompt.Document) []prompt.Suggest {
	if strings.Index(d.TextBeforeCursor(), " ") < 0 {
		return operationCompleter(d)
	}else {
		lastWord = d.TextBeforeCursor()[:strings.LastIndex(d.TextBeforeCursor(), " ")]
	}
	// 数据操作命令之后提示字段
	if isDataOperation(lastWord){
		return fieldNameCompleter(d, target)
	}
	if isChangeOperation(lastWord) {
		return objectNameCompleter(d)
	}
	return nilCompleter(d)
}

func isDataOperation(ops string) bool {
	if ops == GET || ops == SET || ops == DEL{
		return true
	}
	return false
}

func isChangeOperation(ops string) bool {
	return ops == USE
}

func getFieldSuggest(fieldNames []string)[]prompt.Suggest{
	var fieldNameSuggest []prompt.Suggest
	for _, item := range fieldNames{
		fieldNameSuggest = append(fieldNameSuggest, prompt.Suggest{
			Text:        item,
		})
	}
	return fieldNameSuggest
}

// 从命令中获取操作符，字段名，值（设置时）
func getOpsAndFieldNameAndValue(token string) (ops, fieldName,value string){
	values := strings.Split(token, " ")
	if len(values) < 2 {
		return values[0], "", ""
	}
	if len(values) < 3 {
		return values[0], values[1], ""
	}
	return values[0], values[1], values[2]
}

// 获取结构体中的所有不为空的字段名
func getStructFieldNames(structObject interface{}) []string  {
	t := reflect.TypeOf(structObject)
	v := reflect.ValueOf(structObject)
	v = reflect.Indirect(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		//log.Println("Check type error not Struct")
		return []string{}
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i:= 0; i < fieldNum; i++ {
		if v.FieldByName(t.Field(i).Name).IsZero() {
			continue
		}
		result = append(result, t.Field(i).Name)
	}
	return result
}

// 获取map的所有key名，[$KeyName]格式
func getMapKeyNames(mapObject interface{}) []string {
	v := reflect.ValueOf(mapObject)
	keys := make([]string, 0, len(v.MapKeys()))
	for _, key := range v.MapKeys(){
		keys = append(keys, fmt.Sprintf("[%s]", key.Interface().(string)))
	}
	return keys
}