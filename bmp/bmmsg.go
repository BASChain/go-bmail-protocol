package bmp

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/translayer"
)

type EnvelopeMsg interface {
	MsgType() uint16
	VerifyHeader(header *Header) bool
	GetBytes() ([]byte, error)
}

type Header struct {
	Ver    uint16 `json:"ver"`
	MsgTyp uint16 `json:"t"`
	MsgLen int    `json:"l"`
}

func (h *Header) GetLen() int8 {
	switch h.Ver {
	case translayer.BMAILVER1, translayer.BMAILVER2:
		return 8
	default:
		return -1
	}
}

func (h *Header) GetBytes() []byte {
	var r, tmp []byte

	tmp = translayer.UInt16ToBuf(h.Ver)
	r = append(r, tmp...)
	tmp = translayer.UInt16ToBuf(h.MsgTyp)
	r = append(r, tmp...)
	tmp = translayer.UInt32ToBuf(uint32(h.MsgLen))
	r = append(r, tmp...)

	return r

}

func (h *Header) Derive(data []byte) (int, error) {
	if len(data) < int(h.GetLen()) {
		return 0, errors.New("data buffer too small")
	}

	h.Ver = binary.BigEndian.Uint16(data)
	h.MsgTyp = binary.BigEndian.Uint16(data[2:])
	h.MsgLen = int(binary.BigEndian.Uint32(data[4:]))

	return int(h.GetLen()), nil
}

type HELO struct {
}

//ErrCode
//0: success
//1: server not support this version.

type HELOACK struct {
	SN             BMailSN       `json:"sn"`
	SrvBca         bmail.Address `json:"srv"`
	ErrCode        int           `json:"errCode"`
	SupportVersion []uint16      `json:"support_version"`
}

func (ha *HELOACK) MsgType() uint16 {
	return translayer.HELLO_ACK
}

func (ha *HELOACK) VerifyHeader(header *Header) bool {
	return header.MsgLen != 0 && header.MsgTyp == translayer.HELLO_ACK
}

func (ha *HELOACK) GetBytes() ([]byte, error) {
	return json.Marshal(*ha)
}

type EnvelopeSyn struct {
	SN   BMailSN        `json:"sn"`
	Sig  []byte         `json:"sig"`
	Hash []byte         `json:"hash"`
	Env  *BMailEnvelope `json:"env"`
}

func (es *EnvelopeSyn) MsgType() uint16 {
	return translayer.SEND_CRYPT_ENVELOPE
}
func (es *EnvelopeSyn) VerifyHeader(header *Header) bool {
	return header.MsgTyp == translayer.SEND_CRYPT_ENVELOPE &&
		header.MsgLen != 0
}

func (ha *EnvelopeSyn) GetBytes() ([]byte, error) {
	return json.Marshal(*ha)
}

type EnvelopeAck struct {
	NextSN    BMailSN `json:"nextSN"`
	Hash      []byte  `json:"hash"`
	Sig       []byte  `json:"sig"`
	ErrorCode int     `json:"errorCode"`
}

func (ea *EnvelopeAck) MsgType() uint16 {
	return translayer.RESP_CRYPT_ENVELOPE
}
func (ea *EnvelopeAck) VerifyHeader(header *Header) bool {
	return header.MsgTyp == translayer.RESP_CRYPT_ENVELOPE &&
		header.MsgLen != 0
}
func (ha *EnvelopeAck) GetBytes() ([]byte, error) {
	return json.Marshal(*ha)
}
