package reedsolomon

import (
	"errors"
	"fmt"

	"github.com/esteth/usenet/pkg/par2/gf"
)

// Matrix type heavily inspired by github.com/klauspost/reedsolomon
type matrix [][]uint16

var errInvalidRowSize = errors.New("invalid row size")
var errInvalidColSize = errors.New("invalid col size")
var errRowSizeMismatch = errors.New("row size is not the same for both matrices")
var errColSizeMismatch = errors.New("column size is not the same for all rows")
var errSingular = errors.New("cannot solve a singular matrix")

func newMatrix(rows, cols int) (matrix, error) {
	if rows <= 0 {
		return nil, errInvalidRowSize
	}
	if cols <= 0 {
		return nil, errInvalidColSize
	}
	m := matrix(make([][]uint16, rows))
	for i := range m {
		m[i] = make([]uint16, cols)
	}
	return m, nil
}

func newMatrixData(data [][]uint16) (matrix, error) {
	m := matrix(data)
	err := m.Check()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func newMatrixColumn(data []uint16) (matrix, error) {
	if len(data) <= 0 {
		return nil, errInvalidColSize
	}
	m := matrix(make([][]uint16, len(data)))
	for i := range m {
		m[i] = []uint16{data[i]}
	}
	return m, nil
}

func newVandermondeMatrix(rows, cols int) (matrix, error) {
	m, err := newMatrix(rows, cols)
	if err != nil {
		return nil, err
	}
	for r, row := range m {
		for c := range row {
			if c == 0 || r == 0 {
				m[r][c] = 1
				continue
			}
			if r == 1 {
				m[r][c] = uint16(c + 1)
				continue
			}
			m[r][c] = gf.Mul(m[r-1][c], uint16(c+1))
		}
	}
	return m, nil
}

func identityMatrix(size int) (matrix, error) {
	m, err := newMatrix(size, size)
	if err != nil {
		return nil, err
	}
	for i := range m {
		m[i][i] = 1
	}
	return m, nil
}

func (m matrix) Check() error {
	rows := len(m)
	if rows <= 0 {
		return errInvalidRowSize
	}
	cols := len(m[0])
	if cols <= 0 {
		return errInvalidColSize
	}

	for _, col := range m {
		if len(col) != cols {
			return errColSizeMismatch
		}
	}
	return nil
}

func (m matrix) Mul(other matrix) (matrix, error) {
	if len(m[0]) != len(other) {
		return nil, fmt.Errorf("matrices cannot be multiplied, cols != rows")
	}
	result, err := newMatrix(len(m), len(other[0]))
	if err != nil {
		return nil, err
	}
	for r, row := range result {
		for c := range row {
			var value uint16
			for i := range m[0] {
				value ^= gf.Mul(m[r][i], other[i][c])
			}
			result[r][c] = value
		}
	}
	return result, nil
}

func (m matrix) swapRows(r1, r2 int) error {
	if r1 < 0 || len(m) <= r1 || r2 < 0 || len(m) <= r2 {
		return errInvalidRowSize
	}
	m[r2], m[r1] = m[r1], m[r2]
	return nil
}

func (m matrix) GaussianElimination() error {
	rows := len(m)
	cols := len(m[0])
	for r := 0; r < rows; r++ {
		// We can't work with rows which have 0 on our diagonal slot.
		// Find a row below and swap with it.
		if m[r][r] == 0 {
			for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
				if m[rowBelow][r] != 0 {
					err := m.swapRows(r, rowBelow)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		// If we had to swap but we couldn't, then the matrix is singular.
		if m[r][r] == 0 {
			return errSingular
		}
		// Scale the row to have a 1 in the diagonal.
		if m[r][r] != 1 {
			scale := gf.Div(1, m[r][r])
			for c := 0; c < cols; c++ {
				m[r][c] = gf.Mul(m[r][c], scale)
			}
		}
		// Every row below must have a zero in this column, so subtract
		// multiples of this row.
		for rowBelow := r + 1; rowBelow < rows; rowBelow++ {
			if m[rowBelow][r] != 0 {
				scale := m[rowBelow][r]
				for c := 0; c < cols; c++ {
					m[rowBelow][c] ^= gf.Mul(scale, m[r][c])
				}
			}
		}
	}
	// Clear out everything above the diagonal.
	for d := 0; d < rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m[rowAbove][d] != 0 {
				scale := m[rowAbove][d]
				for c := 0; c < cols; c++ {
					m[rowAbove][c] ^= gf.Mul(scale, m[d][c])
				}
			}
		}
	}
	return nil
}

func (m matrix) Augment(other matrix) (matrix, error) {
	if len(m) != len(other) {
		return nil, errRowSizeMismatch
	}
	newM := matrix(make([][]uint16, len(m)))
	for r := range m {
		newM[r] = make([]uint16, 0, len(m[0])+len(other[0]))
		newM[r] = append(newM[r], m[r]...)
		newM[r] = append(newM[r], other[r]...)
	}
	return newM, nil
}

func (m matrix) AugmentVertical(other matrix) (matrix, error) {
	if len(m[0]) != len(other[0]) {
		return nil, errColSizeMismatch
	}
	newM := matrix(make([][]uint16, 0, len(m)+len(other)))
	for r := range m {
		newM = append(newM, make([]uint16, len(m[r])))
		copy(newM[r], m[r])
	}
	for r := range other {
		newM = append(newM, make([]uint16, len(m[r])))
		copy(newM[r+len(m)], other[r])
	}
	return newM, nil
}
