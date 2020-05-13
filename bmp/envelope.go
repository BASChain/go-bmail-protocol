package bmp

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-bmail-account"
	"github.com/google/uuid"
)

const (
	BMailModeP2P = iota
	BMailModeP2S

	RcpTypeTo = iota
	RcpTypeCC
	RcpTypeBcc
)

type Envelope interface {
	Hash() []byte
}

type EnvelopeHead struct {
	Eid      uuid.UUID     `json:"eid"`
	From     string        `json:"from"`
	FromAddr bmail.Address `json:"fromAddr"`
	To       string        `json:"to"`
	ToAddr   bmail.Address `json:"toAddr"`
	IV       BMailIV       `json:"iv"`
}
type EnvelopeBody struct {
	Subject string `json:"subject"`
	MsgBody string `json:"msgBody"`
}

type RawEnvelope struct {
	EnvelopeHead
	EnvelopeBody
}

func (re *RawEnvelope) Seal(key []byte) (Envelope, error) {
	iv, err := NewIV()
	if err != nil {
		return nil, err
	}

	cc := &EnvelopeBody{
		Subject: re.Subject,
		MsgBody: re.MsgBody,
	}

	ccData, err := json.Marshal(cc)
	if err != nil {
		return nil, err
	}

	encoded, err := account.EncryptWithIV(key, iv.Bytes(), ccData)
	if err != nil {
		return nil, err
	}

	obj := &CryptEnvelope{
		EnvelopeHead: re.EnvelopeHead,
		CryptData:    encoded,
	}
	obj.IV = *iv

	return obj, nil
}

func (re *RawEnvelope) Hash() []byte {
	data, _ := json.Marshal(re)
	hash := sha256.Sum256(data)
	return hash[:]
}
