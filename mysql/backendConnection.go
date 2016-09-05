package mysql

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type BackendConnection struct {
	addr                string
	conn                net.Conn
	protocolVersion     uint8
	serverVersion       string
	connectionID        uint32
	capability          ProtocolConst
	status              uint16
	authPluginDataPart1 string
	authPluginDataPart2 string
	authPluginName      string
	packetIO            *PacketIO
	auth                AuthInfo
	dbname              string
	lastCheck           time.Time
}

func NewBackendConnection(addr string, auth *AuthInfo, dbname string) (*BackendConnection, error) {
	conn := new(BackendConnection)
	conn.auth = *auth
	conn.addr = addr
	conn.dbname = dbname
	err := conn.realConnect()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (conn *BackendConnection) realConnect() error {

	if conn.conn != nil {
		conn.conn.Close()
	}

	n := "tcp"
	if strings.Contains(conn.addr, "/") {
		n = "unix"
	}

	netConn, err := net.Dial(n, conn.addr)
	if err != nil {
		return err
	}

	tcpConn := netConn.(*net.TCPConn)

	tcpConn.SetNoDelay(false)
	tcpConn.SetKeepAlive(true)
	conn.conn = tcpConn
	conn.packetIO = NewPacketIO(conn.conn)

	if err := conn.readInitialHanshake(); err != nil {
		conn.conn.Close()
		return err
	}

	if err := conn.writeInitialHanshake(); err != nil {
		conn.conn.Close()
		return err
	}

	responsePacket, err := conn.packetIO.ReadPacket()
	if err != nil {
		conn.conn.Close()
		return err
	}

	packet, err := NewGenericResponsePacket(responsePacket, conn.capability)

	if err != nil {
		conn.conn.Close()
		return err
	}
	if packet.Type() != PacketHeaderOK {
		return errors.New(fmt.Sprintf("%#v", packet))
	}

	return nil

}

func (conn *BackendConnection) ChangeUser(auth *AuthInfo) error {
	if auth.User == conn.auth.User && auth.Password == conn.auth.Password {
		return nil
	}
	return nil
}

func (conn *BackendConnection) Ping() error {
	return nil
}

func (conn *BackendConnection) Reset() error {
	return nil
}
