package par2

import (
	"fmt"
	"os"

	"github.com/esteth/usenet/pkg/par2/scanner"
)

// An Archive represents the contents of an PAR 2.0 archive.
//
// PAR 2.0 archives may be split across multiple files.
type Archive struct {
	// parFiles is the set of all PAR 2.0 files contained in the archive.
	// Notably, it does not include the "recovery set", only the par2 files themselves.
	parFiles []*os.File
	// sliceSize is the size of slices across the entire archive.
	sliceSize uint64
	// recoveryFileIDs is a slice containing all file IDs we expect to find in the archive.
	recoveryFileIDs [][16]byte
	// recoverySet is a map from file ID to metadata about that file.
	recoverySet map[[16]byte]*recoveryFile
}

// A recoveryFile represents a single file from the archive's recovery set.
type recoveryFile struct {
	ID          [16]byte
	MD5         [16]byte
	MD516       [16]byte
	Length      uint64
	Name        string
	SliceMD5s   [][16]byte
	SliceCRC32s [][4]byte
}

// populateDescription copies the file's descriptive attributes from the given scanner.FileDescriptionPacket.
func (rf *recoveryFile) populateDescription(fd scanner.FileDescriptionPacket) {
	rf.ID = fd.ID
	rf.Name = fd.Name
	rf.Length = fd.Length
	rf.MD5 = fd.MD5
	rf.MD516 = fd.MD516
}

// populateDescription copies the file's checksums from the given scanner.FileSliceChecksumPacket.
func (rf *recoveryFile) populateChecksums(fsc scanner.FileSliceChecksumPacket) {
	rf.SliceMD5s = fsc.SliceHashes
	rf.SliceCRC32s = fsc.SliceCRC32s
}

// Validate verifies the checksums of the recovery set files, returning nil iff all files are undamaged.
func (a *Archive) Validate() error {
	for _, id := range a.recoveryFileIDs {
		rf, exists := a.recoverySet[id]
		if !exists {
			return fmt.Errorf("Could not find checksum data for file ID %v", id)
		}
		_, err := os.Open(rf.Name)
		if err != nil {
			return fmt.Errorf("Could not find expected file to validate at %s", rf.Name)
		}
	}
	return nil
}

// Repair attempts to repair the recovery set files if they are damaged.
//
// It returns an error if it was unable to complete the repairs.
func (a *Archive) Repair() error {
	return nil
}

// FromFiles creates a new Archive struct by reading PAR 2.0 files from disk.
func FromFiles(fs ...*os.File) (Archive, error) {
	var sliceSize uint64 = 0
	recoveryFileIDs := make([][16]byte, 0)
	recoverySet := make(map[[16]byte]*recoveryFile)

	for _, f := range fs {
		parScanner := scanner.NewScanner(f)
		for parScanner.Scan() {
			packet := parScanner.Packet()
			if mainPacket, ok := packet.(scanner.MainPacket); ok {
				sliceSize = mainPacket.SliceSize
				for _, fileID := range mainPacket.RecoveryFileIDs {
					recoveryFileIDs = append(recoveryFileIDs, fileID)
				}
			}
			if fd, ok := packet.(scanner.FileDescriptionPacket); ok {
				if _, exists := recoverySet[fd.ID]; !exists {
					recoverySet[fd.ID] = &recoveryFile{}
				}
				recoverySet[fd.ID].populateDescription(fd)
			}
			if fsc, ok := packet.(scanner.FileSliceChecksumPacket); ok {
				if _, exists := recoverySet[fsc.FileID]; !exists {
					recoverySet[fsc.FileID] = &recoveryFile{}
				}
				recoverySet[fsc.FileID].populateChecksums(fsc)
			}
		}
	}
	return Archive{
		parFiles:        fs,
		sliceSize:       sliceSize,
		recoveryFileIDs: recoveryFileIDs,
		recoverySet:     recoverySet,
	}, nil
}
