package bpop

import "github.com/BASChain/go-bmail-protocol/bmp"

type Command interface {
	Hash() []byte
	MsgType() uint16
}

type CommandSyn struct {
	SN  bmp.BMailSN `json:"sn"`
	Sig []byte      `json:"sig"`
	Cmd Command     `json:"cmd"`
}

type CommandContent interface {
	MsgType() uint16
	Hash() []byte
	GetBytes() ([]byte, error)
}

type CommandAck struct {
	NextSN    bmp.BMailSN    `json:"next_sn"`
	Hash      []byte         `json:"hash"`
	Sig       []byte         `json:"sig"`
	ErrorCode int            `json:"error_code"`
	CmdCnt    CommandContent `json:"cmd"`
}
