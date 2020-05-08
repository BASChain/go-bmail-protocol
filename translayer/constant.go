package translayer

const (
	MIN_TYP uint16 = iota
	//bhello
	HELLO
	HELLO_ACK

	//bmtp
	SEND_ENVELOPE
	RESP_ENVELOPE

	SEND_CRYPT_ENVELOPE
	RESP_CRYPT_ENVELOPE

	//bpop
	STAT
	STAT_RESP
	LIST
	LIST_RESP
	RETR
	RETR_RESP
	DELETE
	DELETE_RESP

	CONTACT_HELLO
	CONTACT_HELLO_RESP
	CONTACT_ADD
	CONTACT_DEL
	CONTACT_PULL

	MAX_TYP
)

const (
	Uin8Size   int = 1
	Uint16Size int = 2
	Uint32Size int = 4
	Uint64Size int = 8
)
const BMTP_PORT = 1025
const BPOP3 = 1110

type EnveUniqID [16]byte
