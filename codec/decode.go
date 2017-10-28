package codec

import "errors"

var ErrInvalidCharacter = errors.New("invalid character")

const (
	b001111 byte = 0x0F
	b110000 byte = 0x30
	b111100 byte = 0x3C
	b000011 byte = 0x03
	errByte byte = 0xFF
)

// ConvertPayload Converts a payload string to a byte slice
func ConvertPayload(payload string, fillBits int) (*Payload, error) {
	_ = fillBits // may be used in the future
	builder := newPayloadBuilder()
	for i := range payload {
		b := convertRune(payload[i])
		if b == errByte {
			return nil, ErrInvalidCharacter
		}
		builder.insertSixBits(b)
	}
	return builder.payload(), nil
}

func convertRune(r byte) byte {
	switch {
	case r >= '0' && r <= 'W':
		return byte(r - '0')
	case r >= '`' && r <= 'w':
		return byte(r - '`' + 40)
	default:
		return errByte
	}
}

func byteToRune(r byte) byte {
	switch {
	case r < 32:
		return r + 64
	case r < 64:
		return r + 32
	default:
		return errByte
	}
}

// Payload represents a payload converted to bytes.
type Payload struct {
	Bytes []byte
}

func (p *Payload) GetBits(From, Until int) []byte {
	b0 := From / 8
	nBits := Until - From
	if nBits <= 0 {
		return []byte{}
	}
	size := 1 + (nBits-1)/8
	result := make([]byte, size)

	leftShift := byte(From % 8)
	rightShift := 8 - leftShift

	dropLastBits := byte(8*size - nBits)

	for i := 0; i < size-1; i++ {
		result[i] = p.Bytes[b0+i]<<leftShift | (p.Bytes[b0+i+1] >> rightShift)
	}

	result[size-1] = ((p.Bytes[b0+size-1] >> dropLastBits) << dropLastBits) << leftShift

	_ = leftShift

	return result
}

type payloadBuilder struct {
	sixByteIndex int
	bytes        []byte
}

func newPayloadBuilder() payloadBuilder {
	return payloadBuilder{
		sixByteIndex: 0,
		bytes:        make([]byte, 0, 0),
	}
}

// insert <bitCount> last bits from byte
// bitCount range [1-6]
func (p *payloadBuilder) insertSixBits(b byte) {
	mode := p.sixByteIndex % 4
	switch mode {
	case 0: // [0-6)
		p.appendByte(b << 2)
	case 1: // [6-8) + [0-4)
		p.lastOr((b & b110000) >> 4)
		p.appendByte((b & b001111) << 4)
	case 2: // [4-8) + [0-2)
		p.lastOr((b & b111100) >> 2)
		p.appendByte((b & b000011) << 6)
	case 3: // [2-8)
		p.lastOr(b)
	}
	p.sixByteIndex += 1
}

// last byte in array will be or'ed with input byte
func (p *payloadBuilder) lastOr(b byte) {
	lastIndex := len(p.bytes) - 1
	p.bytes[lastIndex] |= b
}

func (p *payloadBuilder) appendByte(b byte) {
	p.bytes = append(p.bytes, b)
}

// convert to payload
func (p *payloadBuilder) payload() *Payload {
	return &Payload{
		Bytes: p.bytes,
	}
}
