package nzb

import (
	"sort"
	"os"
	"encoding/xml"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
)

type nzb struct {
	File File `xml:"file"`
}

// A File is metadata regarding where to find data in usenet for a particular file
type File struct {
	Subject  string    `xml:"subject,attr"`
	Segments []Segment `xml:"segments>segment"`
}

// A Segment is a pointer to an email message containing binary data
type Segment struct {
	Number int    `xml:"number,attr"`
	Bytes  int    `xml:"bytes,attr"`
	ID     string `xml:",innerxml"`
}

// FromFile creates a new File struct by reading an nzb file from disk
func FromFile(filename string) (File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return File{}, errors.Wrapf(err, "could not open '%s'", filename)
	}
	var nzb nzb
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&nzb)
	if err != nil {
		return File{}, errors.Wrapf(err, "could not parse '%s' as NZB", filename)
	}

	sort.Slice(nzb.File.Segments, func(i, j int) bool {
		return nzb.File.Segments[i].Number  < nzb.File.Segments[j].Number
	})

	return nzb.File, nil
}