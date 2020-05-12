package bmp

const (
	BMailModeP2P = iota
	BMailModeP2S

	RcpTypeTo = iota
	RcpTypeCC
	RcpTypeBcc
)

type Packable interface {
	Pack() ([]byte, error)
	UnPack([]byte) error
}

type Envelope interface {
	Packable
	MsgType() uint16
	VerifyHeader(header *Header) bool
}
