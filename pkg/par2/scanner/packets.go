package scanner

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const mainPacketType = "PAR 2.0\000Main\000\000\000\000"

// MainPacket represents a Par 2.0 Main Packet.
type MainPacket struct {
	SliceSize          uint64
	RecoveryFileIDs    [][16]byte
	NonRecoveryFileIDs [][16]byte
}

// Type implements interface Packet to return the type of the Par 2.0 Main Packet.
func (p MainPacket) Type() string {
	return mainPacketType
}

// NewMainPacket creates and initializes a new MainPacket struct from the given binary packet data.
func NewMainPacket(data []byte) (MainPacket, error) {
	typ := string(data[48:64])
	if typ != mainPacketType {
		return MainPacket{}, fmt.Errorf("Main packet type not as expected. Was %s", typ)
	}
	mainPacketData := data[64:]

	sliceSize := binary.LittleEndian.Uint64(mainPacketData[:8])

	numRecoveryFiles := binary.LittleEndian.Uint32(mainPacketData[8:12])
	recoveryFileIDs := make([][16]byte, numRecoveryFiles)
	for i := uint32(0); i < numRecoveryFiles; i++ {
		copy(recoveryFileIDs[i][:], mainPacketData[12+16*i:28+16*i])
	}

	numNonRecoveryFiles := (binary.LittleEndian.Uint64(data[8:16]) - 28 + 16*uint64(numRecoveryFiles)) / 16
	nonRecoveryFilesBaseIndex := uint64(12 + 16*numRecoveryFiles)
	nonRecoveryFileIDs := make([][16]byte, numNonRecoveryFiles)
	for i := uint64(0); i < numNonRecoveryFiles; i++ {
		copy(nonRecoveryFileIDs[i][:], mainPacketData[nonRecoveryFilesBaseIndex+16*i:nonRecoveryFilesBaseIndex+16+16*i])
	}

	return MainPacket{
		SliceSize:          sliceSize,
		RecoveryFileIDs:    recoveryFileIDs,
		NonRecoveryFileIDs: nonRecoveryFileIDs,
	}, nil
}

const fileDescriptionPacketType = "PAR 2.0\000FileDesc"

// FileDescriptionPacket represents a Par 2.0 File Description Packet.
type FileDescriptionPacket struct {
	ID     [16]byte
	MD5    [16]byte
	MD516  [16]byte
	Length uint64
	Name   string
}

// Type implements interface Packet to return the type of the Par 2.0 File Description Packet.
func (p FileDescriptionPacket) Type() string {
	return fileDescriptionPacketType
}

// NewFileDescriptionPacket creates and initializes a new FileDescriptionPacket struct from the
// given binary packet data.
func NewFileDescriptionPacket(data []byte) (FileDescriptionPacket, error) {
	typ := string(data[48:64])
	if typ != fileDescriptionPacketType {
		return FileDescriptionPacket{}, fmt.Errorf("File Description packet type not as expected. Was %s", typ)
	}
	packetData := data[64:]

	packet := FileDescriptionPacket{}
	copy(packet.ID[:], packetData[0:16])
	copy(packet.MD5[:], packetData[16:32])
	copy(packet.MD516[:], packetData[32:48])
	packet.Length = binary.LittleEndian.Uint64(packetData[48:56])
	packet.Name = strings.TrimRight(string(packetData[56:]), "\000")

	return packet, nil
}

const fileSliceChecksumPacketType = "PAR 2.0\000IFSC\000\000\000\000"

// FileSliceChecksumPacket represents a Par 2.0 File Slice Checksum Packet.
type FileSliceChecksumPacket struct {
	FileID      [16]byte
	SliceHashes [][16]byte
	SliceCRC32s [][4]byte
}

// Type implements interface Packet to return the type of the Par 2.0 Input File Slice Checksum Packet.
func (p FileSliceChecksumPacket) Type() string {
	return fileSliceChecksumPacketType
}

// NewFileSliceChecksumPacket creates and initializes a new FileSliceChecksumPacket struct from the
// given binary packet data.
func NewFileSliceChecksumPacket(data []byte) (FileSliceChecksumPacket, error) {
	typ := string(data[48:64])
	if typ != fileSliceChecksumPacketType {
		return FileSliceChecksumPacket{}, fmt.Errorf("File Slice Checksum packet type not as expected. Was %s", typ)
	}
	packetData := data[64:]

	packet := FileSliceChecksumPacket{}
	copy(packet.FileID[:], packetData[0:16])

	numSlices := (binary.LittleEndian.Uint64(data[8:16]) - 16 - 64) / 20
	packet.SliceHashes = make([][16]byte, numSlices)
	packet.SliceCRC32s = make([][4]byte, numSlices)

	for i := uint64(0); i < numSlices; i++ {
		copy(packet.SliceHashes[i][:], packetData[16+i*20:32+i*20])
		copy(packet.SliceCRC32s[i][:], packetData[32+i*20:36+i*20])
	}

	return packet, nil
}

const recoverySlicePacketType = "PAR 2.0\000RecvSlic"

// RecoverySlicePacket represents a Par 2.0 Recovery Slice packet.
type RecoverySlicePacket struct {
	Exponent uint32
	Data     []byte
}

// Type implements interface Packet to return the type of the Par 2.0 Recovery Slice Packet.
func (p RecoverySlicePacket) Type() string {
	return recoverySlicePacketType
}

// NewRecoverySlicePacket creates and initializes a new RecoverySlicePacket struct from the
// given binary packet data.
func NewRecoverySlicePacket(data []byte) (RecoverySlicePacket, error) {
	typ := string(data[48:64])
	if typ != recoverySlicePacketType {
		return RecoverySlicePacket{}, fmt.Errorf("Recovery Slice packet type not as expected. Was %s", typ)
	}
	packetData := data[64:]

	packet := RecoverySlicePacket{}
	packet.Exponent = binary.LittleEndian.Uint32(packetData[0:4])
	copy(packet.Data[:], packetData[4:])

	return packet, nil
}

const creatorPacketType = "PAR 2.0\000Creator\000"

// CreatorPacket represents a Par 2.0 Creator packet.
type CreatorPacket struct {
	Creator string
}

// Type implements interface Packet to return the type of the Par 2.0 Recovery Slice Packet.
func (p CreatorPacket) Type() string {
	return creatorPacketType
}

// NewCreatorPacket creates and initializes a new CreatorPacket struct from the
// given binary packet data.
func NewCreatorPacket(data []byte) (CreatorPacket, error) {
	typ := string(data[48:64])
	if typ != creatorPacketType {
		return CreatorPacket{}, fmt.Errorf("Creator packet type not as expected. Was %s", typ)
	}
	packetData := data[64:]

	packet := CreatorPacket{}
	packet.Creator = strings.TrimRight(string(packetData[0:]), "\000")

	return packet, nil
}
