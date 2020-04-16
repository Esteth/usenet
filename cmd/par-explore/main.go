package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/esteth/usenet/pkg/par2/scanner"
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

	parScanner := scanner.NewScanner(parFile)
	for parScanner.Scan() {
		packet := parScanner.Packet()
		fmt.Printf("packet found: %v\n", string(packet.Type()))

		if mainPacket, ok := packet.(scanner.MainPacket); ok {
			fmt.Printf("Recovery file IDs: %v\n", mainPacket.RecoveryFileIDs)
		}

		if fileDescriptionPacket, ok := packet.(scanner.FileDescriptionPacket); ok {
			fmt.Printf("File Name: %s\n", fileDescriptionPacket.Name)
			fmt.Printf("File ID: %v\n", fileDescriptionPacket.ID)
			fmt.Printf("File MD5: %v\n", fileDescriptionPacket.MD5)
			fmt.Printf("File MD5-16: %v\n", fileDescriptionPacket.MD516)
			fmt.Printf("File Length: %d\n", fileDescriptionPacket.Length)
		}

		if fileSliceChecksumPacket, ok := packet.(scanner.FileSliceChecksumPacket); ok {
			fmt.Printf("File ID: %v\n", fileSliceChecksumPacket.FileID)
			for _, hash := range fileSliceChecksumPacket.SliceHashes {
				fmt.Printf("Slice hash: %v\n", hash)
			}
			for _, checksum := range fileSliceChecksumPacket.SliceCRC32s {
				fmt.Printf("Slice checksum: %v\n", checksum)
			}
		}

		if recoverySlicePacket, ok := packet.(scanner.RecoverySlicePacket); ok {
			fmt.Printf("Exponent: %d\n", recoverySlicePacket.Exponent)
			fmt.Printf("Data: %v\n", recoverySlicePacket.Data)
		}

		if creatorPacket, ok := packet.(scanner.CreatorPacket); ok {
			fmt.Printf("Creator: %s", creatorPacket.Creator)
		}
	}

	if parScanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error reading PAR file: %v\n", err)
		return
	}
}
