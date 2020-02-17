package par2

import (
	"encoding/binary"
	"fmt"
)

const mainPacketType = "PAR 2.0\000Main\000\000\000\000"

// MainPacket represents a Par 2.0 Main Packet.
type MainPacket struct {
	typ                []byte
	RecoveryFileIDs    [][16]byte
	NonRecoveryFileIDs [][16]byte
}

// Type implements interface Packet to return the type of the Par 2.0 Main Packet.
func (p MainPacket) Type() []byte {
	return p.typ
}

// NewMainPacket creates and initializes a new MainPacket struct from the given binary packet data.
func NewMainPacket(data []byte) (MainPacket, error) {
	typ := string(data[48:64])
	if typ != mainPacketType {
		return MainPacket{}, fmt.Errorf("Main packet type not as expected. Was %s", typ)
	}
	mainPacketData := data[64:]

	numRecoveryFiles := binary.LittleEndian.Uint32(mainPacketData[8:12])
	recoveryFileIDs := make([][16]byte, numRecoveryFiles)
	for i := uint32(0); i < numRecoveryFiles; i++ {
		copy(recoveryFileIDs[i][:], mainPacketData[12+16*i:28+16*i])
	}

	numNonRecoveryFiles := (binary.LittleEndian.Uint64(data[8:16]) - 28 + 16*uint64(numRecoveryFiles)) / 16
	nonRecoveryFilesBaseIndex := uint64(12+16*numRecoveryFiles)
	nonRecoveryFileIDs := make([][16]byte, numNonRecoveryFiles)
	for i := uint64(0); i < numNonRecoveryFiles; i++ {
		copy(nonRecoveryFileIDs[i][:], mainPacketData[nonRecoveryFilesBaseIndex + 16 * i: nonRecoveryFilesBaseIndex + 16 + 16 * i])
	}

	return MainPacket{
		typ:                []byte(mainPacketType),
		RecoveryFileIDs:    recoveryFileIDs,
		NonRecoveryFileIDs: nonRecoveryFileIDs,
	}, nil
}
