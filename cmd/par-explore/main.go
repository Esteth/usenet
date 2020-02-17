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

		if fileDescriptionPacket, ok := packet.(par2.FileDescriptionPacket); ok {
			fmt.Printf("File Name: %s\n", fileDescriptionPacket.Name)
			fmt.Printf("File ID: %v\n", fileDescriptionPacket.ID)
			fmt.Printf("File MD5: %v\n", fileDescriptionPacket.MD5)
			fmt.Printf("File MD5-16: %v\n", fileDescriptionPacket.MD516)
			fmt.Printf("File Length: %d\n", fileDescriptionPacket.Length)
		}

		if fileSliceChecksumPacket, ok := packet.(par2.FileSliceChecksumPacket); ok {
			fmt.Printf("File ID: %v\n", fileSliceChecksumPacket.FileID)
			for _, hash := range fileSliceChecksumPacket.SliceHashes {
				fmt.Printf("Slice hash: %v\n", hash)
			}
			for _, checksum := range fileSliceChecksumPacket.SliceCRC32s {
				fmt.Printf("Slice checksum: %v\n", checksum)
			}
		}
	}

	if parScanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error reading PAR file: %v\n", err)
		return 
	}
}