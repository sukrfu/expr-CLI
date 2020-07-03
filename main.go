package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

var currentObjectPtr SimpleData

// 记录command光标前一个单词
var lastWord string

var operation = NewSimpleDataOperation()

var registry = NewRegistry()

func main() {
	player1 := &Player{
		Name:    "tom",
		Id:      "654321",
		Coin:   123,
		Friends: nil,
	}
	player := &Player{
		Name:    "nick",
		Id:      "12345",
		Coin:   66,
		Friends: []*Player{player1},
	}
	registry.Register("tom", player1)
	registry.Register("nick", player)
	registry.Register("player3", player)
	registry.Register("player4", player)

	p := prompt.New(
		executorFunc,
		completer,
		prompt.OptionTitle("expr: interactive Expr CLI"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.DarkGreen),
	)
	p.Run()
}

func commandErrorHandler(){
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("command error, err: %s\n", err)
		}
	}()
}

func printSysCmdHelp(){
	fmt.Println("   list: 列出当前注册中心所有对象")
	fmt.Println("   quit: 退出程序")
	fmt.Println("command: [object_name].[ops] + [field_name] + [value]")
}

func executorFunc(command string) {
	commandErrorHandler()
	command = strings.TrimSpace(command)
	switch strings.ToLower(command) {
	case "":
		return
	case "?":
		printSysCmdHelp()
		return
	case "list":
		fmt.Println(registry.GetAllNames())
		return
	case "quit":
		fmt.Println("Bye!")
		os.Exit(0)
		return
	}

	ops, fieldName, value := getOpsAndFieldNameAndValue(command)
	ops = strings.ToLower(ops)
	switch ops {
	case operation.GET:
		field, err := operation.Get(currentObjectPtr, fieldName)
		if err != nil {
			log.Debugf("field %s not found, err: %s\n", fieldName, err)
			return
		}
		fmt.Printf("field %s: %+v\n", fieldName, field)
	case operation.SET:
		err := operation.Set(currentObjectPtr, fieldName, value)
		if err != nil {
			fmt.Printf("field %s set failed, err: %s\n", fieldName, err)
			return
		}
		fmt.Println("ok!")
	case operation.DEL:
		err := operation.Delete(currentObjectPtr, fieldName)
		if err != nil {
			fmt.Printf("field %s delete failed, err: %s\n", fieldName, err)
			return
		}
		fmt.Println("ok!")
	case operation.PRT:
		operation.Print(currentObjectPtr)
	default:
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("   input ? for help")
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
	}
	if operation.GetObjectRuntimeType(currentObjectPtr) == reflect.Struct {
		ops = append(ops[:3], ops[4:]...)
	}
	return prompt.FilterHasPrefix(ops, d.GetWordBeforeCursorUntilSeparator("."), true)
}

// map or struct 的字段提示
func fieldNameCompleter(d prompt.Document,object interface{}) []prompt.Suggest{
	if currentObjectPtr == nil {
		return []prompt.Suggest{{Text: "当前操作对象为空"} }
	}
	wordBeforeCursor := d.GetWordBeforeCursor()
	nestFieldNames := strings.Split(wordBeforeCursor, ".")
	switch operation.GetObjectRuntimeType(currentObjectPtr) {
	case reflect.Struct:
		dummyObject := object
		if strings.Contains(wordBeforeCursor, ".") {
			nestFieldName := wordBeforeCursor[:strings.LastIndex(wordBeforeCursor, ".")]
			nestField := operation.GetField(object.(SimpleData), nestFieldName)
			if nestField == nil {
				return nilCompleter(d)
			}
			dummyObject = nestField
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
	names := registry.GetAllNames()
	objectSuggests := make([]prompt.Suggest, 0 , len(names))
	for _, name := range names{
		objectSuggests = append(objectSuggests, prompt.Suggest{
			Text:        name,
			Description: operation.GetObjectFieldType(registry.GetObject(name), ""),
		})
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
		if strings.Count(d.TextBeforeCursor(), ".") == 1 {
			currentObjectPtr = registry.GetObject(
				d.TextBeforeCursor()[:strings.Index(d.TextBeforeCursor(), ".")])
			return operationCompleter(d)
		}
		return objectNameCompleter(d)
	}else {
		lastWord = d.TextBeforeCursor()[
		strings.Index(d.TextBeforeCursor(), ".")+1 :
			strings.LastIndex(d.TextBeforeCursor(), " ")]
	}
	// 数据操作命令之后提示字段
	if isDataOperation(lastWord){
		return fieldNameCompleter(d, currentObjectPtr)
	}
	return nilCompleter(d)
}

func isDataOperation(ops string) bool {
	ops = strings.ToLower(ops)
	if ops == operation.GET || ops == operation.SET || ops == operation.DEL{
		return true
	}
	return false
}

func getFieldSuggest(fieldNames []string)[]prompt.Suggest{
	var fieldNameSuggest []prompt.Suggest
	for _, name := range fieldNames{
		fieldNameSuggest = append(fieldNameSuggest, prompt.Suggest{
			Text: name,
			Description: operation.GetObjectFieldType(currentObjectPtr, name),
		})
	}
	return fieldNameSuggest
}

// 从命令中获取操作符，字段名，值（设置时）
func getOpsAndFieldNameAndValue(token string) (ops, fieldName,value string){
	values := strings.Split(token, " ")
	if strings.Contains(values[0], ".") {
		ops = strings.Split(values[0], ".")[1]
	}
	if len(values) < 2 {
		return ops, "", ""
	}
	if len(values) < 3 {
		return ops, values[1], ""
	}
	return ops, values[1], values[2]
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
	v := reflect.Indirect(reflect.ValueOf(mapObject))
	keys := make([]string, 0, len(v.MapKeys()))
	for _, key := range v.MapKeys(){
		keys = append(keys, fmt.Sprintf("[%s]", key.Interface().(string)))
	}
	return keys
}