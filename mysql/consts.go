package mysql

//ProtocolConst is const type for all mysql protocol state & variables
type ProtocolConst uint32

const (
	//prepare

	ComStmtExecute      ProtocolConst = 0
	ComStmtPrepare      ProtocolConst = 0
	ComStmtClose        ProtocolConst = 0
	ComStmtSendLongData ProtocolConst = 0
	ComStmtReset        ProtocolConst = 0

	//packet

	PacketHeaderOK   ProtocolConst = 0
	PacketHeaderRaw  ProtocolConst = 0xfa
	PacketHeaderNULL ProtocolConst = 0xfb
	PacketHeaderEOF  ProtocolConst = 0xfe
	PacketHeaderERR  ProtocolConst = 0xff

	//text protocol

	ComSleep           ProtocolConst = 0
	ComQuit            ProtocolConst = 0
	ComInitDB          ProtocolConst = 0
	ComQuery           ProtocolConst = 0
	ComFieldList       ProtocolConst = 0
	ComCreateDB        ProtocolConst = 0
	ComDropDB          ProtocolConst = 0
	ComRefresh         ProtocolConst = 0
	ComShutdown        ProtocolConst = 0
	ComStatistics      ProtocolConst = 0
	ComProcessInfo     ProtocolConst = 0
	ComConnect         ProtocolConst = 0
	ComProcessKill     ProtocolConst = 0
	ComDebug           ProtocolConst = 0
	ComPing            ProtocolConst = 0
	ComTime            ProtocolConst = 0
	ComDelayedInsert   ProtocolConst = 0
	ComChangeUser      ProtocolConst = 0
	ComResetConnection ProtocolConst = 0
	ComDaemon          ProtocolConst = 0

	//slave protocol

	ComBinLogDump     ProtocolConst = 0
	ComBinLogDumpGTID ProtocolConst = 0
	ComTableDump      ProtocolConst = 0
	ComConnectOut     ProtocolConst = 0
	ComRegisterSlave  ProtocolConst = 0

	//status flags

	ServerStatusInTrans            ProtocolConst = 0x0001
	ServerStatusAutoCommit         ProtocolConst = 0x0002
	ServerStatusRequestExits       ProtocolConst = 0x0008
	ServerStatusNoGoodIndexUsed    ProtocolConst = 0x0010
	ServerStatusNoIndexUsed        ProtocolConst = 0x0020
	ServerStatusCursorExists       ProtocolConst = 0x0040
	ServerStatusLastRowSent        ProtocolConst = 0x0080
	ServerStatusDBDroped           ProtocolConst = 0x0100
	ServerStatusNoBackslashEscapes ProtocolConst = 0x0200
	ServerStatusMetadataChanged    ProtocolConst = 0x0400
	ServerQueryWasSlow             ProtocolConst = 0x0800
	ServerPSOutParams              ProtocolConst = 0x1000
	ServerStatusInTransReadOnly    ProtocolConst = 0x2000
	ServerSessionStateChanged      ProtocolConst = 0x4000

	CapabilityFlagClientLongPassword               ProtocolConst = 0x00000001
	CapabilityFlagClientFoundRows                  ProtocolConst = 0x00000002
	CapabilityFlagClientLongFlag                   ProtocolConst = 0x00000004
	CapabilityFlagClientConnectWithDB              ProtocolConst = 0x00000008
	CapabilityFlagClientNoSchema                   ProtocolConst = 0x00000010
	CapabilityFlagClientCompress                   ProtocolConst = 0x00000020
	CapabilityFlagClientODBC                       ProtocolConst = 0x00000040
	CapabilityFlagClientLocalFiles                 ProtocolConst = 0x00000080
	CapabilityFlagClientIgnoreSpace                ProtocolConst = 0x00000100
	CapabilityFlagClientProtocol41                 ProtocolConst = 0x00000200
	CapabilityFlagClientInterActive                ProtocolConst = 0x00000400
	CapabilityFlagClientSSL                        ProtocolConst = 0x00000800
	CapabilityFlagClientIgnoreSigpipe              ProtocolConst = 0x00001000
	CapabilityFlagClientTransactions               ProtocolConst = 0x00002000
	CapabilityFlagClientReserved                   ProtocolConst = 0x00004000
	CapabilityFlagClientSecureConnection           ProtocolConst = 0x00008000
	CapabilityFlagClientMultiStatements            ProtocolConst = 0x00010000
	CapabilityFlagClientMultiResults               ProtocolConst = 0x00020000
	CapabilityFlagClientPSMultiResults             ProtocolConst = 0x00040000
	CapabilityFlagClientPluginAuth                 ProtocolConst = 0x00080000
	CapabilityFlagClientConnectAttrs               ProtocolConst = 0x00100000
	CapabilityFlagClientPluginAuthLenencClientData ProtocolConst = 0x00200000
	CapabilityFlagClientCanHandleExpiredPasswords  ProtocolConst = 0x00400000
	CapabilityFlagClientSessionTrack               ProtocolConst = 0x00800000
	CapabilityFlagClientDeprecateEOF               ProtocolConst = 0x01000000
)
