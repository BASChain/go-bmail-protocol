package translayer

const (
	MIN_TYP uint16 = iota
	//bhello
	HELLO
	HELLO_ACK

	SEND_SIGNATURE
	VALIDATE_SIGNATURE

	//bmtp
	SEND_ENVELOPE
	RESP_ENVELOPE

	BMAIL_UNREACHABLE
	RESP_BMAIL_UNRECHABLE

	//bpop
	STAT
	STAT_RESP
	LIST
	LIST_RESP
	LATEST
	LATEST_RESP

	MAX_TYP
)

const (
	Uin8Size   int = 1
	Uint16Size int = 2
	UintSize   int = 4
	Uint32Size int = 4
	Uint64Size int = 8
)
