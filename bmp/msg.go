package bmp

import (
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/translayer"
)

type EnvelopeMsg interface {
	MsgType() uint16
	VerifyHeader(header *Header) bool
}

type Header struct {
	Ver    uint16 `json:"ver"`
	MsgTyp uint16 `json:"t"`
	MsgLen int    `json:"l"`
}

func (h *Header) GetLen() int8 {
	if h.Ver == translayer.BMAILVER1 {
		return 8
	}
	return -1
}

type HELO struct {
}

type HELOACK struct {
	SN     BMailSN       `json:"sn"`
	SrvBca bmail.Address `json:"srv"`
}

func (ha *HELOACK) MsgType() uint16 {
	return translayer.HELLO_ACK
}

func (ha *HELOACK) VerifyHeader(header *Header) bool {
	return header.MsgLen != 0 && header.MsgTyp == translayer.HELLO_ACK
}
