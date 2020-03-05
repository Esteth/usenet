package scanner

import (
	"os"
	"reflect"
	"testing"
)

func TestPacketTypes(t *testing.T) {
	encodedFile, err := os.Open("testdata/sample.mp4.par2")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	scanner := NewScanner(encodedFile)
	packetTypes := make([]string, 0)
	for scanner.Scan() {
		packetTypes = append(packetTypes, scanner.Packet().Type())
	}

	if scanner.Err() != nil {
		t.Fatalf("Could not read packet: %v", scanner.Err())
	}

	if !reflect.DeepEqual(
		packetTypes,
		[]string{
			fileDescriptionPacketType,
			fileSliceChecksumPacketType,
			mainPacketType,
			creatorPacketType}) {
		t.Errorf("Read packet types %q not equal to expected types", packetTypes)
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

func TestMainPacket(t *testing.T) {
	encodedFile, err := os.Open("testdata/sample.mp4.par2")
	defer encodedFile.Close()
	if err != nil {
		t.Fatalf("Could not open encoded par2 file: %v", err)
	}

	scanner := NewScanner(encodedFile)
	var packet Packet = nil
	for scanner.Scan() {
		if scanner.Packet().Type() == mainPacketType {
			packet = scanner.Packet()
			break
		}
	}

	if scanner.Err() != nil {
		t.Fatalf("Could not read packet: %v", scanner.Err())
	}

	mainPacket, ok := packet.(MainPacket)
	if !ok {
		t.Fatalf("Could not read packet as MainPacket")
	}
	if !reflect.DeepEqual(
		mainPacket.RecoveryFileIDs,
		[][16]byte{
			[16]byte{186, 200, 203, 239, 3, 115, 52, 142, 72, 149, 1, 173, 245, 81, 40, 141}}) {
		t.Fatalf("Expected recovery IDs not equal to actual recovery IDs %v", mainPacket.RecoveryFileIDs)
	}
	if !reflect.DeepEqual(
		mainPacket.NonRecoveryFileIDs,
		[][16]byte{
			[16]byte{80, 65, 82, 50, 0, 80, 75, 84, 76, 0, 0, 0, 0, 0, 0, 0},
			[16]byte{73, 201, 138, 29, 80, 154, 229, 188, 13, 245, 176, 216, 48, 157, 196, 51},
			[16]byte{56, 229, 221, 186, 244, 190, 192, 164, 112, 14, 154, 109, 121, 40, 29, 61},
			[16]byte{80, 65, 82, 32, 50, 46, 48, 0, 67, 114, 101, 97, 116, 111, 114, 0},
			[16]byte{81, 117, 105, 99, 107, 80, 97, 114, 32, 48, 46, 57, 0, 0, 0, 0}}) {
		t.Fatalf("Expected non-recovery IDs not equal to actual non-recovery IDs %v", mainPacket.NonRecoveryFileIDs)
	}
}
