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
		worker := nntp.Worker{
			Address: *address,
			User: *user,
			Password: *password,
			Requests: messageIds,
			Completions: completions,
		}
		go worker.Work()
	}
	for _, segment := range segments {
		messageIds <- segment
	}
	close(messageIds)
	for i := 0; i < len(segments); i++ {
		<-completions
	}
}