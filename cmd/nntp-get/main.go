package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/esteth/usenet/pkg/nntp"
)

func main() {
	address := flag.String("server", "", "the address of the server to connect to")
	messageID := flag.String("message", "", "a message ID to download from the server")
	user := flag.String("user", "", "a username to auth to the server")
	password := flag.String("password", "", "a password for auth to the server")
	flag.Parse()

	if *address == "" {
		fmt.Fprintf(os.Stderr, "Must specify server address")
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

	reader, err := conn.ReadMessage(*messageID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return
	}
	fmt.Printf("%v\n", string(bytes))
}
