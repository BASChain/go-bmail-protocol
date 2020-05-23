package bpop

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/google/uuid"
)

type CmdDownload struct {
	Owner      bmail.Address `json:"owner"`
	MailCnt    int           `json:"mail_cnt"`
	BeforeTime int64         `json:"befor_time"`
}

func (cd *CmdDownload) Hash() []byte {
	data, _ := json.Marshal(*cd)

	return (sha256.Sum256(data))[:]
}

func (cd *CmdDownload) MsgType() uint16 {
	return translayer.RETR
}

type CmdState struct {
	Owner     bmail.Address `json:"owner"`
	BeforTime int64         `json:"befor_time"`
}

func (cs *CmdState) Hash() []byte {
	data, _ := json.Marshal(*cs)

	return (sha256.Sum256(data))[:]
}

func (cs *CmdState) MsgType() uint16 {
	return translayer.STAT
}

type CmdDelete struct {
	Owner bmail.Address `json:"owner"`
	Eids  []uuid.UUID   `json:"eid"`
}

func (cd *CmdDelete) Hash() []byte {
	data, _ := json.Marshal(*cd)

	return (sha256.Sum256(data))[:]
}

func (cs *CmdDelete) MsgType() uint16 {
	return translayer.DELETE
}

func (cs *CommandSyn) MsgType() uint16 {
	return cs.Cmd.MsgType()
}

func (cs *CommandSyn) GetBytes() ([]byte, error) {
	return json.Marshal(*cs)
}

func (cs *CommandSyn) VerifyHeader(header *bmp.Header) bool {
	return header.MsgTyp == cs.Cmd.MsgType() &&
		header.MsgLen != 0
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
	return (sha256.Sum256(data))[:]
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
	return (sha256.Sum256(data))[:]
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

type CmdDeleteAclk struct {
	Result []CmdResult `json:"result"`
}

func (cda *CmdDeleteAclk) MsgType() uint16 {
	return translayer.DELETE_RESP
}

func (cda *CmdDeleteAclk) GetBytes() ([]byte, error) {
	return json.Marshal(*cda)
}

func (cda *CmdDeleteAclk) Hash() []byte {
	data, _ := json.Marshal(*cda)
	return (sha256.Sum256(data))[:]
}
