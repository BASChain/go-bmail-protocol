package bmprotocol

import (
	"encoding/binary"
	"errors"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
)

//client -> server
type SendEnvelope struct {
	translayer.BMTransLayer
	CryptEnvelope
}

func (se *SendEnvelope) String() string {
	s := se.BMTransLayer.HeadString()
	s += se.CryptEnvelope.String()

	return s
}

func NewSendEnvelope() *SendEnvelope {
	se := &SendEnvelope{}

	bmact := translayer.NewBMTL(translayer.SEND_ENVELOPE, nil)
	se.BMTransLayer = *bmact

	return se
}

func (se *SendEnvelope) Pack() ([]byte, error) {

	if len(se.CipherTxt) == 0 {
		return nil, errors.New("Cipher Text is nil")
	}

	var r []byte

	rh, err := se.EnvelopeHead.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, rh...)

	var ret []byte
	ret, err = se.CryptEnvelope.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, ret...)

	se.BMTransLayer.SetData(r)

	return se.BMTransLayer.Pack()
}

func (se *SendEnvelope) UnPack(data []byte) (int, error) {
	eh := &se.EnvelopeHead
	offset, err := eh.UnPack(data)

	if err != nil {
		return 0, err
	}

	ce := &se.CryptEnvelope
	var let int

	let, err = ce.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += let

	return offset, nil
}

//server -> client
type RespSendEnvelope struct {
	translayer.BMTransLayer
	EnvelopeHead
	IV []byte //same as SendEnvelope
}

func (rse *RespSendEnvelope) String() string {
	s := rse.BMTransLayer.HeadString()
	eh := &rse.EnvelopeHead
	s += eh.String()
	s += "\r\n"
	s += base58.Encode(rse.IV)
	return s
}

func NewRespSendEnvelope() *RespSendEnvelope {
	rse := &RespSendEnvelope{}

	bmact := translayer.NewBMTL(translayer.RESP_ENVELOPE, nil)

	rse.BMTransLayer = *bmact

	return rse
}

func (rse *RespSendEnvelope) Pack() ([]byte, error) {

	if len(rse.IV) == 0 {
		return nil, errors.New("IV is not set")
	}

	var (
		r, tmp []byte
		err    error
	)

	eh := &rse.EnvelopeHead
	tmp, err = eh.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	var iv []byte
	iv, err = PackShortBytes(rse.IV)
	if err != nil {
		return nil, err
	}
	r = append(r, iv...)

	rse.BMTransLayer.SetData(r)

	return rse.BMTransLayer.Pack()
}

func (rse *RespSendEnvelope) UnPack(data []byte) (int, error) {
	var (
		l   int
		err error
	)
	offset := 0
	eh := &rse.EnvelopeHead

	l, err = eh.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	rse.IV, l, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += l

	return offset, nil

}

//server -> client
type SendEnvelopeFail struct {
	translayer.BMTransLayer
	CryptEnvelope
	ErrorCode int
}

func (sef *SendEnvelopeFail) String() string {
	s := sef.BMTransLayer.HeadString()
	s += sef.CryptEnvelope.String()
	s += "\r\n"
	s += strconv.Itoa(sef.ErrorCode)

	return s
}

func NewSendEnvelopeFail() *SendEnvelopeFail {
	sef := &SendEnvelopeFail{}

	bmact := translayer.NewBMTL(translayer.SEND_ENVELOPE_FAILED, nil)

	sef.BMTransLayer = *bmact

	return sef
}

func (sef *SendEnvelopeFail) Pack() ([]byte, error) {
	if len(sef.CipherTxt) == 0 {
		return nil, errors.New("Cipher Text is not set")
	}
	var (
		r, tmp []byte
		err    error
	)

	ce := &sef.CryptEnvelope

	tmp, err = ce.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	bufl := translayer.UInt32ToBuf(uint32(sef.ErrorCode))

	r = append(r, bufl...)

	sef.BMTransLayer.SetData(r)

	return sef.BMTransLayer.Pack()
}

func (sef *SendEnvelopeFail) UnPack(data []byte) (int, error) {
	var (
		offset, l int
		err       error
	)

	ec := &sef.CryptEnvelope
	l, err = ec.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	sef.ErrorCode = int(binary.BigEndian.Uint32(data[offset:]))

	offset += translayer.Uint32Size

	return offset, nil

}

//client -> server
type RespSendEnvelopeFail struct {
	translayer.BMTransLayer
	EnvelopeHead
	IV []byte
}

func (rsef *RespSendEnvelopeFail) String() string {
	s := rsef.BMTransLayer.HeadString()
	eh := &rsef.EnvelopeHead
	s += eh.String()
	s += "\r\n"
	s += base58.Encode(rsef.IV)
	return s
}

func NewRespSendEnvelopeFail() *RespSendEnvelopeFail {
	rsef := &RespSendEnvelopeFail{}
	bmact := translayer.NewBMTL(translayer.RESP_SEND_ENVELOPE_FAILED, nil)

	rsef.BMTransLayer = *bmact

	return rsef
}

func (rsef *RespSendEnvelopeFail) Pack() ([]byte, error) {
	if len(rsef.IV) == 0 {
		return nil, errors.New("IV is not set")
	}

	var (
		r, tmp []byte
		err    error
	)

	eh := &rsef.EnvelopeHead
	tmp, err = eh.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	tmp, err = PackShortBytes(rsef.IV)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	rsef.BMTransLayer.SetData(r)

	return rsef.BMTransLayer.Pack()
}

func (rsef *RespSendEnvelopeFail) UnPack(data []byte) (int, error) {

	var (
		offset, l int
		err       error
	)
	eh := &rsef.EnvelopeHead

	l, err = eh.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	rsef.IV, l, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	return offset, nil

}
