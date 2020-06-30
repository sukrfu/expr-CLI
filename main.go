package main


import (
	"bufio"
	expr "expr"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Obj struct {
	Name   string
	Friend *Obj
	Id     int64
}

var obj = &Obj{
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
	DEL string = "del"
	PRT string = "prt"
	USE string = "use"
)

var target = obj

const Delim = '\n'

func main() {
	fmt.Println("************************************")
	fmt.Println("command: ops + [name] + [value]")
	fmt.Println("get: 获取操作对象某个字段")
	fmt.Println("set: 设置操作对象某个字段上的值")
	fmt.Println("del: 删除操作对象某个字段(map或slice)")
	fmt.Println("use: 切换操作对象(map, slice, struct)")
	fmt.Println("prt: 打印当前操作对象")
	fmt.Println("************************************")

	reader := bufio.NewReader(os.Stdin)
	for {
		command, err := reader.ReadString(Delim)
		command = removeDelim(command)
		if err != nil {
			continue
		}

		ops, fieldName, value := getOpsAndFieldNameAndValue(command)
		ops = strings.ToLower(ops)
		switch ops {
		case GET:
			field, err := expr.GetField(&target, fieldName)
			if err != nil {
				log.Debugf("field %s not found", fieldName)
				continue
			}
			fmt.Printf("field %s: %+v\n", fieldName, field)
		case SET:
			if reflect.TypeOf(target).Elem().Kind()==reflect.Int {
				value,err = strconv.Atoi(value.(string))
				if err != nil {
					fmt.Printf("field %s set failed, err: %s\n", fieldName, err)
					continue
				}
			}
			expr.SetField(&target, fieldName, value, nil)
		case DEL:
			err := expr.Del(&target, fieldName)
			if err != nil {
				fmt.Printf("field %s delete failed, err: %s\n", fieldName, err)
				continue
			}
			log.Debugf("field %s delete success", fieldName)
		case PRT:
			fmt.Printf("obj: %+v\n", target)
		//case USE:
		//	objType := fieldName
		//	target = globalTestMap[objType]
		default:
			log.Debugf("operation %s invalid", ops)
		}
	}
}

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


