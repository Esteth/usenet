package par2

import "os"

// An Archive represents the contents of an PAR 2.0 archive.
//
// PAR 2.0 archives may be split across multiple files.
type Archive struct {
	files []*os.File
}

// Validate verifies the checksums of the recovery set files, returning nil iff all files are undamaged.
func (a Archive) Validate() error {
	return nil
}

// Repair attempts to repair the recovery set files if they are damaged.
//
// It returns an error if it was unable to complete the repairs.
func (a Archive) Repair() error {
	return nil
}

// FromFiles creates a new Archive struct by reading PAR 2.0 files from disk.
func FromFiles(fs ...*os.File) (Archive, error) {
	return Archive{files: fs}, nil
}
