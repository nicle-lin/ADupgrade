package test

func Absolute(x,y int) int{
	if x > 0 && y > 0 {
		return x + y
	}
	if x > 0 && y < 0{
		return x - y
	}
	if x < 0 && y > 0{
		return y - x
	}
	if x < 0 && y < 0 {
		return -(x + y)
	}
	return x + y
}

func Compare(x, y int) bool{
	if x == y {
		return true
	}
	return false
}
