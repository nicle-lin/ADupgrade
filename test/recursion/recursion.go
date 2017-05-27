package main

import "fmt"

func main() {
	fmt.Println(recursion(4))
	for n,i := range GetList(){
		fmt.Println("n:",n, "i:",i)
	}
}

var i int
func recursion(n int) (sum int) {
	i++
	fmt.Println("i:",i)
	if n == 1 {
		return 1
	}
	sum = n + recursion(n-1)
	return sum
}
/*
func recursion2(n int) (sum int){
	switch n {
	case :

	}
}
*/
func GetList()[]string{
	fmt.Println("in the GetList")
	return []string{"a", "b","c"}
}
