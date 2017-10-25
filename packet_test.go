package ais

import "testing"
import "github.com/stretchr/testify/assert"

var packet1 = "!AIVDM,1,1,,B,13@nocPP0427vl<`JO2``gwj08RD,0*11"
var packet1noChecksum = "!AIVDM,1,1,,B,13@nocPP0427vl<`JO2``gwj08RD,0"
var packet1wrongChecksum = "!AIVDM,1,1,,B,13@nocPP0427vl<`JO2``gwj08RD,0*12"

func TestParsePacket(t *testing.T) {
	correctPacket, err := ParsePacket(packet1)
	expected := &Packet{
		Talker:"AI",
		PacketType:"VDM",
		Channel:"B",
		FragCount:1,
		FragNo:1,
		SeqId:"",
		FillBits:0,
		Payload: "13@nocPP0427vl<`JO2``gwj08RD",
	}
	assert.Nil(t, err)
	assert.EqualValues(t, expected, correctPacket)
}

func TestWrongChecksum(t *testing.T) {
	_, err := ParsePacket(packet1wrongChecksum)
	assert.Equal(t, ErrIncorrectChecksum, err)
}

func TestMissingChecksum(t *testing.T) {
	_, err := ParsePacket(packet1noChecksum)
	assert.Equal(t, ErrMissingChecksum, err)
}