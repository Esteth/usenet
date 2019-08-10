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
		t.Errorf("Read data %v not equal to expected data %v", decoded, expected)
	}

	multipart, err := yencReader.Multipart()
	if err != nil {
		t.Fatalf("Could not determine if file is multipart: %v", err)
	}
	if multipart {
		t.Errorf("Incorrectly determined single-part file to be multi-part")
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

func TestFilename(t *testing.T) {
	encodedFile, err := os.Open("testdata/encoded.txt")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded data file: %v", err)
	}

	yencReader, err := NewReader(encodedFile)
	if err != nil {
		t.Fatalf("Could not initialize yenc Reader: %v", err)
	}

	_, err = ioutil.ReadAll(yencReader)
	if err != nil {
		t.Fatalf("Failed to read encoded data file: %v", err)
	}

	filename, err :=  yencReader.Filename()
	if err != nil {
		t.Fatalf("Failed to read filename: %v", err)
	}
	if filename != "testfile.txt" {
		t.Fatalf("Read filename '%s' not equal to expected filename 'testfile.txt'", filename)
	}
}


func TestFilenameBeforeRead(t *testing.T) {
	encodedFile, err := os.Open("testdata/encoded.txt")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded data file: %v", err)
	}

	yencReader, err := NewReader(encodedFile)
	if err != nil {
		t.Fatalf("Could not initialize yenc Reader: %v", err)
	}

	filename, err :=  yencReader.Filename()
	if err != nil {
		t.Fatalf("Failed to read filename: %v", err)
	}
	if filename != "testfile.txt" {
		t.Fatalf("Read filename '%s' not equal to expected filename 'testfile.txt'", filename)
	}
}

func TestMultipart(t *testing.T) {
	encodedFile, err := os.Open("testdata/00000021.ntx")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded data file: %v", err)
	}

	yencReader, err := NewReader(encodedFile)
	if err != nil {
		t.Fatalf("Could not initialize yenc Reader: %v", err)
	}

	multipart, err := yencReader.Multipart()
	if err != nil {
		t.Fatalf("Failed to read multipart information: %v", err)
	}
	if !multipart {
		t.Error("Multipart information not detected")
	}

	offset, err := yencReader.Offset()
	if err != nil {
		t.Fatalf("Failed to read offset: %v", err)
	}
	if offset != 11250 {
		t.Errorf("Offset expected to be 11250, was %d", offset)
	}
}