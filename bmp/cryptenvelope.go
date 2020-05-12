package bmp

import (
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"github.com/google/uuid"
)

type CryptEnvelope struct {
	IV           BMailIV       `json:"iv"`
	EId          uuid.UUID     `json:"eid"`
	From         bmail.Address `json:"from"`
	PeerAddr     bmail.Address `json:"peer"`
	PeerMailName string        `json:"peerName"`
	CryptData    []byte        `json:"cryptBody"`
}

func (cb *CryptEnvelope) Pack() ([]byte, error) {

	_, err := bmprotocol.PackShortBytes(cb.IV.Bytes())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

type CryptContent struct {
	Subject []byte `json:"subject"`
	MsgBody []byte `json:"msgBody"`
}

func (cc *CryptContent) Pack() ([]byte, error) {
	data, err := bmprotocol.PackShortBytes(cc.Subject)
	if err != nil {
		return nil, err
	}

	data2, err := bmprotocol.PackShortBytes(cc.MsgBody)
	if err != nil {
		return nil, err
	}

	data = append(data, data2...)
	return data, nil
}
