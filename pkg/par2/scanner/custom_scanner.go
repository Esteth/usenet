package scanner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type CustomScanner struct {
	reader io.Reader
	packet Packet
	err    error
}

type packetHeader struct {
	length        uint64
	md5Hash       [16]byte
	recoverySetId [16]byte
	packetType    [16]byte
}

func (c *CustomScanner) Next() {
	magicSequenceBytes := make([]byte, 8)
	if _, err := io.ReadFull(c.reader, magicSequenceBytes); err != nil {
		c.err = fmt.Errorf("could not peek for magic sequence: %w", err)
		return
	}
	if bytes.Compare(magicSequenceBytes, magicSequence) != 0 {
		c.err = fmt.Errorf("could not find magic packet header in data")
		return
	}

	header, err := readHeader(c.reader)
	if err != nil {
		c.err = fmt.Errorf("could not read packet header: %w", err)
	}

	packetTypeString := string(header.packetType[:])
	if packetTypeString == mainPacketType {
		c.packet, c.err = scanMainPacket(c.reader)
	} else if packetTypeString == fileDescriptionPacketType {
		c.packet, c.err = scanFileDescriptionPacket(c.reader)
	} else if packetTypeString == fileSliceChecksumPacketType {
		c.packet, c.err = scanFileSliceChecksumPacket(c.reader)
	} else if packetTypeString == recoverySlicePacketType {
		c.packet, c.err = scanRecoverySlicePacket(c.reader)
	} else if packetTypeString == creatorPacketType {
		c.packet, c.err = scanCreatorPacket(c.reader)
	} else {
		c.packet = &unknownPacket{
			typ: packetTypeString,
		}
	}
}

func scanMainPacket(reader io.Reader) (MainPacket, error) {
	var packet MainPacket

	sliceSize, err := readUint64(reader)
	if err != nil {
		err = fmt.Errorf("could not read slice size from main packet: %w", err)
		return packet, err
	}

	numRecoveryFiles, err := readUint32(reader)
	if err != nil {
		err = fmt.Errorf("could not read number of recovery files from main packet: %w", err)
		return packet, err
	}

	recoveryFileIDs := make([][16]byte, numRecoveryFiles)
	for i := uint32(0); i < numRecoveryFiles; i++ {
		if _, err = io.ReadFull(reader, recoveryFileIDs[i][:]); err != nil {
			return packet, fmt.Errorf("could not read md5 hash from packet header: %w", err)
		}
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

func scanFileDescriptionPacket(reader io.Reader) (packet FileDescriptionPacket, err error) {

}

func scanFileSliceChecksumPacket(reader io.Reader) (packet FileSliceChecksumPacket, err error) {

}

func scanRecoverySlicePacket(reader io.Reader) (packet RecoverySlicePacket, err error) {

}

func scanCreatorPacket(reader io.Reader) (packet CreatorPacket, err error) {

}

func readHeader(reader io.Reader) (header packetHeader, err error) {
	packetLengthBytes := make([]byte, 8)
	if _, err = io.ReadFull(reader, packetLengthBytes); err != nil {
		return
	}
	header.length = binary.LittleEndian.Uint64(packetLengthBytes)

	if _, err = io.ReadFull(reader, header.md5Hash[:]); err != nil {
		err = fmt.Errorf("could not read md5 hash from packet header: %w", err)
		return
	}

	if _, err = io.ReadFull(reader, header.recoverySetId[:]); err != nil {
		err = fmt.Errorf("could not read recover set ID from packet header: %w", err)
		return
	}

	if _, err = io.ReadFull(reader, header.packetType[:]); err != nil {
		err = fmt.Errorf("could not read packet type from packet header: %w", err)
		return
	}

	return
}

func readUint32(reader io.Reader) (result uint32, err error) {
	bytes := make([]byte, 4)
	if _, err = io.ReadFull(reader, bytes); err != nil {
		err = fmt.Errorf("could not little endian number: %w", err)
		return
	}
	result = binary.LittleEndian.Uint32(bytes)
	return
}

func readUint64(reader io.Reader) (result uint64, err error) {
	bytes := make([]byte, 8)
	if _, err = io.ReadFull(reader, bytes); err != nil {
		err = fmt.Errorf("could not little endian number: %w", err)
		return
	}
	result = binary.LittleEndian.Uint64(bytes)
	return
}
