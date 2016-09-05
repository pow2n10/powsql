package mysql

type ComHandler func(conn *BackendConnection, packet Packet)

var comHandler map[ProtocolConst]ComHandler
