package packet

import (
	"fmt"

	"github.com/piot/hasty-protocol/deserialize"
)

// Type : The packet type (Publish, Subscribe, Unsubscribe, etc)
type Type byte

// Nop : The different types
const (
	Nop Type = 128 + iota
	PublishStream
	SubscribeStream
	UnsubscribeStream
	CreateStream
	CreateStreamResult
	StreamData
	Connect
	ConnectResult
	Ping
	Pong
)

// Header : A Packet header
type Header struct {
	packetType      Type
	payloadSize     int
	headerOctetSize int
}

func (in Header) String() string {
	return fmt.Sprintf("[Header type:%02X payloadSize:%d]", in.packetType, in.payloadSize)
}

func convertFromOctetToPacketConst(t byte) (Type, error) {
	packetType := Type(t)
	if packetType < Nop || packetType > Pong {
		return Nop, fmt.Errorf("Illegal packet type:%d", t)
	}

	return packetType, nil
}

func tryReadLength(buf []byte) (int, int, error) {
	length, octetsUsed, err := deserialize.ToSmallLength(buf)
	return int(length), octetsUsed, err
}

// CheckIfWeHavePacketHeader : Returns a packet if it is ready
func CheckIfWeHavePacketHeader(buf []byte) (packet Header, packetWasReady bool, err error) {
	bufSize := len(buf)
	if bufSize < 1 {
		return Header{}, false, nil
	}
	packetLength, octetsUsedForLengthEncoding, _ := tryReadLength(buf)
	octetSize := len(buf)
	if octetSize < packetLength+octetsUsedForLengthEncoding {
		return Header{}, false, nil
	}
	t, convertErr := convertFromOctetToPacketConst(buf[octetsUsedForLengthEncoding])
	if convertErr != nil {
		return Header{}, false, convertErr
	}
	headerOctetSize := octetsUsedForLengthEncoding + 1

	return Header{packetType: t, payloadSize: packetLength - 1, headerOctetSize: headerOctetSize}, true, nil
}
