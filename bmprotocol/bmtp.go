package bmprotocol

import (
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"strconv"
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
	Sn    []byte
	NewSn []byte
	EId   translayer.EnveUniqID
	ErrId int
}

func (rse *RespSendCryptEnvelope) String() string {
	s := rse.BMTransLayer.String()
	s += "sn:" + base58.Encode(rse.Sn)
	s += "\r\n"
	s += "NewSn:" + base58.Encode(rse.NewSn)
	s += "\r\n"
	s += base58.Encode(rse.EId[:])
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

	tmp, err := PackShortBytes(rse.Sn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rse.NewSn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rse.EId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return AddPackHead(&(rse.BMTransLayer), r)
}

func (rse *RespSendCryptEnvelope) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	rse.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	rse.NewSn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	copy(rse.EId[:], tmp)

	return offset, nil

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
	Sn    []byte
	NewSn []byte
	EId   translayer.EnveUniqID
	ErrId int
}

func (rse *RespSendEnvelope) String() string {
	s := rse.BMTransLayer.String()
	s += "sn:" + base58.Encode(rse.Sn)
	s += "\r\n"
	s += "NewSn:" + base58.Encode(rse.NewSn)
	s += "\r\n"
	s += base58.Encode(rse.EId[:])
	s += "\r\n"
	s += strconv.Itoa(rse.ErrId)
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

	tmp, err := PackShortBytes(rse.Sn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rse.NewSn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rse.EId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return AddPackHead(&(rse.BMTransLayer), r)
}

func (rse *RespSendEnvelope) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	rse.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	rse.NewSn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	copy(rse.EId[:], tmp)

	return offset, nil

}
