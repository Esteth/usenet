package par2

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/esteth/usenet/pkg/par2/reedsolomon"
	"github.com/esteth/usenet/pkg/par2/scanner"
)

// An Archive represents the contents of an PAR 2.0 archive.
//
// PAR 2.0 archives may be split across multiple files.
type Archive struct {
	// baseDirectory is the absolute path to the directory containing the recovery set.
	baseDirectory string
	// parFiles is the set of all PAR 2.0 files contained in the archive.
	// Notably, it does not include the "recovery set", only the par2 files themselves.
	parFiles []*os.File
	// sliceSize is the size of slices across the entire archive.
	sliceSize uint64
	// recoveryFileIDs is a slice containing all file IDs we expect to find in the archive.
	recoveryFileIDs [][16]byte
	// recoverySet is a map from file ID to metadata about that file.
	recoverySet map[[16]byte]*recoveryFile
	// creator is the arbitrary text identifying the creator of the archive.
	creator string
}

// recoveryFile represents a single file from the archive's recovery set.
type recoveryFile struct {
	ID           [16]byte
	MD5          [16]byte
	MD516        [16]byte
	Length       uint64
	Name         string
	SliceMD5s    [][16]byte
	SliceCRC32s  [][4]byte
	RecoveryData []recoveryData
}

// recoveryData represents a single piece of recovery data on disk.
type recoveryData struct {
	exponent   uint32
	filePath   string
	fileOffset uint32
}

// populateDescription copies the file's descriptive attributes from the given scanner.FileDescriptionPacket.
func (rf *recoveryFile) populateDescription(fd scanner.FileDescriptionPacket) {
	rf.ID = fd.ID
	rf.Name = fd.FileName
	rf.Length = fd.FileLength
	rf.MD5 = fd.MD5
	rf.MD516 = fd.MD516
}

// populateDescription copies the file's checksums from the given scanner.FileSliceChecksumPacket.
func (rf *recoveryFile) populateChecksums(fsc scanner.FileSliceChecksumPacket) {
	rf.SliceMD5s = fsc.SliceHashes
	rf.SliceCRC32s = fsc.SliceCRC32s
}

func (rf *recoveryFile) addRecoveryData(rs scanner.RecoverySlicePacket) {
	rf.RecoveryData = append(rf.RecoveryData, recoveryData{
		exponent:   rs.Exponent,
		filePath:   rs.RecoveryDataFilePath,
		fileOffset: rs.RecoveryDataFileOffset,
	})
}

// Validate verifies the checksums of the recovery set files, returning nil iff all files are undamaged.
func (a *Archive) Validate() error {
	for _, id := range a.recoveryFileIDs {
		recoveryFile, exists := a.recoverySet[id]
		if !exists {
			return fmt.Errorf("Could not find checksum data for file ID %v", id)
		}
		f, err := os.Open(filepath.Join(a.baseDirectory, recoveryFile.Name))
		if err != nil {
			return fmt.Errorf("Could not find expected file to validate at %s", recoveryFile.Name)
		}
		defer f.Close()

		buf := make([]byte, a.sliceSize)
		for i, expectedChecksum := range recoveryFile.SliceMD5s {
			bytesRead, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return fmt.Errorf("Could not read from recovery file %s: %w", recoveryFile.Name, err)
				}
			}
			// The specification says that "empty" bytes should be zeroed.
			for i := uint64(bytesRead); i < a.sliceSize; i++ {
				buf[i] = 0
			}

			actualChecksum := md5.Sum(buf)
			if actualChecksum != expectedChecksum {
				return fmt.Errorf("Checksum failed for file %s in slice %d", recoveryFile.Name, i)
			}
		}
	}
	return nil
}

// Repair attempts to repair the recovery set files if they are damaged.
//
// It returns an error if it was unable to complete the repairs.
func (a *Archive) Repair() error {
	data := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	checksums := []uint16{11, 60570, 57778}
	identity, err := reedsolomon.IdentityMatrix(len(data))
	if err != nil {
		return fmt.Errorf("could not create identity matrix: %w", err)
	}
	vandermonde, err := reedsolomon.NewVandermondePar2Matrix(3, len(data))
	if err != nil {
		return fmt.Errorf("could not create vandermonde matrix: %w", err)
	}
	parityIdentity, err := identity.AugmentVertical(vandermonde)
	if err != nil {
		return fmt.Errorf("could not stack identity and vandermonde matrix: %w", err)
	}
	sourceColumn, err := reedsolomon.NewMatrixColumn(append(data, checksums...))
	if err != nil {
		return fmt.Errorf("could not create column with data and checksums: %w", err)
	}
	solve, err := parityIdentity.Augment(sourceColumn)
	if err != nil {
		return fmt.Errorf("could not create problem matrix: %w", err)
	}

	// Delete some rows to pretend we lost some data
	solve = append(solve[:4], solve[7:]...)

	err = solve.GaussianElimination()
	if err != nil {
		return fmt.Errorf("could not solve problem matrix: %w", err)
	}

	recoveredData := make([]uint16, len(data))
	for r, row := range solve {
		recoveredData[r] = row[len(row)-1]
	}

	// TODO: Do something with the recoveredData.

	return nil
}

// FromFiles creates a new Archive struct by reading PAR 2.0 files from disk.
func FromFiles(baseDirectory string, fs ...*os.File) (Archive, error) {
	baseDirectory, err := filepath.Abs(baseDirectory)
	if err != nil {
		return Archive{}, fmt.Errorf("Could not convert base directory %s to absolute path: %w", baseDirectory, err)
	}
	var sliceSize uint64 = 0
	recoveryFileIDs := make([][16]byte, 0)
	recoverySet := make(map[[16]byte]*recoveryFile)
	creatorText := ""

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
			if rsp, ok := packet.(scanner.RecoverySlicePacket); ok {
				if rf, exists := recoverySet[rsp.FileID]; !exists {
					rf.addRecoveryData(rsp)
				}
			}
			if creatorPacket, ok := packet.(scanner.CreatorPacket); ok {
				creatorText = creatorPacket.Creator
			}
		}
	}
	return Archive{
		baseDirectory:   baseDirectory,
		parFiles:        fs,
		sliceSize:       sliceSize,
		recoveryFileIDs: recoveryFileIDs,
		recoverySet:     recoverySet,
		creator:         creatorText,
	}, nil
}
