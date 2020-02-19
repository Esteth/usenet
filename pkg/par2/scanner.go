package par2

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var magicSequence = []byte{'P', 'A', 'R', '2', '\000', 'P', 'K', 'T'}

type Packet interface {
	Type() []byte
}

type unknownPacket struct {
	typ []byte
}

func (p *unknownPacket) Type() []byte {
	return p.typ
}

// A Scanner allows a PAR2 stream to be parsed packet-by-packet in the style
// of bufio.Scanner
type Scanner struct {
	scanner bufio.Scanner
	packet  Packet
	err     error
}

// NewScanner creates a new Scanner reading the given reader.
func NewScanner(r io.Reader) *Scanner {
	z := new(Scanner)
	z.Reset(r)
	return z
}

// Reset discards the Scanner z's state and makes it equivalent to the
// result of it's original state from NewScanner, but reading from r instead.
// This permits reusing a reader rather than allocating a new one.
func (z *Scanner) Reset(r io.Reader) {
	*z = Scanner{
		scanner: *bufio.NewScanner(r),
		packet:  nil,
		err:     nil,
	}
	z.scanner.Split(scanPackets)
}

// Scan advances the scanner to the next packet, which will then be available through the Packet
// method.
// It returns false when the scan stops, either by reaching the end of the input or an error.
// After Scan returns false, the Err method will return any error that occurred during scanning,
// except that if it was io.EOF, Err will return nil.
func (z *Scanner) Scan() bool {
	if z.scanner.Scan() {
		packetType := z.scanner.Bytes()[48:64]
		if string(packetType) == mainPacketType {
			z.packet, z.err = NewMainPacket(z.scanner.Bytes())
		} else if string(packetType) == fileDescriptionPacketType {
			z.packet, z.err = NewFileDescriptionPacket(z.scanner.Bytes())
		} else if string(packetType) == fileSliceChecksumPacketType {
			z.packet, z.err = NewFileSliceChecksumPacket(z.scanner.Bytes())
		} else if string(packetType) == recoverySlicePacketType {
			z.packet, z.err = NewRecoverySlicePacket(z.scanner.Bytes())
		} else if string(packetType) == creatorPacketType {
			z.packet, z.err = NewCreatorPacket(z.scanner.Bytes())
		} else {
			z.packet = &unknownPacket{
				typ: packetType,
			}
		}
		if z.scanner.Err() != nil && z.scanner.Err() != io.EOF {
			z.err = z.scanner.Err()
		}
		return true
	}
	return false
}

// Err returns the first non-EOF error that was encountered by the scanner.
func (z *Scanner) Err() error {
	return z.err
}

// Packet returns the most recent packet generated by a call to Scan.
func (z *Scanner) Packet() Packet {
	return z.packet
}

func scanPackets(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// We can only identify a packet if we can read the full magic sequence and length
	if len(data) < 16 {
		return 0, nil, nil
	}
	if bytes.Compare(data[:8], magicSequence) != 0 {
		return 0, nil, fmt.Errorf("Could not find magic packet header in data")
	}
	// We've found the magic header. The next 8 bytes define the length of the packet.
	pktLength := binary.LittleEndian.Uint64(data[8:])
	if uint64(len(data)) < pktLength {
		return 0, nil, nil
	}
	return int(pktLength), data[:pktLength], nil
}
