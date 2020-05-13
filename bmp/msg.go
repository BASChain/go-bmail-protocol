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

type EnvelopeSyn struct {
	Mode uint16   `json:"mode"`
	SN   BMailSN  `json:"sn"`
	Sig  []byte   `json:"sig"`
	Hash []byte   `json:"hash"`
	Env  Envelope `json:"env"`
}

func (es *EnvelopeSyn) MsgType() uint16 {
	return translayer.SEND_CRYPT_ENVELOPE
}
func (es *EnvelopeSyn) VerifyHeader(header *Header) bool {
	return header.MsgTyp == translayer.SEND_CRYPT_ENVELOPE &&
		header.MsgLen != 0
}

type EnvelopeAck struct {
	NextSN BMailSN `json:"nextSN"`
	Hash   []byte  `json:"hash"`
	Sig    []byte  `json:"sig"`
}

func (ea *EnvelopeAck) MsgType() uint16 {
	return translayer.RESP_CRYPT_ENVELOPE
}
func (ea *EnvelopeAck) VerifyHeader(header *Header) bool {
	return header.MsgTyp == translayer.RESP_CRYPT_ENVELOPE &&
		header.MsgLen != 0
}
