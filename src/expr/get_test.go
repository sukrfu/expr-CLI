package expr

import (
	"fmt"
	"reflect"
	"testing"
)

type Pb struct {
	A int16
	B int32
}
type address struct {
	Streat  string
	ZipCode string
	State   string
	City    []*city
	Nums    []Num
	Emails  []string
	Maps    map[string]interface{}
	Ph      map[int16]Pb
	In      interface{}
}

type Num struct {
	Name string
	InUS bool
	N    int16
	Park park
}
type city struct {
	Name string
	InUS bool
	Park *park
}

type park struct {
	Name     string
	Location string
	Maps     map[string]Pb
	Emails   []string
}

func TestGetStruct(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77479"
	s.State = "Taxes"
	s.Emails = []string{"123@123.com", "456@456.com"}
	m := make(map[string]interface{})
	m["dd"] = "dd"
	m["cc"] = "cc"
	m["bb"] = "bb"
	m1 := make(map[string]Pb)

	s.Maps = m

	s.City = append(s.City, &city{Name: "Sugar Land", InUS: true, Park: &park{Name: "Name", Location: "location", Maps: m1}})
	s.In = "string222"

	field, err := GetField(s, "Emails")
	if err != nil {
		panic(err)
	}

	fmt.Println("Emails:", FormatObj(field))

	field, err = GetField(s, "Emails[0]")
	if err != nil {
		panic(err)
	}

	fmt.Println("Emails[0]:", FormatObj(field))
	field, err = GetField(s, `City[0].Park.Maps`)
	if err != nil {
		panic(err)
	}
	fmt.Println("city[0].Park.Maps:", FormatObj(field))
}

//////////////////////////
type CharData struct {
	Name string
}

type CharOp interface {
	GetId() int64
}

type ModChar struct {
	Char []CharOp
}

type Char1 struct {
	Id   int64
	Data *CharData
}

func (c *Char1) GetId() int64 {
	return c.Id
}

func TestInterface(t *testing.T) {

	fun := func(v reflect.Value) reflect.Value {
		switch v.Type().Name() {
		case "CharOp":
			msg := v.Interface().(*Char1)
			return reflect.ValueOf(msg)
		}

		return v
	}

	mods := &ModChar{
		Char: make([]CharOp, 0),
	}
	c1 := &Char1{1, &CharData{"1"}}
	mods.Char = append(mods.Char, c1)
	field, err := GetField(mods, `Char[0].Data`)
	if err != nil {
		panic(err)
	}
	fmt.Println(FormatObj(field))
	err = SetField(mods, "Char[0].Data.Name", "555@125553.com", fun)
	if err != nil {
		panic(err)
	}

	field, err = GetField(mods, `Char[0].Data`)
	if err != nil {
		panic(err)
	}
	fmt.Println(FormatObj(field))
}

/////////////////////////////////////////////////

type MapSlice struct {
	Datas map[int32][]int32
}

func TestMapSlice(t *testing.T) {
	data := &MapSlice{
		Datas: make(map[int32][]int32),
	}
	ls := []int32{1, 2, 3}
	data.Datas[1] = ls

	field, err := GetField(data, `Datas[1][2]`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("=======%v\n", field)
	fmt.Println(FormatObj(field))

	field, err = GetField(ls, `[1]`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("=======%v\n", field)
	fmt.Println(FormatObj(field))

	ls2 := [][]int32{[]int32{1, 2, 3}}
	field, err = GetField(ls2, `[0][1]`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("=======%v\n", field)
	fmt.Println(FormatObj(field))

	ls3 := [][]MapSlice{[]MapSlice{*data}}
	field, err = GetField(ls3, `[0][0]Datas[1][2]`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("=======%v\n", field)
	fmt.Println(FormatObj(field))

}
