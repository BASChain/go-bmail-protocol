package bmp

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/BASChain/go-bmail-account"
)

const (
	BMailModeUnknown = iota
	BMailModeP2P
	BMailModeP2S

	RcpTypeTo = iota
	RcpTypeCC
	RcpTypeBcc
	RcpMonitor
)

type Recipient struct {
	ToName   string        `json:"to"`
	ToAddr   bmail.Address `json:"toAddr"`
	RcptType int8          `json:"rcptType"`
	AESKey   []byte        `json:"aesKey"`
}

type BMailEnvelope struct {
	Eid           string        `json:"eid"`
	From          string        `json:"from"`
	FromAddr      bmail.Address `json:"fromAddr"`
	RCPTs         []*Recipient  `json:"rcpts"`
	DateSince1970 uint64        `json:"timeSince1970"`
	Subject       string        `json:"subject"`
	MailBody      string        `json:"mailBody"`
	PreEid        string        `json:"preEid"`
}

func (re *BMailEnvelope) Hash() []byte {
	data, _ := json.Marshal(re)
	hash := sha256.Sum256(data)
	return hash[:]
}
