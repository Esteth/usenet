package par2

import (
	"os"
	"testing"
)

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
	f, err := os.Open("testdata/sample.broken.mp4.par2")
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

	// TODO: Test that the files were actually fixed
	//   Will require unpacking to a new directory instead of fixing in-place.
	//   Use go's testing package temporary directory, and setup/teardown to reset.
}
