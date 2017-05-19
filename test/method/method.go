package main

import (
	"reflect"
	"fmt"
)
type S struct{}
type T struct {
	S
}

func (S) sVal(){}
func (*S)sPtr(){}

func(T)tVal(){}
func (*T)tPtr(){}

func methodSet(a interface{}){
	t := reflect.TypeOf(a)
	fmt.Println("t.NumMethod:",t.NumMethod())
	for i, n := 0, t.NumMethod(); i< n; i++{
		m := t.Method(i)
		fmt.Println(m.Name,m.Type)
	}
}

func main() {

	var t T
	methodSet(t)
	fmt.Println("-------------------")
	methodSet(&t)
}
