package yenc

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"hash"
	"hash/crc32"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type header struct {
	lineLength int
	name       string
	size 	   int64
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
	hash           hash.Hash
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
		s:    *bufio.NewScanner(r),
		hash: crc32.New(crc32.MakeTable(crc32.IEEE)),
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
		if n == len(buf) {
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
			if n > 0 {
				// Most Reader clients expect 0, EOF, so save the EOF for the next Read call
				return n, nil
			}
			return n, io.EOF
		}
		if !z.foundHeader {
			// Ignore all text until we find the yEnc begin header
			if strings.HasPrefix(z.s.Text(), "=ybegin") {
				z.foundHeader = true
				z.header, err = parseBegin(z.s.Text())
				if err != nil {
					return n, errors.Wrap(err, "Failed to parse ybegin header")
				}
			}
		} else {
			if strings.HasPrefix(z.s.Text(), "=yend") {
				z.foundHeader = false
				err = z.validateEnd(z.s.Text())
				if err != nil {
					return n, errors.Wrap(err, "Failed to validate footer")
				}
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

// Filename returns the filename specified in the ybegin header.
//
// If no header has been read, it returns an error.
func (z *Reader) Filename() (string, error) {
	if z.header.name == "" {
		return "", errors.New("Cannot determine filename until ybegin header has been read")
	}
	return z.header.name, nil
}

// readLine reads a single line of input data from intput into output.
// It returns the number of bytes written to output and and error.
//
// Note: readLine should only be called when the Reader is positioned between ybegin and yend. 
func (z *Reader) readLine(output []byte, input []byte) (n int, err error) {
	// Before we return, add all the bytes we wrote to the ongoing CRC32
	defer func() { z.hash.Write(output[:n]) }()

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
			// If we've run out of space in the output buffer, save the overflow in the Reader
			z.overflowBuffer = input[i:]
			return
		}
	}
	return
}

// parseBegin parses a "=ybegin" header line, returning it and an error
func parseBegin(beginLine string) (h header, err error) {
	fields, err := parseHeader(beginLine)
	if err != nil {
		return header{}, errors.Wrapf(err, "Failed to parse ybegin line: %v", beginLine)
	}

	h.lineLength, err = strconv.Atoi(fields["line"])
	if err != nil {
		return header{}, errors.Wrapf(err, "could not convert 'line' to int: %s", fields["line"])
	}

	if size, ok := fields["size"]; ok {
		h.size, err = strconv.ParseInt(size, 10, 0)
		if err != nil {
			return header{}, errors.Wrapf(err, "could not convert 'size' to int: %s", fields["size"])
		}
	} else {
		return header{}, errors.New("ybegin header does not contain size field")
	}

	if name, ok := fields["name"]; ok {
		h.name = name	
	} else {
		return header{}, errors.New("ybegin header does not contain name field")
	}
	
	return
}

// validateEnd validates a "=yend" header line, returning an error if it does not validate
func (z *Reader) validateEnd(endLine string) error {
	fields, err := parseHeader(endLine)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse yend line: %v", endLine)
	}

	// Only conduct a CRC32 check if the checksum is present in the footer
	if expectedString, ok := fields["crc32"]; ok {
		expected, err := hex.DecodeString(expectedString)
		if err != nil {
			return errors.Wrapf(err, "CRC32 Check Failure. Could not parse checksum %s", expectedString)
		}
		actual := z.hash.Sum(nil)

		if !bytes.Equal(expected, actual) {
			return errors.Errorf("CRC32 Check failure. Expected %v, Actual %v", expected, actual)
		}
	}

	if sizeString, ok := fields["size"]; ok {
		size, err := strconv.ParseInt(sizeString, 10, 0)
		if err != nil {
			return errors.Wrap(err, "size validation failure: Could not parse size in footer")
		}
		if size != z.header.size {
			return errors.New("header and foter do not agree on size. Could not validate")
		}
	} else {
		return errors.New("no size found in footer. Could not validate")
	}
	return nil
}

// parseHeader parses a yenc header line, returning a map of the fields contained in it and an error.
func parseHeader(line string) (m map[string]string, err error) {
	fields := strings.Fields(line)[1:]
	m = make(map[string]string, len(fields))
	for _, field := range fields {
		re := regexp.MustCompile(`(\w+)=([^\s]+)`)
		result := re.FindSubmatch([]byte(field))
		if result == nil || len(result) == 0 {
			return nil, errors.Wrapf(err, "Failed to parse header field \"%v\"", field)
		}
		key := string(result[1])
		val := string(result[2])
		m[key] = val
	}
	return
}
