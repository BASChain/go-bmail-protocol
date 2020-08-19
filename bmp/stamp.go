package bmp

import (
	"encoding/json"
	"github.com/realbmail/go-bmail-protocol/translayer"
)

type StampReceiptSyn struct {
	StampAddr string `json:"stampAddress"`
	UserAddr  string `json:"userAddress"`
}

func (srs *StampReceiptSyn) MsgType() uint16 {
	return translayer.STAMP_RECEIPT_QUERY
}
func (srs *StampReceiptSyn) VerifyHeader(header *Header) bool {
	return true
}
func (srs *StampReceiptSyn) GetBytes() ([]byte, error) {
	return json.Marshal(srs)
}

type StampOptsAck struct {
	IssuerName string   `json:"issuerName"`
	HomePage   string   `json:"homePage"`
	StampAddr  []string `json:"stampAddr"`
}

func (sa *StampOptsAck) MsgType() uint16 {
	return translayer.RESP_STAMP_LIST
}
func (sa *StampOptsAck) VerifyHeader(header *Header) bool {
	return header.MsgTyp == translayer.RESP_STAMP_LIST &&
		header.MsgLen != 0
}
func (sa *StampOptsAck) GetBytes() ([]byte, error) {
	return json.Marshal(*sa)
}

func (sa *StampOptsAck) String() string {
	j, _ := json.Marshal(*sa)
	return string(j)
}

type StampTX struct {
}

type StampReceipt struct {
	StampAddr string `json:"stampAddress"`
	UserAddr  string `json:"userAddress"`
	Credit    int64  `json:"credit"`
}
type StampReceiptACK struct {
	Sig  []byte `json:"signature"`
	Hash []byte `json:"hash"`
	*StampReceipt
}

func (sra *StampReceiptACK) MsgType() uint16 {
	return translayer.RESP_STAMP_RECEIPT
}
func (sra *StampReceiptACK) VerifyHeader(header *Header) bool {
	return header.MsgTyp == translayer.RESP_STAMP_RECEIPT &&
		header.MsgLen != 0
}
func (sra *StampReceiptACK) GetBytes() ([]byte, error) {
	return json.Marshal(sra)
}

func (sra *StampReceiptACK) String() string {
	j, _ := json.Marshal(sra)
	return string(j)
}
