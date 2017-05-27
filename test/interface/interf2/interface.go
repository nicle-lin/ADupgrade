
package main

import "fmt"

type Person struct {
	name string
	age int
}

func (p Person) printMsg() {
	fmt.Printf("I am %s, and my age is %d.\n", p.name, p.age)
}

func (p Person) eat(s string) {
	fmt.Printf("%s is eating %s ...\n", p.name, s)
}

func (p Person) drink(s string) {
	fmt.Printf("%s is drinking %s ...\n", p.name, s)
}

type People interface {
	printMsg()
	PeopleEat    //组合
	PeopleDrink
	//eat() //不能出现重复的方法
}
/*  
//与上面等价  
type People interface {  
    printMsg()  
    eat()  
    drink()  
}  
*/
type PeopleDrink interface {
	drink(s string)
}

type PeopleEat interface {
	eat(s string)
}

type PeopleEatDrink interface {
	eat(s string)
	drink(s string)
}

//以上 Person 类[型]就实现了 People/PeopleDrink/PeopleEat/PeopleEatDrink interface 类型  

type Foodie struct {
	name string
}

func (f Foodie) eat(s string) {
	fmt.Printf("I am foodie, %s. My favorite food is the %s.\n", f.name, s)
}

//Foodie 类实现了 PeopleEat interface 类型  

func main() {
	//定义一个 People interface 类型的变量p1
	var p1 People
	p1 = Person{"Rain", 23}
	p1.printMsg()           //I am Rain, and my age is 23.
	p1.drink("orange juice")//print result: Rain is drinking orange juice

	//同一类可以属于多个 interface, 只要这个类实现了这个 interface中的方法
	var p2 PeopleEat
	p2 = Person{"Sun", 24}
	p2.eat("chaffy dish")//print result: Sun is eating chaffy dish ...

	//不同类也可以实现同一个 interface
	var p3 PeopleEat
	p3 = Foodie{"James"}
	p3.eat("noodle")//print result: I am foodie, James. My favorite food is the noodle

	//interface 赋值
	p3 = p1  //p3 中的方法会被 p1 中的覆盖
	p3.eat("noodle")
	/************************************/
	/*print result                      */
	/*Rain is eating noodle ...         */
	/************************************/

	//interface 查询
	//将(子集) PeopleEat 转为 People 类型
	if p4, ok := p2.(People); ok {
		p4.drink("water") //调用 People interface 中有而 PeopleEat 中没有的方法
		fmt.Println(p4)
	}
	/************************************/
	/*print result                      */
	/*Sun is drink water ...            */
	/*{Sun 24}                          */
	/************************************/

	//查询 p2 是否为 Person 类型变量
	if p5, ok := p2.(Person); ok {
		fmt.Println(p5, "type is Person")
		p5.drink("***")  //此时也可以调用 Person 所有的方法
	}
	/************************************/
	/*print result                      */
	/*{Sun 24} type is Person           */
	/*Sun is drink *** ...              */
	/************************************/

	var p6 PeopleEat = Foodie{"Tom"}

	if p7, ok := p6.(People); ok {
		fmt.Println(p7)
	} else {
		fmt.Println("Error: can not convert")
	}
	//result: Error: can not convert

	if p8, ok := p6.(Foodie); ok {
		fmt.Println(p8, "type is Foodie")
	}
	//result: {Tom} type is Foodie
}  