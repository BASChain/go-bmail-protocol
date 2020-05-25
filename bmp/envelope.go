package bmp

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-bmail-account"
	"github.com/google/uuid"
	"time"
)

const (
	BMailModeUnknown = iota
	BMailModeP2P
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
	Date     time.Time     `json:"time"`
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

	encodeSub, err := account.EncryptWithIV(key, iv.Bytes(), ([]byte)(re.Subject))
	if err != nil {
		return nil, err
	}
	encodeMsg, err := account.EncryptWithIV(key, iv.Bytes(), ([]byte)(re.MsgBody))
	if err != nil {
		return nil, err
	}

	obj := &CryptEnvelope{
		EnvelopeHead: re.EnvelopeHead,
		CryptSub:     encodeSub,
		CryptBody:    encodeMsg,
	}
	obj.IV = *iv

	return obj, nil
}

func (re *RawEnvelope) Hash() []byte {
	data, _ := json.Marshal(re)
	hash := sha256.Sum256(data)
	return hash[:]
}
