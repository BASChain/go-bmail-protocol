package bmp

import (
	"encoding/json"
	"fmt"
	"github.com/realbmail/go-bmail-protocol/translayer"
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

func (bc *BMailConn) Helo() error {
	header := Header{
		Ver:    translayer.BMAILVER1,
		MsgTyp: translayer.HELLO,
		MsgLen: 0,
	}
	data := header.GetBytes()
	if _, err := bc.Write(data); err != nil {
		return err
	}

	return nil
}

func (bc *BMailConn) QueryStamp() error {
	header := Header{
		Ver:    translayer.BMAILVER1,
		MsgTyp: translayer.STAMP_QUERY,
		MsgLen: 0,
	}
	data := header.GetBytes()
	if _, err := bc.Write(data); err != nil {
		return err
	}

	return nil
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

	data := header.GetBytes()
	if n, err := bc.Write(data); err != nil {
		fmt.Println("write header len:", n)
		return err
	}
	fmt.Println("send with header: body:=>", string(dataV))
	if n, err := bc.Write(dataV); err != nil {
		fmt.Println("write body len:", n)
		return err
	}
	return nil
}

func (bc *BMailConn) ReadWithHeader(v EnvelopeMsg) error {
	header := &Header{Ver: translayer.BMAILVER1}
	buf := make([]byte, header.GetLen())
	if _, err := bc.Read(buf); err != nil {
		return err
	}
	if _, err := header.Derive(buf); err != nil {
		fmt.Println("header.Derive:", err)
		return err
	}

	if !v.VerifyHeader(header) {
		return fmt.Errorf("unexcept data")
	}

	if header.MsgLen == 0 {
		return nil
	}
	buf = make([]byte, header.MsgLen)
	offset := 0
	for {
		n, err := bc.Read(buf[offset:])
		if err != nil {
			fmt.Println("bc.Read:", err)
			return err
		}
		offset += n
		if offset >= header.MsgLen {
			break
		}
	}

	fmt.Println("read with header: body:=>", string(buf))

	if err := json.Unmarshal(buf, v); err != nil {
		fmt.Println("json.Unmarshal:", err)
		return err
	}
	return nil
}
