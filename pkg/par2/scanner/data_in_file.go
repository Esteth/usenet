package scanner

import (
	"fmt"
	"os"
)

type DataInFile struct {
	FilePath   string
	FileOffset uint32
}

func (d *DataInFile) Open() (*os.File, error) {
	return nil, fmt.Errorf("Not Implemented")
}
