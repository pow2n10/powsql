package mysql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Packet []byte

func (packet Packet) Header() ProtocolConst {
	return ProtocolConst(packet[0])
}

type GenericResponsePacket interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte, capability ProtocolConst) error
	Type() ProtocolConst
}

type PacketOK struct {
	AffectedRows        uint64
	LastInsertID        uint64
	StatusFlags         uint16
	Warnings            uint16
	Info                string
	SessionStateChanges string
}

func (packet *PacketOK) Marshal() ([]byte, error) {
	return nil, nil
}

func (packet *PacketOK) Unmarshal(data []byte, capability ProtocolConst) error {

	pos := 1

	n, m, e := ReadUint64(data[pos:])

	if e != nil {
		return e
	}

	packet.AffectedRows = n
	pos += m

	n, m, e = ReadUint64(data[pos:])

	if e != nil {
		return e
	}

	packet.LastInsertID = n
	pos += m

	if capability&CapabilityFlagClientProtocol41 > 0 {
		packet.StatusFlags = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
		packet.Warnings = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
	} else if capability&CapabilityFlagClientTransactions > 0 {
		packet.StatusFlags = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
	}

	if capability&CapabilityFlagClientSessionTrack > 0 {
		index := bytes.IndexByte(data[pos:], 0x00)
		packet.Info = string(data[pos : pos+index])
		pos += index
	}

	if ProtocolConst(packet.StatusFlags)&ServerSessionStateChanged > 0 {
		index := bytes.IndexByte(data[pos:], 0x00)
		packet.SessionStateChanges = string(data[pos : pos+index])
		pos += index
	}

	return nil
}

func (packet *PacketOK) Type() ProtocolConst {
	return PacketHeaderOK
}

type PacketERR struct {
	ErrCode        uint16
	SqlStateMarker string
	SqlState       string
	ErrorMessage   string
}

func (packet *PacketERR) Marshal() ([]byte, error) {
	return nil, nil
}

func (packet *PacketERR) Unmarshal(data []byte, capability ProtocolConst) error {

	pos := 1

	packet.ErrCode = binary.LittleEndian.Uint16(data[pos : pos+2])
	pos += 2

	if capability&CapabilityFlagClientProtocol41 > 0 {
		packet.SqlStateMarker = string(data[pos : pos+1])
		pos++

		packet.SqlState = string(data[pos : pos+5])
		pos += 5
	}

	packet.ErrorMessage = string(data[pos:])
	return nil
}

func (packet *PacketERR) Type() ProtocolConst {
	return PacketHeaderERR
}

type PacketEOF struct {
	Warnings    uint16
	StatusFlags uint16
}

func (packet *PacketEOF) Marshal() ([]byte, error) {
	return nil, nil
}

func (packet *PacketEOF) Unmarshal(data []byte, capability ProtocolConst) error {
	pos := 1
	if capability&CapabilityFlagClientProtocol41 > 0 {

		packet.Warnings = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		packet.StatusFlags = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
	}
	return nil
}

func (packet *PacketEOF) Type() ProtocolConst {
	return PacketHeaderEOF
}

func NewGenericResponsePacket(packet Packet, capability ProtocolConst) (GenericResponsePacket, error) {

	var p GenericResponsePacket

	switch packet[0] {
	case 0x00:
		p = new(PacketOK)
	case 0xff:
		p = new(PacketERR)
	case 0xfe:
		p = new(PacketEOF)
	default:
		return nil, errors.New("unknow packet type:" + fmt.Sprintf("%#v", packet[0]))
	}
	err := p.Unmarshal([]byte(packet), capability)
	return p, err
}
