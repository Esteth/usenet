package reedsolomon

import (
	"reflect"
	"testing"
)

func TestMultiplyByIdentityMatrix(t *testing.T) {
	m, err := newMatrixData(
		[][]uint16{
			{1, 2, 3},
			{2, 3, 5},
			{3, 4, 5},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	id, err := identityMatrix(3)
	if err != nil {
		t.Fatalf("could not create identity matrix: %v", err)
	}
	mulResult, err := m.Mul(id)
	if err != nil {
		t.Fatalf("could not multiply: %v", err)
	}
	if !reflect.DeepEqual(m, mulResult) {
		t.Fatalf("Identity multiplication does not give original back")
	}
}

func TestMultiply(t *testing.T) {
	m1, err := newMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	m2, err := newMatrixData(
		[][]uint16{
			{2, 5},
			{3, 6},
			{4, 7},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	mulResult, err := m1.Mul(m2)
	if err != nil {
		t.Fatalf("could not multiply: %v", err)
	}
	expected, err := newMatrixData(
		[][]uint16{
			{8, 0},
			{35, 59},
		},
	)
	if err != nil {
		t.Fatalf("could not create expected result matrix: %v", err)
	}
	if !reflect.DeepEqual(mulResult, expected) {
		t.Fatalf("Result matrix %v not expected %v", mulResult, expected)
	}
}

func TestGaussianEliminationSimple(t *testing.T) {
	m, err := newMatrixData(
		[][]uint16{
			{4, 2, 3, 1},
			{2, 3, 5, 0},
			{3, 4, 5, 0},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	err = m.GaussianElimination()
	if err != nil {
		t.Fatalf("gaussian elimination failed: %v", err)
	}

	expected, err := newMatrixData(
		[][]uint16{
			{1, 0, 0, 43393},
			{0, 1, 0, 14427},
			{0, 0, 1, 21091},
		},
	)
	if err != nil {
		t.Fatalf("could not create expected matrix: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("Elimination did not provide correct results: %v", m)
	}
}
