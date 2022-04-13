package nzb

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"sync"

	"golang.org/x/net/html/charset"
)

// An Nzb represents the contents of an NZB file.
type Nzb struct {
	Files []File `xml:"file"`
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

// FromFile creates a new Nzb struct by reading an nzb file from disk
func FromFile(filename string) (Nzb, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Nzb{}, fmt.Errorf("could not open '%s': %w", filename, err)
	}
	var nzb Nzb
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&nzb)
	if err != nil {
		return Nzb{}, fmt.Errorf("could not parse '%s' as NZB: %w", filename, err)
	}

	// Sort each file's segments into order
	// TODO: Consider whether this belongs here or in the consumer.
	var wg sync.WaitGroup
	for _, f := range nzb.Files {
		wg.Add(1)
		go func(f File) {
			defer wg.Done()
			sort.Slice(f.Segments, func(i, j int) bool {
				return f.Segments[i].Number < f.Segments[j].Number
			})
		}(f)
	}
	wg.Wait()

	return nzb, nil
}
