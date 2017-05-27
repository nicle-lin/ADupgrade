package main


type TestError struct {

}

func (*TestError) Error() string{
	return "error"
}

func test(x int) (int, error) {
	var err *TestError
	if x < 0{
		err = new(TestError)
		x = 0
	}else{
		x += 100
	}
	return x, err
}

func main() {
	//x, err := test(100)
	x, err := test(-100)
	if err != nil {
		println("err != nil")
	}
	println(err.(*TestError))

	println(x)
}
