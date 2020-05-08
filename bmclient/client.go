package bmclient

import (
	"bytes"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/pkg/errors"
	"net"
	"strconv"
	"time"
)

type BMClient struct {
	sn      []byte
	c       *net.TCPConn
	timeout int //second
}

func NewClient(serverIP net.IP, timeout int) *BMClient {
	laddr := &net.TCPAddr{}
	raddr := &net.TCPAddr{IP: serverIP, Port: 1025}
	conn, err := net.DialTCP("tcp4", laddr, raddr)
	if err != nil {
		return nil
	}

	c := &BMClient{}
	c.c = conn
	c.timeout = timeout

	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	return c
}

func (c *BMClient) Close() {
	c.c.Close()
	c.c = nil
	c.sn = nil
}

func (c *BMClient) GetSn() []byte {
	return c.sn
}

func (c *BMClient) SendEnvelope(envelope *bmprotocol.SendEnvelope) (rse *bmprotocol.RespSendEnvelope, err error) {
	if c.c == nil {
		return nil, errors.New("client is not initialized")
	}
	if bytes.Compare(envelope.EnvelopeSig.Sn, c.sn) != 0 {
		return nil, errors.New("Envelope not correct")
	}

	var data []byte

	data, err = envelope.Pack()
	if err != nil {
		return nil, err
	}

	var n int
	n, err = c.c.Write(data)
	if n != len(data) || err != nil {
		return nil, errors.New("Send envelope Failed")
	}
	buf := make([]byte, translayer.BMHeadSize())

	n, err = c.c.Read(buf)
	if n != translayer.BMHeadSize() || err != nil {
		return nil, errors.New("Read a bad bmail head")
	}

	bmtl := &translayer.BMTransLayer{}
	bmtl.UnPack(buf)

	if bmtl.GetMsgType() != translayer.RESP_ENVELOPE || bmtl.GetDataLen() == 0 {
		return nil, errors.New("Received a error message: " + strconv.Itoa(int(bmtl.GetMsgType())))
	}

	buf = make([]byte, bmtl.GetDataLen())

	n, err = c.c.Read(buf)
	if n != int(bmtl.GetDataLen()) || err != nil {
		return nil, errors.New("Read a bad bmail data")
	}

	resp := &bmprotocol.RespSendEnvelope{}
	resp.BMTransLayer = *bmtl
	_, err = resp.UnPack(buf)
	if err != nil {
		return nil, err
	}
	c.sn = resp.NewSn

	return resp, nil

}

func (c *BMClient) HeloSendAndRcv() (err error) {

	if c.c == nil {
		return errors.New("client is not initialized")
	}

	helo := bmprotocol.NewBMHello()
	data, _ := helo.Pack()

	var n int

	n, err = c.c.Write(data)
	if n != len(data) || err != nil {
		return errors.New("Send Helo Failed")
	}

	buf := make([]byte, translayer.BMHeadSize())

	n, err = c.c.Read(buf)
	if n != translayer.BMHeadSize() || err != nil {
		return errors.New("Read a bad bmail head")
	}

	bmtl := &translayer.BMTransLayer{}
	bmtl.UnPack(buf)

	if bmtl.GetMsgType() != translayer.HELLO_ACK || bmtl.GetDataLen() == 0 {
		return errors.New("Received a error message: " + strconv.Itoa(int(bmtl.GetMsgType())))
	}

	//read left
	buf = make([]byte, bmtl.GetDataLen())

	n, err = c.c.Read(buf)
	if n != int(bmtl.GetDataLen()) || err != nil {
		return errors.New("Read a bad bmail data")
	}

	ha := &bmprotocol.BMHelloACK{}
	ha.BMTransLayer = *bmtl
	_, err = ha.UnPack(buf)
	if err != nil {
		return err
	}

	c.sn = ha.GetSn()

	return nil
}
