package bmp

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/realbmail/go-bmail-account"
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

func (r *Recipient) ToString() string {
	return fmt.Sprintf("\n=================================="+
		"\n\tToName:\t%20s"+
		"\n\tToAddr:\t%20s"+
		"\n\tRcptType:\t%20d"+
		"\n\tAESKey:\t%20x"+
		"\n==================================",
		r.ToName,
		r.ToAddr,
		r.RcptType,
		r.AESKey)
}

type BMailEnvelope struct {
	Eid           string        `json:"eid"`
	FromName      string        `json:"fromName"`
	FromAddr      bmail.Address `json:"fromAddr"`
	RCPTs         []*Recipient  `json:"rcpts"`
	DateSince1970 uint64        `json:"timeSince1970"`
	Subject       string        `json:"subject"`
	MailBody      string        `json:"mailBody"`
	SessionID     string        `json:"sessionID"`
}

func (re *BMailEnvelope) Hash() []byte {
	data, _ := json.Marshal(re)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (re *BMailEnvelope) ToString() string {

	str := fmt.Sprintf("\n======================BMailEnvelope========================"+
		"\n\tEid:\t%20s"+
		"\n\tFrom:\t%20s"+
		"\n\tFromAddr:\t%20s"+
		"\n\tSessionEid:\t%20s",
		re.Eid,
		re.FromName,
		re.FromAddr,
		re.SessionID)

	for _, r := range re.RCPTs {
		str += r.ToString()
	}

	str += "\n==========================================================="
	return str
}
