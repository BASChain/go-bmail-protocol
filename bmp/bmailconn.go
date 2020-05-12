package bmp

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"net"
)

type BMailConn struct {
	*net.TCPConn
}

func NewBMConn(ip net.IP) *BMailConn {
	rAddr := &net.TCPAddr{IP: ip, Port: translayer.BMTP_PORT}
	conn, err := net.DialTCP("tcp4", nil, rAddr)
	if err != nil {
		return nil
	}
	return &BMailConn{conn}
}

func (bc *BMailConn) SendWithHeader(v Envelope) error {
	data, err := v.Pack()
	if err != nil {
		return err
	}

	header := Header{
		Ver:     translayer.BMAILVER1,
		MsgTyp:  v.MsgType(),
		DataLen: uint32(len(data)),
	}

	if _, err := bc.Write(header.Pack()); err != nil {
		return err
	}

	if header.DataLen == 0 {
		return nil
	}

	if _, err := bc.Write(data); err != nil {
		return err
	}

	return nil
}

func (bc *BMailConn) ReadWithHeader(v Envelope) error {
	header := &Header{}
	buf := make([]byte, header.GetLen())

	if _, err := bc.Read(buf); err != nil {
		return err
	}

	if err := header.Unpack(buf); err != nil {
		return err
	}

	if !v.VerifyHeader(header) {
		return fmt.Errorf("unexcept data")
	}

	if header.DataLen == 0 {
		return nil
	}

	buf = make([]byte, header.DataLen)
	if _, err := bc.Read(buf); err != nil {
		return err
	}

	if err := v.UnPack(buf); err != nil {
		return err
	}

	return nil
}
