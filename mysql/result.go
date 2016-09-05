package mysql

import "encoding/binary"

//https://dev.mysql.com/doc/internals/en/com-query-response.html#packet-COM_QUERY_Response
type Field struct {
	Catalog      string
	Schema       string
	Table        string
	OrgTable     string
	Name         string
	OrgName      string
	Length       uint32
	Charset      uint16
	ColumnLength uint32
	Type         uint8
	Flags        uint16
	Decimals     uint8
	Filler       uint16
	DefaultValue []byte
}

func (field *Field) Unmarshal(data []byte, capability ProtocolConst) error {

	offset := 0

	if s, n, err := ReadLenencStr(data[offset:]); err == nil {
		field.Catalog = string(s)
		offset += n
	} else {
		return err
	}

	if s, n, err := ReadLenencStr(data[offset:]); err == nil {
		field.Schema = string(s)
		offset += n
	} else {
		return err
	}

	if s, n, err := ReadLenencStr(data[offset:]); err == nil {
		field.Table = string(s)
		offset += n
	} else {
		return err
	}

	if s, n, err := ReadLenencStr(data[offset:]); err == nil {
		field.OrgTable = string(s)
		offset += n
	} else {
		return err
	}

	if s, n, err := ReadLenencStr(data[offset:]); err == nil {
		field.Name = string(s)
		offset += n
	} else {
		return err
	}

	if s, n, err := ReadLenencStr(data[offset:]); err == nil {
		field.OrgName = string(s)
		offset += n
	} else {
		return err
	}

	if s, n, err := ReadUint64(data[offset:]); err == nil {
		field.Length = uint32(s)
		offset += n
	} else {
		return err
	}

	field.Charset = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	field.Length = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	field.Type = byte(data[offset])
	offset++

	field.Flags = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	field.Decimals = byte(data[offset])
	offset++

	field.Filler = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2
	if capability&ComFieldList > 0 {
		if s, n, err := ReadUint64(data[offset:]); err == nil {
			field.Length = uint32(s)
			offset += n
		} else {
			return err
		}
		field.DefaultValue = data[offset:field.Length]
	}
	return nil
}

type Row []byte

type Result struct {
	Fields map[string]*Field
	Rows   []Row

	Errcode uint16
	Errmsg  string

	AffectedRows uint64
	LastInsertID uint64

	Status   uint16
	Warnings uint16
}
