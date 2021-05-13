package gf

import (
	"testing"
)

func TestAddIsXor(t *testing.T) {
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

func TestDiv(t *testing.T) {
	res := Div(11, 7)
	if res != 55007 {
		t.Fatalf("11 / 7 != 4. Actual: %d", res)
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
