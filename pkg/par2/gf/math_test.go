package gf

import (
	"testing"
)

func TestAdd(t *testing.T) {
	res := Add(11, 7)
	if res != 12 {
		t.Fatalf("11 + 7 != 12. Actual: %d", res)
	}
}

func TestMul(t *testing.T) {
	res := Mul(11, 7)
	if res != 49 {
		t.Fatalf("11 * 7 != 4. Actual: %d", res)
	}
}

func TestMulByZero(t *testing.T) {
	res := Mul(11, 0)
	if res != 0 {
		t.Fatalf("Multiplication by zero gave nonzero result. Actual: %d", res)
	}
}

func TestMulOverflow(t *testing.T) {
	res := Mul(3, 3)
	if res != 5 {
		t.Fatalf("5487 * 39905 != 61347. Actual: %d", res)
	}
}

func TestDiv(t *testing.T) {
	res := Div(15, 7)
	if res != 27500 {
		t.Fatalf("11 / 7 != 4. Actual: %d", res)
	}
}

func TestBigDiv(t *testing.T) {
	res := Div(5487, 39905)
	if res != 29454 {
		t.Fatalf("5487 * 39905 != 61347. Actual: %d", res)
	}
}

func TestDivZeroBuSomethingIsZero(t *testing.T) {
	res := Div(0, 1)
	if res != 0 {
		t.Fatalf("0 / 1 != 0. Actual: %d", res)
	}
}

func TestDivByZeroPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Divide by 0 did not panic")
		}
	}()
	Div(11, 0)
}
