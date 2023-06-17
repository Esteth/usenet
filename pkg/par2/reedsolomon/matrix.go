package reedsolomon

import (
	"errors"
	"fmt"

	"github.com/esteth/usenet/pkg/par2/gf"
)

// Matrix type heavily inspired by github.com/klauspost/reedsolomon
type matrix struct {
	data []uint16
	rows int
	cols int
}

func (m matrix) cell(row, col int) uint16 {
	return m.data[row*m.cols+col]
}

func (m matrix) row(row int) []uint16 {
	return m.data[row*m.cols : (row+1)*m.cols]
}

var invalidMatrix = matrix{
	data: nil,
	rows: 0,
}

var errInvalidRowSize = errors.New("invalid row size")
var errInvalidColSize = errors.New("invalid col size")
var errRowSizeMismatch = errors.New("row size is not the same for both matrices")
var errColSizeMismatch = errors.New("column size is not the same for all rows")
var errSingular = errors.New("cannot solve a singular matrix")

// NewMatrix creates a new matrix of the given size filled with zeroes.
func NewMatrix(rows, cols int) matrix {
	if rows <= 0 {
		panic(errInvalidRowSize)
	}
	if cols <= 0 {
		panic(errInvalidColSize)
	}
	m := matrix{
		data: make([]uint16, rows*cols),
		rows: rows,
		cols: cols,
	}
	return m
}

// NewMatrix creates a new matrix backed by the given slice.
func NewMatrixData(data [][]uint16) (matrix, error) {
	m := NewMatrix(len(data), len(data[0]))
	offset := 0
	for i := range data {
		width := len(data[i])
		copy(m.data[offset:offset+width], data[i])
		offset += width
	}
	return m, nil
}

// NewMatrixColumn creates a new single-column matrix with a copy of the given
// data.
func NewMatrixColumn(data []uint16) (matrix, error) {
	if len(data) <= 0 {
		return invalidMatrix, errInvalidColSize
	}
	m := NewMatrix(len(data), 1)
	copy(m.data, data)
	return m, nil
}

// NewVandermondePar2Matrix creates a new Vandermonde matrix using the Par2
// specification's custom rules for generating Vandermonde matrices.
func NewVandermondePar2Matrix(rows, cols int) (matrix, error) {
	m := NewMatrix(rows, cols)
	var constantPool constantPool = newConstantPool()
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if c == 0 || r == 0 {
				m.data[r*cols+c] = 1
				continue
			}
			if r == 1 {
				m.data[r*cols+c] = constantPool.Next()
				continue
			}
			m.data[r*cols+c] = gf.Mul(m.cell(r-1, c), m.cell(1, c))
		}
	}
	return m, nil
}

// IdentityMatrix creates a square identity matrix of the given size.
func IdentityMatrix(size int) (matrix, error) {
	m := NewMatrix(size, size)
	for i := 0; i < size*size; i += (size + 1) {
		m.data[i] = 1
	}
	return m, nil
}

// Mul performs matrix multiplication with other.
// It returns the result, modifying neither matrix in-place.
func (m matrix) Mul(other matrix) (matrix, error) {
	if m.cols != other.rows {
		return invalidMatrix, fmt.Errorf("matrices cannot be multiplied, cols != rows")
	}
	result := NewMatrix(m.rows, len(other.data)/other.rows)
	for r := 0; r < result.rows; r++ {
		for c := 0; c < result.cols; c++ {
			var value uint16
			for i := 0; i < m.cols; i++ {
				value ^= gf.Mul(m.cell(r, i), other.cell(i, c))
			}
			result.data[r*result.cols+c] = value
		}
	}
	return result, nil
}

func (m matrix) swapRows(r1, r2 int, buf []uint16) error {
	if r1 < 0 || m.rows <= r1 || r2 < 0 || m.rows <= r2 {
		return errInvalidRowSize
	}
	copy(buf, m.row(r2))
	copy(m.row(r2), m.row(r1))
	copy(m.row(r1), buf)
	return nil
}

// GaussianElimination performs a gaussian elimination in-place on m.
// After gaussian elimination, each row should have a single column equal to 1
// and the final column represents the value of that variable in the system
// of linear equations.
func (m matrix) GaussianElimination() error {
	buf := make([]uint16, m.cols)
	for r := 0; r < m.rows; r++ {
		// We can't work with rows which have 0 on our diagonal slot.
		// Find a row below and swap with it.
		if m.cell(r, r) == 0 {
			for rowBelow := r + 1; rowBelow < m.rows; rowBelow++ {
				if m.cell(rowBelow, r) != 0 {
					err := m.swapRows(r, rowBelow, buf)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		// If we had to swap but we couldn't, then the matrix is singular.
		if m.cell(r, r) == 0 {
			return errSingular
		}
		// Scale the row to have a 1 in the diagonal.
		if m.cell(r, r) != 1 {
			scale := gf.Div(1, m.cell(r, r))
			for c := 0; c < m.cols; c++ {
				m.data[r*m.cols+c] = gf.Mul(m.cell(r, c), scale)
			}
		}
		// Every row below must have a zero in this column, so subtract
		// multiples of this row.
		for rowBelow := r + 1; rowBelow < m.rows; rowBelow++ {
			if m.cell(rowBelow, r) != 0 {
				scale := m.cell(rowBelow, r)
				for c := 0; c < m.cols; c++ {
					m.data[rowBelow*m.cols+c] ^= gf.Mul(scale, m.cell(r, c))
				}
			}
		}
	}
	// Clear out everything above the diagonal.
	for d := 0; d < m.rows; d++ {
		for rowAbove := 0; rowAbove < d; rowAbove++ {
			if m.cell(rowAbove, d) != 0 {
				scale := m.cell(rowAbove, d)
				for c := 0; c < m.cols; c++ {
					m.data[rowAbove*m.cols+c] ^= gf.Mul(scale, m.cell(d, c))
				}
			}
		}
	}
	return nil
}

// Augment returns a new matrix by putting other to the right of this matrix.
// Both matrices MUST have the same number of rows.
func (m matrix) Augment(other matrix) (matrix, error) {
	if m.rows != other.rows {
		return invalidMatrix, errRowSizeMismatch
	}
	newM := NewMatrix(m.rows, m.cols+other.cols)
	for r := 0; r < m.rows; r++ {
		copy(newM.row(r), m.row(r))
		copy(newM.row(r)[m.cols:], other.row(r))
	}
	return newM, nil
}

// Augment returns a new matrix by putting other to the bottom of this matrix.
// Both matrices MUST have the same number of columns.
func (m matrix) AugmentVertical(other matrix) (matrix, error) {
	if m.cols != other.cols {
		return invalidMatrix, errColSizeMismatch
	}
	newM := NewMatrix(m.rows+other.rows, m.cols)

	copy(newM.data, m.data)
	copy(newM.data[len(m.data):], other.data)
	return newM, nil
}
