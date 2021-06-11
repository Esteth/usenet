package par2

import "testing"

func TestFirstFewConstants(t *testing.T) {
	expected := []uint16{
		2,
		4,
		16,
		128,
		256,
		2048,
		8192,
		16384,
		4107,
		32856,
		17132,
	}
	pool := NewConstantPool()
	for i, expected := range expected {
		constant := pool.Next()
		if constant != expected {
			t.Fatalf("failed at index %d: actually %d, expected %d", i, constant, expected)
		}
	}
}
