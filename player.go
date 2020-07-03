package main

import (
	"fmt"
	"gitee.com/legou-lib/expr"
	"strconv"
	"strings"
)

type Player struct {
	Name string
	Id string
	Coin int
	Friends []*Player
}

func (player *Player)Get(fieldName string) (string, error){
	return expr.Get(player, fieldName, strconv.Itoa(len(strings.Split(fieldName, "."))))
}

func (player *Player)Set(fieldName ,valueEval string) error{
	return expr.Set(player, fieldName, valueEval,nil, nil)
}

func (player *Player)Delete(key string) error{
	return expr.Del(player, key)
}

func (player *Player)Print(){
	object, _ :=player.Get("")
	fmt.Println(object)
}


