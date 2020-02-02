package nntp

import "fmt"
import "io"
import "os"

import "github.com/esteth/usenet/pkg/yenc"

type Worker struct {
	Address     string
	User        string
	Password    string
	Requests    <-chan string
	Completions chan<- bool
}

func (w *Worker) Work() {
	conn, err := DialTLS(w.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return
	}

	if w.User != "" && w.Password != "" {
		err = conn.Authenticate(w.User, w.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to Authenticate: %v\n", err)
			return
		}
	}

	for messageID := range w.Requests {
		reader, err := conn.ReadMessage(messageID)
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
		w.Completions <- true
	}
}
