package par2

import (
	"os"
	"testing"
)

func TestPacketCount(t *testing.T) {
	encodedFile, err := os.Open("testdata/sample.mp4.par2")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	scanner := NewScanner(encodedFile)
	packetCount := 0
	for scanner.Scan() {
		packetCount++
	}
	
	if scanner.Err() != nil {
		t.Fatalf("Could not read packet: %v", scanner.Err())
	}

	if packetCount != 4 {
		t.Errorf("Read packet count %v not equal to expected count 4", packetCount)
	}
}