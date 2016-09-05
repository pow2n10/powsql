package mysql

import (
	"bytes"
	"encoding/binary"
)

func (conn *BackendConnection) parseInitialHanshake(data []byte) error {

	offset := 0
	nextOffset := 0

	protocolVersion := uint8(data[0])
	conn.protocolVersion = protocolVersion
	offset++

	has := bytes.IndexByte(data[offset:], 0x00)
	nextOffset = offset + has
	backendVersion := string(data[offset:nextOffset])
	nextOffset++
	conn.serverVersion = backendVersion
	offset = nextOffset
	nextOffset = offset + 4

	connectionID := binary.LittleEndian.Uint32(data[offset:nextOffset])
	conn.connectionID = connectionID
	offset = nextOffset

	//auth-plugin-data-part-1
	nextOffset = offset + 8
	authPluginDataPart1 := string(data[offset:nextOffset])
	conn.authPluginDataPart1 = authPluginDataPart1
	offset = nextOffset

	//skip filler
	nextOffset++
	//	filler := uint8(data[nextOffset])
	offset = nextOffset
	nextOffset = offset + 2

	capability := uint32(binary.LittleEndian.Uint16(data[offset:nextOffset]))
	offset = nextOffset

	conn.capability = ProtocolConst(capability)

	nextOffset++
	//charset := uint8(data[offset])
	offset = nextOffset
	nextOffset = offset + 2
	status := binary.LittleEndian.Uint16(data[offset:nextOffset])
	conn.status = status
	offset = nextOffset

	nextOffset = offset + 2
	capability = uint32(binary.LittleEndian.Uint16(data[offset:nextOffset]))<<16 | capability
	conn.capability = ProtocolConst(capability)
	offset = nextOffset

	nextOffset++
	authPluginDataLength := uint8(data[offset])

	offset = nextOffset

	nextOffset = offset + 10

	//skip reserved

	offset = nextOffset

	if capability&uint32(CapabilityFlagClientSecureConnection) > 0 {
		partLen := authPluginDataLength - 8
		if partLen > 12 {
			partLen = 12
		}
		nextOffset = offset + int(partLen)
		authPluginDataPart2 := string(data[offset:nextOffset])
		offset = nextOffset
		conn.authPluginDataPart1 += authPluginDataPart2
	}

	authPluginName := ""

	if capability&uint32(CapabilityFlagClientPluginAuth) > 0 {
		has := bytes.IndexByte(data[offset:], 0x00)
		nextOffset = offset + has
		authPluginName = string(data[offset:nextOffset])
		offset = nextOffset
	}
	conn.authPluginName = authPluginName

	return nil
}

func (conn *BackendConnection) writeInitialHanshake() error {

	capability := CapabilityFlagClientProtocol41 | CapabilityFlagClientSecureConnection | CapabilityFlagClientLongPassword | CapabilityFlagClientTransactions | CapabilityFlagClientLongFlag

	capability = capability & conn.capability

	auth := CalcPassword([]byte(conn.authPluginDataPart1), []byte(conn.auth.Password))

	packetLength := 4 + 4 + 1 + 23 + len(conn.auth.User) + 1 + len(auth) + 1

	if len(conn.dbname) > 0 {
		capability |= CapabilityFlagClientConnectWithDB
		packetLength += len(conn.dbname) + 1
	}

	conn.capability = capability
	packetData := make([]byte, packetLength+4)

	packetData[4] = byte(capability)
	packetData[5] = byte(capability >> 8)
	packetData[6] = byte(capability >> 16)
	packetData[7] = byte(capability >> 24)

	packetData[12] = byte(0x00)

	pos := 13 + 23

	if len(conn.auth.User) > 0 {
		pos += copy(packetData[pos:], conn.auth.User)
	}
	//	packetData[pos] = 0
	pos++

	packetData[pos] = byte(len(auth))
	pos += 1 + copy(packetData[pos+1:], auth)

	if capability&CapabilityFlagClientConnectWithDB > 0 {
		pos += copy(packetData[pos:], conn.dbname)
	}

	return conn.packetIO.WritePacket(packetData)
}

func (conn *BackendConnection) readInitialHanshake() error {

	data, err := conn.packetIO.ReadPacket()
	if err != nil {
		return err
	}
	return conn.parseInitialHanshake([]byte(data))
}

func (conn *BackendConnection) UesDB(dbname string) error {

	if dbname == conn.dbname {
		return nil
	}

	return nil
}
