package bpop

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/google/uuid"
)

const (
	DefaultMailCount int = 20
	DirectionToLeft bool = true
	DirectionToRight bool = false
)


type CmdDownload struct {
	MailAddr  string        `json:"mail_addr"`
	Owner     bmail.Address `json:"owner"`
	MailCnt   int           `json:"mail_cnt"`
	Direction bool          `json:"direction"` //false -> after TimePivot, true -> before TimePivot
	TimePivot int64         `json:"time_pivot"`
}

func (cd *CmdDownload) Hash() []byte {
	data, _ := json.Marshal(*cd)

	hash := sha256.Sum256(data)

	return hash[:]
}

func (cd *CmdDownload) MsgType() uint16 {
	return translayer.RETR
}

type CmdState struct {
	MailAddr  string        `json:"mail_addr"`
	Owner     bmail.Address `json:"owner"`
	BeforTime int64         `json:"before_time"`
}

func (cs *CmdState) Hash() []byte {
	data, _ := json.Marshal(*cs)

	hash := sha256.Sum256(data)

	return hash[:]
}

func (cs *CmdState) MsgType() uint16 {
	return translayer.STAT
}

type CmdDelete struct {
	MailAddr string        `json:"mail_addr"`
	Owner    bmail.Address `json:"owner"`
	Eids     []uuid.UUID   `json:"eid"`
}

func (cd *CmdDelete) Hash() []byte {
	data, _ := json.Marshal(*cd)

	hash := sha256.Sum256(data)

	return hash[:]
}

func (cs *CmdDelete) MsgType() uint16 {
	return translayer.DELETE
}

type CmdDownloadAck struct {
	CryptEps []bmp.CryptEnvelope
}

func (cda *CmdDownloadAck) MsgType() uint16 {
	return translayer.RETR_RESP
}

func (cda *CmdDownloadAck) GetBytes() ([]byte, error) {
	return json.Marshal(*cda)
}

func (cda *CmdDownloadAck) Hash() []byte {
	data, _ := json.Marshal(*cda)

	hash := sha256.Sum256(data)

	return hash[:]
}

type State struct {
	TotalSpace int64 `json:"total_space"`
	UsedSize   int64 `json:"used_size"`
	TotalCount int   `json:"total_count"`
}

type CmdStateAck struct {
	SendMail    State `json:"send_mail_space"`
	ReceiptMail State `json:"receipt_mail"`
}

func (csa *CmdStateAck) MsgType() uint16 {
	return translayer.STAT_RESP
}

func (csa *CmdStateAck) GetBytes() ([]byte, error) {
	return json.Marshal(*csa)
}

func (csa *CmdStateAck) Hash() []byte {
	data, _ := json.Marshal(*csa)

	hash := sha256.Sum256(data)

	return hash[:]
}

//Result
const (
	MailDeleteSuccess int = iota
	MailNotFound
	MailDeleteFailed
)

type CmdResult struct {
	Eid    uuid.UUID `json:"eid"`
	Result int       `json:"result"`
}

type CmdDeleteAck struct {
	Result []CmdResult `json:"result"`
}

func (cda *CmdDeleteAck) MsgType() uint16 {
	return translayer.DELETE_RESP
}

func (cda *CmdDeleteAck) GetBytes() ([]byte, error) {
	return json.Marshal(*cda)
}

func (cda *CmdDeleteAck) Hash() []byte {
	data, _ := json.Marshal(*cda)

	hash := sha256.Sum256(data)

	return hash[:]
}
