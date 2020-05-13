package bmp

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"net"
)

type BMailConn struct {
	*net.TCPConn
}

func NewBMConn(ip net.IP) (*BMailConn, error) {
	rAddr := &net.TCPAddr{IP: ip, Port: translayer.BMTP_PORT}
	conn, err := net.DialTCP("tcp4", nil, rAddr)
	if err != nil {
		return nil, err
	}
	return &BMailConn{conn}, nil
}

func (bc *BMailConn) SendWithHeader(v EnvelopeMsg) error {
	dataV, err := json.Marshal(v)
	if err != nil {
		return err
	}

	header := Header{
		Ver:    translayer.BMAILVER1,
		MsgTyp: v.MsgType(),
		MsgLen: len(dataV),
	}

	data, err := json.Marshal(header)
	if err != nil {
		return err
	}
	if _, err := bc.Write(data); err != nil {
		return err
	}
	if _, err := bc.Write(dataV); err != nil {
		return err
	}
	return nil
}

func (bc *BMailConn) ReadWithHeader(v EnvelopeMsg) error {
	header := &Header{}
	buf := make([]byte, header.GetLen())
	if _, err := bc.Read(buf); err != nil {
		return err
	}
	if err := json.Unmarshal(buf, header); err != nil {
		return err
	}

	if !v.VerifyHeader(header) {
		return fmt.Errorf("unexcept data")
	}

	if header.MsgLen == 0 {
		return nil
	}
	buf = make([]byte, header.MsgLen)
	if _, err := bc.Read(buf); err != nil {
		return err
	}

	if err := json.Unmarshal(buf, v); err != nil {
		return err
	}
	return nil
}
