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

func TestCreatorPacket(t *testing.T) {
	encodedFile, err := os.Open("testdata/sample.mp4.par2")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	scanner := NewScanner(encodedFile)
	var packet Packet = nil
	for scanner.Scan() {
		if scanner.Packet().Type() == creatorPacketType {
			packet = scanner.Packet()
			break
		}
	}

	if scanner.Err() != nil {
		t.Fatalf("Could not read packet: %v", scanner.Err())
	}

	creatorPacket, ok := packet.(CreatorPacket)
	if !ok {
		t.Fatalf("Could not read packet as CreatorPacket")
	}
	if creatorPacket.Creator != "QuickPar 0.9" {
		t.Fatalf("Expected creator QuickPar 0.9 not equal to actual creator %s", creatorPacket.Creator)
	}
}