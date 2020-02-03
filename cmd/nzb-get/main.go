package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/esteth/usenet/pkg/nntp"
	"github.com/esteth/usenet/pkg/nzb"
)

func main() {
	address := flag.String("server", "", "the address of the server to connect to")
	user := flag.String("user", "", "a username to auth to the server")
	password := flag.String("password", "", "a password for auth to the server")
	nzbPath := flag.String("nzb", "", "an NZB file to download the articles from")
	maxConnections := flag.Int("connections", 1, "the number of simultaneous connections to use during the download")
	flag.Parse()

	if *address == "" {
		fmt.Fprintf(os.Stderr, "Must specify server address\n")
		return
	}

	if *nzbPath == "" {
		fmt.Fprint(os.Stderr, "Must specify the path to an NZB file\n")
		return
	}

	nzb, err := nzb.FromFile(*nzbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse nzb file: %v\n", err)
		return
	}
	
	segments := make([]string, 0)
	for _, file := range nzb.Files {
		for _, segment := range file.Segments {
			segments = append(segments, segment.ID)
		}
	}

	messageIds := make(chan string, len(segments))
	completions := make(chan bool, len(segments))
	for c := 0; c < *maxConnections; c++ {
		go worker(*address, *user, *password, messageIds, completions)
	}
	for _, segment := range segments {
		messageIds <- segment
	}
	close(messageIds)
	for i := 0; i < len(segments); i++ {
		<-completions
	}
}

func worker(address string, user string, password string, requests <-chan string, completions chan<- bool) {
	conn, err := nntp.DialTLS(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return
	}

	if user != "" && password != "" {
		err = conn.Authenticate(user, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to Authenticate: %v\n", err)
			return
		}
	}

	for messageID := range requests {
		bytesWritten, err := conn.ReadMessageToFile(messageID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message to file: %v\n", err)
		}

		fmt.Printf("Written %d bytes\n", bytesWritten)
		completions <- true
	}
}