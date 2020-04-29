package bmprotocol

type BMailHash [32]byte

type BMailAddrss struct {
	MailAddress string
	Alias       string
	Groupid     []byte
}
