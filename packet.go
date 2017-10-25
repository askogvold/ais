package ais

import (
	"errors"
	"strconv"
	"strings"
)

// Packet represents one line of AIS data
type Packet struct {
	Talker     string
	PacketType string
	FragCount  int
	FragNo     int
	SeqId      string
	Channel    string
	Payload    string
	FillBits   int
}

var (
	ErrEmptyPacket         = errors.New("empty packet")
	ErrInvalidPacketPrefix = errors.New("invalid prefix")
	ErrMissingChecksum     = errors.New("missing checksum")
	ErrIncorrectChecksum   = errors.New("incorrect checksum")
)

// ParsePacket parses one line of AIS data
func ParsePacket(rawPacket string) (*Packet, error) {
	l := len(rawPacket)
	if l == 0 {
		return nil, ErrEmptyPacket
	}

	if rawPacket[0] != '!' {
		return nil, ErrInvalidPacketPrefix
	}

	checksum, err := readChecksum(rawPacket)
	if err != nil {
		return nil, err
	}

	innerMessage := rawPacket[1 : l-3]
	calculatedChecksum := calculateChecksum(innerMessage)
	if checksum != calculatedChecksum {
		return nil, ErrIncorrectChecksum
	}

	parts := strings.Split(innerMessage, ",")

	return &Packet{
		Talker:     parts[0][0:2],
		PacketType: parts[0][2:],
		FragCount:  toInt(parts[1]),
		FragNo:     toInt(parts[2]),
		SeqId:      parts[3],
		Channel:    parts[4],
		Payload:    parts[5],
		FillBits:   toInt(parts[6]),
	}, nil
}

func readChecksum(rawPacket string) (byte, error) {
	l := len(rawPacket)
	if rawPacket[l-3] != '*' {
		return 0, ErrMissingChecksum
	}
	checksumN := rawPacket[l-2:]
	checksum, err := strconv.ParseUint(checksumN, 16, 8)
	if err != nil {
		panic(err)
	}
	return byte(checksum), nil
}
