package main

import (
	"expr"
	"fmt"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

type Obj struct {
	Name   string
	Friend *Obj
	Id     int64
}

var obj = Obj{
	Name: "wuwj",
	Id:   22,
	Friend: &Obj{
		Name: "hello",
		Id:   3,
	},
}

var mapObj = map[string]string{
	"hello": "123",
	"test": "456",
}

var sliceObj = []int{1,2,3,4,5}


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

var target = obj

const Delim = '\n'

var lastWord string

func completer(d prompt.Document) []prompt.Suggest {
	ops := []prompt.Suggest{
		{Text: "print", Description: "打印当前操作对象"},
		{Text: "get", Description: "获取操作对象某个字段"},
		{Text: "set", Description: "设置操作对象某个字段上的值"},
		{Text: "delete", Description: "删除操作对象某个字段(map或slice)"},
		{Text: "quit", Description: "退出"},
	}


	if strings.Index(d.TextBeforeCursor(), " ") < 0 {
		return prompt.FilterHasPrefix(ops, d.GetWordBeforeCursor(), true)
	}else {
		lastWord = d.TextBeforeCursor()[:strings.LastIndex(d.TextBeforeCursor(), " ")]
	}

	if reflect.TypeOf(target).Kind() != reflect.Struct {
		return nil
	}

	// must be struct
	fieldNames := getFieldName(target)
	wordBeforeCursor := d.GetWordBeforeCursor()
	fieldNameSuggest := getFieldSuggest(fieldNames)

	if strings.Contains(wordBeforeCursor, ".") {
		currFieldName := wordBeforeCursor[:strings.Index(wordBeforeCursor,".")]
		currField, err := expr.GetField(target, currFieldName)
		if err != nil {
			fmt.Printf("field %s get failed, err: %s\n", currFieldName, err)
			os.Exit(0)
		}
		if currField.Type().Kind() != reflect.Ptr {
			return prompt.FilterHasPrefix(nil, wordBeforeCursor, true)
		}
		fieldNames = getFieldName(currField.Interface())
		fieldNameSuggest = getFieldSuggest(fieldNames)
		return prompt.FilterHasPrefix(fieldNameSuggest, d.GetWordBeforeCursorUntilSeparator("."), true)
	}
	if isStructOps(lastWord){
		return prompt.FilterHasPrefix(fieldNameSuggest, wordBeforeCursor, true)
	}
	return prompt.FilterHasPrefix(nil, wordBeforeCursor, true)
}

func isStructOps(lastWord string) bool{
	if lastWord == GET || lastWord == SET || lastWord == DEL{
		return true
	}
	return false
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
		field, err := expr.GetField(&target, fieldName)
		if err != nil {
			log.Debugf("field %s not found", fieldName)
			return
		}
		fmt.Printf("field %s: %+v\n", fieldName, field)
	case SET:
		if reflect.TypeOf(target).Kind()==reflect.Slice &&
			reflect.TypeOf(target).Elem().Kind()==reflect.Int {
			var err error
			value,err = strconv.Atoi(value.(string))
			if err != nil {
				fmt.Printf("field %s set failed, err: %s\n", fieldName, err)
				return
			}
		}
		expr.SetField(&target, fieldName, value, nil)
	case DEL:
		err := expr.Del(&target, fieldName)
		if err != nil {
			fmt.Printf("field %s delete failed, err: %s\n", fieldName, err)
			return
		}
		log.Debugf("field %s delete success", fieldName)
	case PRT:
		fmt.Printf("obj: %+v\n", target)
	//case USE:
	//	objType := fieldName
	//	target = globalTestMap[objType]
	default:
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Got error: %s\n", err.Error())
		}
	}
	return
}

func main() {
	p := prompt.New(
		executorFunc,
		completer,
		prompt.OptionTitle("expr: interactive Expr CLI"),
		prompt.OptionPrefix(">>> "),
		prompt.OptionInputTextColor(prompt.Black),
	)
	p.Run()
}
//
//func main() {
//	fmt.Println("************************************")
//	fmt.Println("command: ops + [name] + [value]")
//	fmt.Println("get: 获取操作对象某个字段")
//	fmt.Println("set: 设置操作对象某个字段上的值")
//	fmt.Println("del: 删除操作对象某个字段(map或slice)")
//	fmt.Println("use: 切换操作对象(map, slice, struct)")
//	fmt.Println("prt: 打印当前操作对象")
//	fmt.Println("************************************")
//
//	reader := bufio.NewReader(os.Stdin)
//	for {
//		command, err := reader.ReadString(Delim)
//		command = removeDelim(command)
//		if err != nil {
//			continue
//		}
//
//		ops, fieldName, value := getOpsAndFieldNameAndValue(command)
//		ops = strings.ToLower(ops)
//		switch ops {
//		case GET:
//			field, err := expr.GetField(&target, fieldName)
//			if err != nil {
//				log.Debugf("field %s not found", fieldName)
//				continue
//			}
//			fmt.Printf("field %s: %+v\n", fieldName, field)
//		case SET:
//			if reflect.TypeOf(target).Elem().Kind()==reflect.Int {
//				value,err = strconv.Atoi(value.(string))
//				if err != nil {
//					fmt.Printf("field %s set failed, err: %s\n", fieldName, err)
//					continue
//				}
//			}
//			expr.SetField(&target, fieldName, value, nil)
//		case DEL:
//			err := expr.Del(&target, fieldName)
//			if err != nil {
//				fmt.Printf("field %s delete failed, err: %s\n", fieldName, err)
//				continue
//			}
//			log.Debugf("field %s delete success", fieldName)
//		case PRT:
//			fmt.Printf("obj: %+v\n", target)
//		//case USE:
//		//	objType := fieldName
//		//	target = globalTestMap[objType]
//		default:
//			log.Debugf("operation %s invalid", ops)
//		}
//	}
//}

func getOpsAndFieldNameAndValue(token string) (ops, fieldName string, value interface{}){
	values := strings.Split(token, " ")
	if len(values) < 2 {
		return values[0], "", nil
	}
	if len(values) < 3 {
		return values[0], values[1], nil
	}
	return values[0], values[1], values[2]
}

func removeDelim(s string) string {
	return s[: strings.Index(s, string(Delim))]
}

// 获取结构体中的所有字段名
func getFieldName(structName interface{}) []string  {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i:= 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}