package mysql

import (
	"bufio"
	"fmt"
	"io"
)

const (
	MaxPayloadLen int = 1<<24 - 1
)

//Packet is a mapping to mysql packet
// len(Payload) = PayloadLength

type PacketIO struct {
	wb       io.Writer
	rb       io.Reader
	Sequence uint8
}

func NewPacketIO(conn Connection) *PacketIO {

	p := new(PacketIO)
	p.rb = bufio.NewReaderSize(conn, 1024)
	p.wb = conn
	p.Sequence = 0
	return p
}

func (p *PacketIO) Reset() {
	p.Sequence = 0
}

func (p *PacketIO) ReadPacket() (Packet, error) {
	header := []byte{0, 0, 0, 0}

	if _, err := io.ReadFull(p.rb, header); err != nil {
		return nil, err
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	if length < 1 {
		return nil, fmt.Errorf("invalid payload length %d", length)
	}

	sequence := uint8(header[3])

	if sequence != p.Sequence {
		return nil, fmt.Errorf("invalid sequence %d != %d", sequence, p.Sequence)
	}

	p.Sequence++

	data := make([]byte, length)
	if _, err := io.ReadFull(p.rb, data); err != nil {
		return nil, err
	} else {
		if length < MaxPayloadLen {
			return Packet(data), nil
		}

		var buf []byte
		buf, err = p.ReadPacket()
		if err != nil {
			return nil, err
		} else {
			return Packet(append(data, buf...)), nil
		}
	}
}

//data already have header
func (p *PacketIO) WritePacket(data Packet) error {
	length := len(data) - 4

	for length >= MaxPayloadLen {

		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff

		data[3] = p.Sequence

		if n, err := p.wb.Write(data[:4+MaxPayloadLen]); err != nil {
			return err
		} else if n != (4 + MaxPayloadLen) {
			return err
		} else {
			p.Sequence++
			length -= MaxPayloadLen
			data = data[MaxPayloadLen:]
		}
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)
	data[3] = p.Sequence

	if n, err := p.wb.Write(data); err != nil {
		return err
	} else if n != len(data) {
		return err
	} else {
		p.Sequence++
		return nil
	}
}
