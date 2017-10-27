package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	b101010   byte = 0x2A
	b10101010 byte = 0xAA
	b10000000 byte = 0x80
	b10100000 byte = 0xA0
)

func TestPayloadBuilder(t *testing.T) {
	p := newPayloadBuilder()
	p.insertSixBits(b101010)
	assert.Equal(t, p.bytes[0], b101010<<2)

	p.insertSixBits(b101010)
	assert.Equal(t, p.bytes[0], b10101010)

	assert.Equal(t, p.bytes[1], b10100000)

	p.insertSixBits(b101010)
	assert.Equal(t, p.bytes[1], b10101010)
	assert.Equal(t, p.bytes[2], b10000000)

	p.insertSixBits(b101010)
	assert.Equal(t, p.bytes[2], b10101010)
}

func TestDecodePayload(t *testing.T) {
	p, err := ConvertPayload("13@nocPP0427vl<`JO2``gwj08RDr", 0)
	assert.Nil(t, err)

	expected := []byte{
		0x04, 0x34, 0x36, 0xde,
		0xb8, 0x20, 0x00, 0x40,
		0x87, 0xfb, 0x43, 0x28,
		0x69, 0xf0, 0xa8, 0xa2,
		0xff, 0xf2, 0x00, 0x88,
		0x94, 0xe8,
	}
	assert.EqualValues(t, expected, p.Bytes)

	q, err := ConvertPayload("88888888880", 2)
	assert.Nil(t, err)
	expected2 := []byte{
		0x20, 0x82, 0x8,
		0x20, 0x82, 0x8,
		0x20, 0x80, 0x0,
	}
	assert.EqualValues(t, expected2, q.Bytes)
}

func TestFailOnBogusCharacter(t *testing.T) {
	_, err := ConvertPayload("13@nocPP0427vl<`JO2``gwj08RDæ", 0)
	assert.EqualValues(t, ErrInvalidCharacter, err)
}
