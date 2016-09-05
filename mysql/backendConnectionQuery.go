package mysql

import "errors"

func (conn *BackendConnection) Command(cmd ProtocolConst, data []byte) (*Result, error) {
	conn.packetIO.Reset()
	packet := make([]byte, len(data)+1+4)
	packet[4] = byte(cmd)
	copy(packet[5:], data)
	err := conn.packetIO.WritePacket(packet)

	if err != nil {
		return nil, err
	}

	readPacket, err := conn.packetIO.ReadPacket()

	if err != nil {
		return nil, err
	}

	result := &Result{}

	if readPacket[0] == byte(PacketHeaderERR) {
		packet, err := NewGenericResponsePacket(readPacket, conn.capability)
		if err != nil {
			return nil, err
		}
		result.Errcode = packet.(*PacketERR).ErrCode
		result.Errmsg = packet.(*PacketERR).ErrorMessage
		return result, nil
	}

	if readPacket[0] == byte(PacketHeaderOK) {

		packet, err := NewGenericResponsePacket(readPacket, conn.capability)
		if err != nil {
			return nil, err
		}
		result.AffectedRows = packet.(*PacketOK).AffectedRows
		result.LastInsertID = packet.(*PacketOK).LastInsertID
		return result, nil
	}

	if readPacket[0] == byte(PacketHeaderNULL) {
		return nil, errors.New("unsupport loca infile")
	}

	cols, _, err := ReadUint64(readPacket)

	if err != nil {
		return nil, err
	}

	if err = conn.readColumns(cmd, result, uint32(cols)); err != nil {
		return nil, err
	}

	if err = conn.readRows(cmd, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (conn *BackendConnection) readColumns(cmd ProtocolConst, result *Result, cols uint32) error {

	result.Fields = make(map[string]*Field, 0)

	for {
		readPacket, err := conn.packetIO.ReadPacket()

		if err != nil {
			return err
		}

		if readPacket[0] == byte(PacketHeaderEOF) {
			packet, err := NewGenericResponsePacket(readPacket, conn.capability)
			if err != nil {
				return err
			}
			result.Status = packet.(*PacketEOF).StatusFlags
			result.Warnings = packet.(*PacketEOF).Warnings
			break
		}

		field := new(Field)

		err = field.Unmarshal(readPacket, cmd)
		if err != nil {
			return err
		}
		result.Fields[field.Name] = field
	}
	return nil
}

func (conn *BackendConnection) readRows(cmd ProtocolConst, result *Result) error {

	result.Rows = make([]Row, 0)

	for {

		readPacket, err := conn.packetIO.ReadPacket()

		if err != nil {
			return err
		}

		if readPacket[0] == byte(PacketHeaderEOF) {

			packet, err := NewGenericResponsePacket(readPacket, conn.capability)
			if err != nil {
				return err
			}
			result.Status = packet.(*PacketEOF).StatusFlags
			result.Warnings = packet.(*PacketEOF).Warnings
			break
		}

		result.Rows = append(result.Rows, Row(readPacket))

	}
	return nil

}

func (conn *BackendConnection) ComQuery(query string) (*Result, error) {
	return conn.Command(ComQuery, []byte(query))
}
