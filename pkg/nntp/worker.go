package nntp

import "fmt"
import "os"

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
		bytesWritten, err := conn.ReadMessageToFile(messageID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message to file: %v\n", err)
		}

		fmt.Printf("Written %d bytes\n", bytesWritten)
		w.Completions <- true
	}
}
