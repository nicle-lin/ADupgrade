package main

import (
	"reflect"
	"fmt"
)
type S struct{
	good string
}
type T struct {
	S
	bad string
}

func (S) sVal(){
	println("good S")
}
func (*S)sPtr(){
	println("good ss")
}

func(T)tVal(){
	println("bad T")
}
func (*T)tPtr(){
	println("bad TT")
}

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
