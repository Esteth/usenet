package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/esteth/usenet/pkg/par2"
)

func main() {
	parPath := flag.String("par", "", "an PAR2 file to read")
	flag.Parse()

	if *parPath == "" {
		fmt.Fprint(os.Stderr, "Must specify the path to an PAR file\n")
		return
	}

	parFile, err := os.Open(*parPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open PAR file: %v\n", err)
		return
	}
	
	parScanner := par2.NewScanner(parFile)
	for parScanner.Scan() {
		packet := parScanner.Packet()
		fmt.Printf("packet found: %v\n", string(packet.Type()))

		if mainPacket, ok := packet.(par2.MainPacket); ok {
			fmt.Printf("Recovery file IDs: %v\n", mainPacket.RecoveryFileIDs)
		}
	}

	if parScanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error reading PAR file: %v\n", err)
		return 
	}
}