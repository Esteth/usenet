package yenc

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestSinglePart(t *testing.T) {
	encodedFile, err := os.Open("testdata/encoded.txt")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded data file: %v", err)
	}

	yencReader, err := NewReader(encodedFile)
	if err != nil {
		t.Fatalf("Could not initialize yenc Reader: %v", err)
	}

	decoded, err := ioutil.ReadAll(yencReader)
	if err != nil {
		t.Fatalf("Failed to read encoded data file: %v", err)
	}

	expectedFile, err := os.Open("testdata/expected.txt")
	if err != nil {
		t.Fatalf("Could not open expected data file: %v", err)
	}
	expected, err := ioutil.ReadAll(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read expected data file: %v", err)
	}

	if !bytes.Equal(decoded, expected) {
		t.Fatalf("Read data %v not equal to expected data %v", decoded, expected)
	}
}

func TestSmallBuffer(t *testing.T) {
	encodedFile, err := os.Open("testdata/encoded.txt")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded data file: %v", err)
	}

	yencReader, err := NewReader(encodedFile)
	if err != nil {
		t.Fatalf("Could not initialize yenc Reader: %v", err)
	}

	expectedFile, err := os.Open("testdata/expected.txt")
	if err != nil {
		t.Fatalf("Could not open expected data file: %v", err)
	}

	for {
		decoded := make([]byte, 5)
		_, err1 := yencReader.Read(decoded)

		expected := make([]byte, 5)
		_, err2 := expectedFile.Read(expected)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return
			} else if err1 == io.EOF || err2 == io.EOF {
				t.Fatalf("One read finished, one is not. actual: %v, expected: %v", err1, err2)
			} else {
				t.Fatalf("Errors reading: %v, %v", err1, err2)
			}
		}

		if !bytes.Equal(decoded, expected) {
			t.Fatalf("Read data %v not equal to expected data %v", decoded, expected)
		}
	}
}
