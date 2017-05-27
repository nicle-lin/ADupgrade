package main

type good interface {
	string2()
}

type tester interface {
	test()
	string() string
	//good()
	//name string
	good
}

/*
func (tester)update(){
	println("update")
}
*/

type data struct {

}

func (*data) test(){ println("test")}
func (data) string() string {
	println("string")
	return ""
}
func (data) string2(){
	println("string2")
}


func PP(a good){
	a.string2()
}


func main() {

	var d data

	var t tester = &d
	//var t tester = d

	t.test()
	t.string()
	t.string2()
	println("=======================")
	PP(t)
	var s  good = t

	PP(s)
	//t.good()
	//t.update()

	var t1, t2  interface{}

	println(t1==nil, t2 == nil)

	t1, t2 = 100, 100

	println(t1 == t2)

	//t1, t2 = map[string]int{}, map[string]int{}
	//println(t1 == t2)
}
