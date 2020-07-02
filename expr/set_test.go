package expr

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSetSlice(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77477"

	err := SetField(s, "Emails[]", []string{"yy123@123.com", "456@456.com"}, nil)
	if err != nil {
		panic(err)
	}
	err = SetField(s, "Emails[0]", "555@125553.com", nil)
	err = SetField(s, "Emails[1]", "555@125553.com", nil)

	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(s)
	fmt.Println(string(v))
}

type B struct {
	Num Num
}

func TestSetMap(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77477"
	m := make(map[string]interface{})
	m["dd"] = "dd"
	m["cc"] = "cc"
	m["bb"] = "bb"
	p := make(map[int16]Pb)
	err := SetField(s, "Maps", m, nil)
	if err != nil {
		panic(err)
	}
	err = SetField(s, "Ph", p, nil)
	if err != nil {
		panic(err)
	}

	err = SetField(s, "Streat", "ddddd", nil)
	if err != nil {
		panic(err)
	}

	err = SetField(s, "Ph[1]", Pb{1, 1}, nil)
	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(s)
	fmt.Println(string(v))
	err = SetField(s, "Ph[1].A", int16(188), nil)
	if err != nil {
		panic(err)
	}
	mapinfo := make(map[string]Pb)
	mapinfo["1"] = Pb{A: 1}
	err = SetField(s, "Nums[]", []Num{Num{Name: "jake1", Park: park{Maps: mapinfo}}}, nil)
	if err != nil {
		panic(err)
	}

	err = SetField(s, "Nums", Num{"jake2", false, 1, park{Name: "new york"}}, nil)
	if err != nil {
		panic(err)
	}
	err = SetField(s, "Nums", Num{"jake3", false, 1, park{Name: "new york"}}, nil)
	if err != nil {
		panic(err)
	}

	err = SetField(s, "Nums", Num{"jake4", false, 1, park{Name: "new york"}}, nil)
	if err != nil {
		panic(err)
	}
	v, _ = json.Marshal(s)
	fmt.Println(string(v))
	err = SetField(s, "Nums[1].N", int16(1), nil)
	if err != nil {
		panic(err)
	}

	err = Del(s, "Ph[1]")
	if err != nil {
		panic(err)
	}

	v, _ = json.Marshal(s)
	fmt.Println(string(v))
	err = Del(s, "Nums[1]")
	if err != nil {
		panic(err)
	}

	err = Del(s, "Nums[2]")
	if err != nil {
		panic(err)
	}

	str, err := Get(s, "Nums[0].Park")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(str))
	err = Set(s, "Nums[0].Park.Maps[1].A", "\"122\"", nil, nil)
	if err != nil {
		panic(err)
	}
	err = Set(s, "Nums[0].Park.Maps[1].B", "\"123\"", nil, nil)
	if err != nil {
		panic(err)
	}
	v, _ = json.Marshal(s)
	fmt.Println(string(v))

	err = Del(s, "Nums[0].Park.Maps[1]")
	if err != nil {
		panic(err)
	}

	err = Set(s, "Nums[0].Name", "\"sfsdf\"", nil, nil)
	if err != nil {
		panic(err)
	}
	v, _ = json.Marshal(s)
	fmt.Println(string(v))
}
