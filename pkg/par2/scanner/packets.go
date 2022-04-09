package scanner

type Packet interface {
	Type() string
}

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

const fileDescriptionPacketType = "PAR 2.0\000FileDesc"

// FileDescriptionPacket represents a Par 2.0 File Description Packet.
type FileDescriptionPacket struct {
	ID         [16]byte
	MD5        [16]byte
	MD516      [16]byte
	FileLength uint64
	FileName   string
}

// Type implements interface Packet to return the type of the Par 2.0 File Description Packet.
func (p FileDescriptionPacket) Type() string {
	return fileDescriptionPacketType
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

const recoverySlicePacketType = "PAR 2.0\000RecvSlic"

// RecoverySlicePacket represents a Par 2.0 Recovery Slice packet.
type RecoverySlicePacket struct {
	FileID                 [16]byte
	Exponent               uint32
	RecoveryDataFilePath   string
	RecoveryDataFileOffset uint32
}

// Type implements interface Packet to return the type of the Par 2.0 Recovery Slice Packet.
func (p RecoverySlicePacket) Type() string {
	return recoverySlicePacketType
}

const creatorPacketType = "PAR 2.0\000Creator\000"

// CreatorPacket represents a Par 2.0 Creator packet.
type CreatorPacket struct {
	Creator string
}

// Type implements interface Packet to return the type of the Par 2.0 Creator Packet.
func (p CreatorPacket) Type() string {
	return creatorPacketType
}

// unknownPacket represents a packet of a type we cannot parse.
type unknownPacket struct {
	typ string
}

// Type implements interface Packet to return the type of the Par 2.0 Recovery Slice Packet.
func (p *unknownPacket) Type() string {
	return p.typ
}
