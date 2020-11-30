package par2

import (
	"bytes"
	"io"
	"os"
	"path"
	"testing"
)

func copyFile(t *testing.T, src string, dst string) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		t.Fatalf("Could not stat copy src %s: %v", src, err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		t.Fatalf("%s is not a regular file", src)
	}

	destFileStat, err := os.Stat(dst)
	if err != nil {
		t.Fatalf("Could not stat copy dst %s: %v", dst, err)
	}

	if destFileStat.IsDir() {
		dst = path.Join(dst, path.Base(src))
	}

	source, err := os.Open(src)
	if err != nil {
		t.Fatalf("Could not open copy src %s: %v", src, err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		t.Fatalf("Could not create copy dst %s: %v", dst, err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		t.Fatalf("Could not copy src %s to dst %s: %v", src, dst, err)
	}
}

func areFileContentsEqual(t *testing.T, file1 string, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		t.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		t.Fatal(err)
	}
	defer f2.Close()

	for {
		b1 := make([]byte, 1024)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, 1024)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				t.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func TestValidateValidArchive(t *testing.T) {
	f, err := os.Open("testdata/sample.mp4.par2")
	defer f.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	archive, err := FromFiles("testdata", f)
	if err != nil {
		t.Fatalf("Could not create Archive from file: %v", err)
	}

	if err = archive.Validate(); err != nil {
		t.Fatalf("Intact archive did not validate as expected: %v", err)
	}
}

func TestValidateBrokenFiles(t *testing.T) {
	f, err := os.Open("testdata/sample.broken.mp4.par2")
	defer f.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	archive, err := FromFiles("testdata", f)
	if err != nil {
		t.Fatalf("Could not create Archive from file: %v", err)
	}

	if err = archive.Validate(); err == nil {
		t.Fatalf("Broken archive unexpectedly validated.")
	}
}

func TestRepairValidArchive(t *testing.T) {
	f, err := os.Open("testdata/sample.mp4.par2")
	defer f.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	archive, err := FromFiles("testdata", f)
	if err != nil {
		t.Fatalf("Could not create Archive from file: %v", err)
	}

	if err = archive.Repair(); err != nil {
		t.Fatalf("Intact archive threw error when asked to repair: %v", err)
	}
}

func TestRepairBrokenFiles(t *testing.T) {
	tempDir := t.TempDir()
	copyFile(t, "testdata/sample.broken.mp4", tempDir)
	copyFile(t, "testdata/sample.broken.mp4.par2", tempDir)
	copyFile(t, "testdata/sample.broken.mp4.vol0+1.PAR2", tempDir)
	copyFile(t, "testdata/sample.broken.mp4.vol1+2.PAR2", tempDir)
	copyFile(t, "testdata/sample.broken.mp4.vol3+2.PAR2", tempDir)

	f, err := os.Open(path.Join(tempDir, "sample.broken.mp4.par2"))
	defer f.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	archive, err := FromFiles("testdata", f)
	if err != nil {
		t.Fatalf("Could not create Archive from file: %v", err)
	}

	if err = archive.Repair(); err != nil {
		t.Fatalf("Intact archive threw error when asked to repair: %v", err)
	}

	if !areFileContentsEqual(t, path.Join(tempDir, "sample.broken.mp4"), "testdata/sample.mp4") {
		t.Errorf("Repaired file contents not same as reference file")
	}
}
