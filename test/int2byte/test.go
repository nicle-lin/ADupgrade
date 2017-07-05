package main

import "fmt"

func main() {
	value := IntToByte(0xfffefd)
	fmt.Println(value)
	fmt.Println(((0xfefdfc) >> 16) & 0xff)
	fmt.Println(1<<10)
}

func IntToByte(data int64) []byte {
	var result []byte
	mask := int64(0xFF)
	fmt.Println("mask:", mask)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((data>>shift)&mask))
		fmt.Println("shift:", byte(data>>shift))
		fmt.Println("shift and mask:", byte((data>>shift)&mask))
	}
	return result
}
