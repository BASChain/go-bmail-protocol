package bmp

import (
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"github.com/BASChain/go-bmail-protocol/translayer"
)

const (
	BMailModeP2P = iota
	BMailModeP2S

	RcpTypeTo  = iota
	RcpTypeCC  = iota
	RcpTypeBcc = iota
)

type Envelope struct {
	bmprotocol.EnvelopeSig `json:"sigData"`
	EId                    translayer.EnveUniqID `json:"eid"`
	From                   bmail.Address         `json:"from"`
	Mode                   int                   `json:"mode"`
}

type Receipt struct {
	RcpAddr bmail.Address `json:"rcpAddr"`
	RcpTyp  int           `json:"rcpType"`
}

type PlainBody struct {
	Receipts []Receipt `json:"rcps"`
	Subject  []byte    `json:"subject"`
	MsgBody  []byte    `json:"msgBody"`
}

func (pb *PlainBody) CryptBy(peer bmail.Address) []byte {
	return nil
}

type CryptBody struct {
	PeerAddr  bmail.Address `json:"receiver"`
	CryptData []byte        `json:"cryptBody"`
}
