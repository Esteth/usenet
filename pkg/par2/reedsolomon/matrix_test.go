package reedsolomon

import (
	"reflect"
	"testing"
)

func TestVandermondePar2Matrix(t *testing.T) {
	m, err := NewVandermondePar2Matrix(5, 6)
	if err != nil {
		t.Fatalf("failed to create Vandermonde matrix: %v", err)
	}
	expected, err := NewMatrixData(
		[][]uint16{
			{1, 1, 1, 1, 1, 1},
			{1, 2, 4, 16, 128, 256},
			{1, 4, 16, 256, 16384, 4107},
			{1, 8, 64, 4096, 8566, 7099},
			{1, 16, 256, 4107, 43963, 7166},
		},
	)
	if err != nil {
		t.Fatalf("failed to create expected matrix: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("created vandermonde matrix not as expected: %v", m)
	}
}

func TestMultiplyByIdentityMatrix(t *testing.T) {
	m, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{2, 3, 5},
			{3, 4, 5},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	id, err := IdentityMatrix(3)
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
	m1, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	m2, err := NewMatrixData(
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
	expected, err := NewMatrixData(
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

func TestAugment(t *testing.T) {
	m1, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	m2, err := NewMatrixData(
		[][]uint16{
			{2, 5},
			{3, 6},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	augmentResult, err := m1.Augment(m2)
	if err != nil {
		t.Fatalf("could not multiply: %v", err)
	}
	expected, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3, 2, 5},
			{5, 7, 8, 3, 6},
		},
	)
	if err != nil {
		t.Fatalf("could not create expected result matrix: %v", err)
	}
	if !reflect.DeepEqual(augmentResult, expected) {
		t.Fatalf("Result matrix %v not expected %v", augmentResult, expected)
	}
}

func TestAugmentVertical(t *testing.T) {
	m1, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	m2, err := NewMatrixData(
		[][]uint16{
			{2, 5, 1},
			{3, 6, 0},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	augmentResult, err := m1.AugmentVertical(m2)
	if err != nil {
		t.Fatalf("could not multiply: %v", err)
	}
	expected, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
			{2, 5, 1},
			{3, 6, 0},
		},
	)
	if err != nil {
		t.Fatalf("could not create expected result matrix: %v", err)
	}
	if !reflect.DeepEqual(augmentResult, expected) {
		t.Fatalf("Result matrix %v not expected %v", augmentResult, expected)
	}
}

func TestAugmentVerticalCopies(t *testing.T) {
	m1, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	m2, err := NewMatrixData(
		[][]uint16{
			{2, 5, 1},
			{3, 6, 0},
		},
	)
	if err != nil {
		t.Fatalf("could not create matrix: %v", err)
	}
	augmentResult, err := m1.AugmentVertical(m2)
	if err != nil {
		t.Fatalf("could not multiply: %v", err)
	}
	expected, err := NewMatrixData(
		[][]uint16{
			{1, 2, 3},
			{5, 7, 8},
			{2, 5, 1},
			{3, 6, 0},
		},
	)
	if err != nil {
		t.Fatalf("could not create expected result matrix: %v", err)
	}
	// Mutate an element of one of the original matrices before comparison
	m2.data[0] = 0
	if !reflect.DeepEqual(augmentResult, expected) {
		t.Fatalf("Result matrix %v not expected %v", augmentResult, expected)
	}
}

func TestGaussianElimination(t *testing.T) {
	m, err := NewMatrixData(
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

	expected, err := NewMatrixData(
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

func TestPlankPaperErrorRecovery(t *testing.T) {
	data := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	checksums := []uint16{11, 60570, 57778}
	i, err := IdentityMatrix(len(data))
	if err != nil {
		t.Fatalf("could not create identity matrix: %v", err)
	}
	f, err := NewVandermondePar2Matrix(3, len(data))
	if err != nil {
		t.Fatalf("could not create vandermonde matrix: %v", err)
	}
	a, err := i.AugmentVertical(f)
	if err != nil {
		t.Fatalf("could not stack identity and vandermonde matrix: %v", err)
	}
	e, err := NewMatrixColumn(append(data, checksums...))
	if err != nil {
		t.Fatalf("could not create column with data and checksums: %v", err)
	}
	solve, err := a.Augment(e)
	if err != nil {
		t.Fatalf("could not create problem matrix: %v", err)
	}

	// Delete some rows to pretend we lost some data
	copy(solve.data[4*solve.cols:], solve.data[7*solve.cols:])
	solve.rows -= 3

	err = solve.GaussianElimination()
	if err != nil {
		t.Fatalf("could not solve problem matrix: %v", err)
	}

	recoveredData := make([]uint16, len(data))
	for r := 0; r < solve.rows; r++ {
		recoveredData[r] = solve.cell(r, solve.cols-1)
	}

	if !reflect.DeepEqual(recoveredData, data) {
		t.Fatalf("recovered data %v not equal to expected data %v", recoveredData, data)
	}
}

func BenchmarkAugment(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m1 := NewMatrix(1024, 1024)
		m2 := NewMatrix(7, 1024)
		m1.Augment(m2)
	}
}

func BenchmarkAugmentVertical(b *testing.B) {
	for n := 0; n < b.N; n++ {
		m1 := NewMatrix(1024, 1024)
		m2 := NewMatrix(1024, 7)
		m1.AugmentVertical(m2)
	}
}
