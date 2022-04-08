package scanner

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const HEADER_LENGTH = 64

type CustomScanner struct {
	source ScannerSource
	packet Packet
	err    error
}

type ScannerSource interface {
	io.Reader
	io.Seeker
}

type packetHeader struct {
	packetLength  uint64
	md5Hash       [16]byte
	recoverySetId [16]byte
	packetType    [16]byte
}

func (c *CustomScanner) Next() {
	magicSequenceBytes := make([]byte, 8)
	if _, err := io.ReadFull(c.source, magicSequenceBytes); err != nil {
		c.err = fmt.Errorf("could not peek for magic sequence: %w", err)
		return
	}
	if bytes.Compare(magicSequenceBytes, magicSequence) != 0 {
		c.err = fmt.Errorf("could not find magic packet header in data")
		return
	}

	header, err := readHeader(c.source)
	if err != nil {
		c.err = fmt.Errorf("could not read packet header: %w", err)
	}

	packetTypeString := string(header.packetType[:])
	if packetTypeString == mainPacketType {
		c.packet, c.err = scanMainPacket(c.source, header)
	} else if packetTypeString == fileDescriptionPacketType {
		c.packet, c.err = scanFileDescriptionPacket(c.source)
	} else if packetTypeString == fileSliceChecksumPacketType {
		c.packet, c.err = scanFileSliceChecksumPacket(c.source, header)
	} else if packetTypeString == recoverySlicePacketType {
		c.packet, c.err = scanRecoverySlicePacket(c.source, filePath)
	} else if packetTypeString == creatorPacketType {
		c.packet, c.err = scanCreatorPacket(c.source, header)
	} else {
		c.packet = &unknownPacket{
			typ: packetTypeString,
		}
	}
}

func scanMainPacket(source ScannerSource, header packetHeader) (MainPacket, error) {
	var packet MainPacket

	sliceSize, err := readUint64(source)
	if err != nil {
		err = fmt.Errorf("could not read slice size from main packet: %w", err)
		return packet, err
	}

	numRecoveryFiles, err := readUint32(source)
	if err != nil {
		err = fmt.Errorf("could not read number of recovery files from main packet: %w", err)
		return packet, err
	}

	recoveryFileIDs := make([][16]byte, numRecoveryFiles)
	for i := uint32(0); i < numRecoveryFiles; i++ {
		if _, err = io.ReadFull(source, recoveryFileIDs[i][:]); err != nil {
			return packet, fmt.Errorf("could not read recovery file ID %i from main packet: %w", i, err)
		}
	}

	// The number of non-recovery file IDs is the remaining space in the packet
	numNonRecoveryFiles := (header.packetLength - HEADER_LENGTH - 12 - 16*uint64(numRecoveryFiles)) / 16
	nonRecoveryFileIDs := make([][16]byte, numNonRecoveryFiles)
	for i := uint64(0); i < numNonRecoveryFiles; i++ {
		if _, err = io.ReadFull(source, nonRecoveryFileIDs[i][:]); err != nil {
			return packet, fmt.Errorf("could not read non recovery file ID %i from main packet: %w", i, err)
		}
	}

	return MainPacket{
		SliceSize:          sliceSize,
		RecoveryFileIDs:    recoveryFileIDs,
		NonRecoveryFileIDs: nonRecoveryFileIDs,
	}, nil
}

func scanFileDescriptionPacket(reader io.Reader) (packet FileDescriptionPacket, err error) {
	if _, err = io.ReadFull(reader, packet.ID[:]); err != nil {
		return packet, fmt.Errorf("could not read file ID from file description packet: %w", err)
	}
	if _, err = io.ReadFull(reader, packet.MD5[:]); err != nil {
		return packet, fmt.Errorf("could not read MD5 hash from file description packet: %w", err)
	}
	if _, err = io.ReadFull(reader, packet.MD516[:]); err != nil {
		return packet, fmt.Errorf("could not read MD5-16 from file description packet: %w", err)
	}
	packet.Length, err = readUint64(reader)
	if err != nil {
		err = fmt.Errorf("could not read file length from file description packet: %w", err)
		return packet, err
	}

	return
}

func scanFileSliceChecksumPacket(source ScannerSource, header packetHeader) (packet FileSliceChecksumPacket, err error) {
	if _, err = io.ReadFull(source, packet.FileID[:]); err != nil {
		return packet, fmt.Errorf("could not read file ID from file slice checksum packet: %w", err)
	}
	numSlices := (header.packetLength - 16 - HEADER_LENGTH) / 20
	packet.SliceHashes = make([][16]byte, numSlices)
	packet.SliceCRC32s = make([][4]byte, numSlices)

	for i := uint64(0); i < numSlices; i++ {
		if _, err = io.ReadFull(source, packet.SliceHashes[i][:]); err != nil {
			return packet, fmt.Errorf("could not read hash %i from file slice checksum packet: %w", i, err)
		}
		if _, err = io.ReadFull(source, packet.SliceCRC32s[i][:]); err != nil {
			return packet, fmt.Errorf("could not read CRC32 %i from file slice checksum packet: %w", i, err)
		}
	}

	return
}

func scanRecoverySlicePacket(source ScannerSource, filePath string) (packet RecoverySlicePacket, err error) {
	if packet.Exponent, err = readUint32(source); err != nil {
		err = fmt.Errorf("could not read exponent from recovery slice packet", err)
		return
	}
	currentOffset, err := source.Seek(0, io.SeekCurrent)
	if err != nil {
		err = fmt.Errorf("could not read current file position while parsing recovery slice packet", err)
		return
	}
	packet.Data.FileOffset = uint32(currentOffset)
	packet.Data.FilePath = filePath

	return
}

func scanCreatorPacket(reader io.Reader, header packetHeader) (packet CreatorPacket, err error) {
	identifierBytes := make([]byte, header.packetLength-HEADER_LENGTH)
	if _, err = io.ReadFull(reader, identifierBytes); err != nil {
		err = fmt.Errorf("could not read creator packet identifier: %w", err)
		return
	}

	packet.Creator = strings.TrimRight(string(identifierBytes[:]), "\000")
	return
}

func readHeader(reader io.Reader) (header packetHeader, err error) {
	packetLengthBytes := make([]byte, 8)
	if _, err = io.ReadFull(reader, packetLengthBytes); err != nil {
		return
	}
	header.packetLength = binary.LittleEndian.Uint64(packetLengthBytes)

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
