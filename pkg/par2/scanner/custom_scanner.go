package scanner

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

type CustomScanner struct {
	reader  bufio.Reader
	packets []Packet
}

func (c *CustomScanner) Scan() error {
	expectedMagicSequence, err := c.reader.Peek(8)
	if err != nil {
		return fmt.Errorf("could not peek for magic sequence: %w", err)
	}

	if bytes.Compare(expectedMagicSequence, magicSequence) != 0 {
		return fmt.Errorf("Could not find magic packet header in data")
	}
	// We've found the magic header. The next 8 bytes define the length of the packet.
	pktLength := binary.LittleEndian.Uint64(data[8:])
	if uint64(len(data)) < pktLength {
		return 0, nil, nil
	}
	return int(pktLength), data[:pktLength], nil
}
