package test

import (
	"testing"
	"time"
)
func BenchmarkAbsolute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Absolute(1,2)
	}
}

func TestAbsolute1(t *testing.T) {
	var tests = []struct{
		x int
		y int
		expect int
	}{
		{1,1,2},
		{0,0,0},
		{-1,-1,2},
		{-2,-2, 4},
		{-2,-1, 3},
	}
	for _, tt := range tests{
		if Absolute(tt.x, tt.y) != tt.expect{
			t.Errorf("error, %d + %d = %d",tt.x, tt.y, tt.expect)
		}
	}

}

func TestAbsolute(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Absolute(-1, -2) != 3{
		t.FailNow()
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}
func TestCompare1(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}

func TestCompare2(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}
func TestCompare3(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}
func TestCompare4(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}
func TestCompare5(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}
func TestCompare6(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
	if Compare(1, 2) != false {
		t.FailNow()
	}
}
