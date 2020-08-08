package bpop

import (
	"encoding/json"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-bmail-protocol/translayer"
)

type Command interface {
	Hash() []byte
	MsgType() uint16
}

type CommandSyn struct {
	SN  bmp.BMailSN `json:"sn"`
	Sig []byte      `json:"sig"`
	Cmd Command     `json:"cmd"`
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

type CommandContent interface {
	MsgType() uint16
	Hash() []byte
}

const (
	EC_Success int = iota
	EC_No_Mail
)

type CommandAck struct {
	NextSN    bmp.BMailSN    `json:"next_sn"`
	Hash      []byte         `json:"hash"`
	Sig       []byte         `json:"sig"`
	ErrorCode int            `json:"error_code"`
	CmdCxt    CommandContent `json:"cmd"`
}

func (cs *CommandAck) MsgType() uint16 {
	return cs.CmdCxt.MsgType()
}

func (cs *CommandAck) GetBytes() ([]byte, error) {
	return json.Marshal(*cs)
}

func (cs *CommandAck) VerifyHeader(header *bmp.Header) bool {
	return header.MsgTyp == translayer.RETR_RESP &&
		header.MsgLen != 0
}
