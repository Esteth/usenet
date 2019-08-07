package yenc

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type header struct {
	lineLength int
}

// A Reader is an io.Reader that can be read to retrieve
// yenc decoded data from a reader containing yenc encoded
// data
type Reader struct {
	s              bufio.Scanner
	err            error
	foundHeader    bool
	header         header
	overflowBuffer []byte
}

// NewReader creates a new reader reading the given reader.
func NewReader(r io.Reader) (*Reader, error) {
	z := new(Reader)
	if err := z.Reset(r); err != nil {
		return nil, err
	}
	return z, nil
}

// Reset discards the Reader z's state and makes it equivalent to the
// result of it's original state from NewReader, but reading from r instead.
// This permits reusing a reader rather than allocating a new one.
func (z *Reader) Reset(r io.Reader) error {
	*z = Reader{
		s: *bufio.NewScanner(r),
	}
	return nil
}

// Read implements io.Reader, reading encoded bytes from its underlying Reader.
func (z *Reader) Read(buf []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}
	if z.overflowBuffer != nil {
		overflowBuffer := z.overflowBuffer
		z.overflowBuffer = nil
		n, err = z.readLine(buf, overflowBuffer)
		// Need to continue to read if we have exhausted the overflow but not the input
		if (n == len(buf)) {
			return	
		}
	}
	for {
		scanSucceeded := z.s.Scan()
		if !scanSucceeded {
			if z.s.Err() != nil {
				return 0, z.s.Err()
			}
			// We found the end of the file
			z.err = io.EOF
			if (n > 0) {
				// Most Reader clients expect 0, EOF, so save the EOF for the next Read call
				return n, nil
			}
			return n, io.EOF
		}
		if !z.foundHeader {
			// Ignore all text until we find the yEnc begin header
			if strings.HasPrefix(z.s.Text(), "=ybegin") {
				z.foundHeader = true
				z.header, err = parseHeader(z.s.Text())
				if err != nil {
					return n, errors.Wrap(err, "Failed to parse ybegin header")
				}
			}
		} else {
			if strings.HasPrefix(z.s.Text(), "=yend") {
				z.foundHeader = false
				// TODO: Parse Footer and perform CRC/length checks
				// _, err = parseFooter(z.s.Text())
				// if err != nil {
				// 	return 0, errors.Wrap(err, "Failed to parse yend header")
				// }
				continue
			} else {
				break
			}
		}
	}

	encodedLine := z.s.Bytes()
	n, err = z.readLine(buf[n:], encodedLine)
	return
}

func (z *Reader) readLine(output []byte, input []byte) (n int, err error) {
	escapeNext := false
	for i, b := range input {
		if b == '=' && !escapeNext {
			// '=' is the escape character in yEnc. It shouldn't appear in the
			// output, only modify the next character.
			escapeNext = true
			continue
		}
		if escapeNext {
			// Escaped characters must be shifted an extra 64 to avoid critical
			// control characters appearing in encoded text.
			// TODO: Log error if attempting to escape an unnessecary character
			b = b - 64
			escapeNext = false
		}
		// Most of yEnc encoding just adds 42 to each byte. Reverse that.
		b -= 42
		if n < len(output) {
			output[n] = b
			n++
		} else {
			z.overflowBuffer = input[i:]
			return
		}
	}
	return
}

func parseHeader(headerText string) (h header, err error) {
	fields := strings.Fields(headerText)[1:]
	for _, field := range fields {
		re := regexp.MustCompile(`(\w+)=(\w+)`)
		result := re.FindSubmatch([]byte(field))
		if result == nil || len(result) == 0 {
			return header{}, errors.Wrapf(err, "Failed to parse header field \"%v\"", field)
		}
		key := string(result[0][1])
		val := string(result[0][2])

		if key == "line" {
			h.lineLength, err = strconv.Atoi(val)
		}
	}

	return
}
