package bmprotocol

import (
	"github.com/realbmail/go-bmail-protocol/translayer"
)

//client -> server
type SendCryptEnvelope struct {
	translayer.BMTransLayer
	CryptEnvelope
}

func (se *SendCryptEnvelope) String() string {
	s := se.BMTransLayer.String()
	s += se.CryptEnvelope.String()

	return s
}

func NewSendCryptEnvelope() *SendCryptEnvelope {
	se := &SendCryptEnvelope{}

	bmact := translayer.NewBMTL(translayer.SEND_CRYPT_ENVELOPE)
	se.BMTransLayer = *bmact

	return se
}

func (se *SendCryptEnvelope) Pack() ([]byte, error) {

	r := NewHeadBuf()

	tmp, err := se.CryptEnvelope.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return AddPackHead(&(se.BMTransLayer), r)
}

func (se *SendCryptEnvelope) UnPack(data []byte) (int, error) {

	ce := &se.CryptEnvelope

	return ce.UnPack(data)
}

/*
ErrId:
0: success
1: Receipt Address not found
2: Receipt Address is forbidden
3: Data is forbidden
4: bmail command Not Found
5: Mail to large
6: Mail refused by server
7: Connect to Server failed

*/
//server -> client
type RespSendCryptEnvelope struct {
	translayer.BMTransLayer
	ConfirmEnvelope
}

func (rse *RespSendCryptEnvelope) String() string {
	s := rse.BMTransLayer.String()
	s += rse.ConfirmEnvelope.String()

	return s
}

func NewRespSendCryptEnvelope() *RespSendCryptEnvelope {
	rse := &RespSendCryptEnvelope{}

	bmact := translayer.NewBMTL(translayer.RESP_CRYPT_ENVELOPE)

	rse.BMTransLayer = *bmact

	return rse
}

func (rse *RespSendCryptEnvelope) Pack() ([]byte, error) {

	r := NewHeadBuf()

	ce := &rse.ConfirmEnvelope

	tmp, err := ce.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return AddPackHead(&(rse.BMTransLayer), r)
}

func (rse *RespSendCryptEnvelope) UnPack(data []byte) (int, error) {

	ce := &rse.ConfirmEnvelope

	return ce.UnPack(data)
}

type SendEnvelope struct {
	translayer.BMTransLayer
	Envelope
}

func (se *SendEnvelope) String() string {
	s := se.BMTransLayer.String()
	s += se.Envelope.String()

	return s
}

func NewSendEnvelope() *SendEnvelope {
	se := &SendEnvelope{}

	bmact := translayer.NewBMTL(translayer.SEND_ENVELOPE)
	se.BMTransLayer = *bmact

	return se
}

func (se *SendEnvelope) Pack() ([]byte, error) {

	r := NewHeadBuf()

	tmp, err := se.Envelope.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return AddPackHead(&(se.BMTransLayer), r)
}

func (se *SendEnvelope) UnPack(data []byte) (int, error) {

	ce := &se.Envelope

	return ce.UnPack(data)
}

//server -> client
type RespSendEnvelope struct {
	translayer.BMTransLayer
	ConfirmEnvelope
}

func (rse *RespSendEnvelope) String() string {
	s := rse.BMTransLayer.String()
	s += rse.ConfirmEnvelope.String()

	return s
}

func NewRespSendEnvelope() *RespSendEnvelope {
	rse := &RespSendEnvelope{}

	bmact := translayer.NewBMTL(translayer.RESP_ENVELOPE)

	rse.BMTransLayer = *bmact

	return rse
}

func (rse *RespSendEnvelope) Pack() ([]byte, error) {

	r := NewHeadBuf()

	ce := &rse.ConfirmEnvelope

	tmp, err := ce.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return AddPackHead(&(rse.BMTransLayer), r)
}

func (rse *RespSendEnvelope) UnPack(data []byte) (int, error) {
	ce := &rse.ConfirmEnvelope

	return ce.UnPack(data)

}
