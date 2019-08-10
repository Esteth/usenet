package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/esteth/usenet/pkg/nntp"
	"github.com/esteth/usenet/pkg/yenc"
	"github.com/esteth/usenet/pkg/nzb"
)

func main() {
	address := flag.String("server", "", "the address of the server to connect to")
	user := flag.String("user", "", "a username to auth to the server")
	password := flag.String("password", "", "a password for auth to the server")
	nzbPath := flag.String("nzb", "", "an NZB file to download the articles from")
	flag.Parse()

	if *address == "" {
		fmt.Fprintf(os.Stderr, "Must specify server address\n")
		return
	}

	if *nzbPath == "" {
		fmt.Fprint(os.Stderr, "Must specify the path to an NZB file\n")
		return
	}

	nzbFile, err := nzb.FromFile(*nzbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse nzb file: %v\n", err)
		return
	}

	conn, err := nntp.DialTLS(*address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return
	}

	if *user != "" && *password != "" {
		err = conn.Authenticate(*user, *password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to Authenticate: %v\n", err)
			return
		}
	}

	for _, s := range nzbFile.Segments {
		reader, err := conn.ReadMessage(s.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return
		}
		yencReader, err := yenc.NewReader(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not create reader: %v\n", err)
			return
		}

		filename, err := yencReader.Filename()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not get filename: %v\n", err)
			return
		}

		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		defer file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not open output file: %v\n", err)
			return
		}

		offset, err := yencReader.Offset()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not read offset from file: %v\n", err)
			return
		}
		file.Seek(offset, 0)

		bytesWritten, err := io.Copy(file, yencReader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not copy data to file: %v\n", err)
			return
		}

		fmt.Printf("Written %d bytes to %s\n", bytesWritten, filename)
	}
}
