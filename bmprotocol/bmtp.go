package bmprotocol

import (
	"encoding/binary"
	"errors"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"strconv"
)

const (
	PeerUnreachable int = iota + 1
	AddressUnavailable
)

const (
	ErrMsg_PeerUnreachable    string = "Peer is unreachable"
	ErrMsg_AddressUnavailable string = "Recipient is not available"
)

type EnvelopeHead struct {
	From     string
	RecpAddr string //recipient
	LPubKey  []byte //local public key
}

func (eh *EnvelopeHead) Pack() ([]byte, error) {
	if eh.From == "" || eh.RecpAddr == "" || len(eh.LPubKey) == 0 {
		return nil, errors.New("")
	}

	var r []byte

	bufl := translayer.UInt16ToBuf(uint16(len(eh.From)))
	r = append(r, bufl...)
	r = append(r, []byte(eh.From)...)

	bufl = translayer.UInt16ToBuf(uint16(len(eh.RecpAddr)))
	r = append(r, bufl...)
	r = append(r, []byte(eh.RecpAddr)...)

	bufl = translayer.UInt16ToBuf(uint16(len(eh.LPubKey)))
	r = append(r, bufl...)
	r = append(r, []byte(eh.LPubKey)...)

	return r, nil
}

func (eh *EnvelopeHead) UnPack(data []byte) (int, error) {
	if len(data) < translayer.Uint16Size {
		return 0, errors.New("Unpack From Head Failed")
	}

	offset := 0

	lfrom := binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	if len(data) < offset+int(lfrom) || lfrom == 0 {
		return 0, errors.New("Unpack From Failed")
	}

	eh.From = string(data[offset : offset+int(lfrom)])

	offset += int(lfrom)

	if len(data) < offset+translayer.Uint16Size {
		return 0, errors.New("Unpack RecpAddr Head Failed")
	}

	lrcp := binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	if len(data) < offset+int(lrcp) || lrcp == 0 {
		return 0, errors.New("Unpack RecpAddr Failed")
	}

	eh.RecpAddr = string(data[offset : offset+int(lrcp)])

	offset += int(lrcp)

	if len(data) < offset+translayer.Uint16Size {
		return 0, errors.New("Unpack LPubKey Head Failed")
	}

	lpk := binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	if len(data) < offset+int(lpk) || lpk == 0 {
		return 0, errors.New("Unpack RecpAddr Failed")
	}

	eh.LPubKey = data[offset : offset+int(lpk)]

	return offset + int(lpk), nil
}

type EnvelopeContent struct {
	To      []string
	CC      []string
	BC      []string
	Subject string
	Data    string
}

func PackShortString(s string) ([]byte, error) {

	return PackShortBytes([]byte(s))
}

func PackShortBytes(data []byte) ([]byte, error) {
	var r []byte

	bufl := translayer.UInt16ToBuf(uint16(len(data)))

	r = append(r, bufl...)

	if len(data) > 0 {
		r = append(r, data...)
	}

	return r, nil
}

func UnPackShortString(data []byte) (string, int, error) {
	bd, of, err := UnPackShortBytes(data)

	s := ""

	if err == nil {
		s = string(bd)
	}

	return s, of, err

}

func UnPackShortBytes(data []byte) ([]byte, int, error) {
	if len(data) < translayer.Uint16Size {
		return nil, 0, errors.New("Unpack Short String Failed")
	}

	offset := 0

	l := binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	if len(data) < offset+int(l) {
		return nil, 0, errors.New("Unpack Short String Failed")
	}

	var s []byte

	if l > 0 {
		s = data[offset : offset+int(l)]
	}

	offset += int(l)

	return s, offset, nil
}

func PackLongString(s string) ([]byte, error) {

	return PackLongBytes([]byte(s))
}

func PackLongBytes(data []byte) ([]byte, error) {

	var r []byte

	bufl := translayer.UInt32ToBuf(uint32(len(data)))

	r = append(r, bufl...)

	if len(data) > 0 {
		r = append(r, data...)
	}

	return r, nil
}

func UnPackLongString(data []byte) (string, int, error) {
	bd, of, err := UnPackLongBytes(data)

	s := ""

	if err == nil {
		s = string(bd)
	}

	return s, of, err
}

func UnPackLongBytes(data []byte) ([]byte, int, error) {
	if len(data) < translayer.Uint32Size {
		return nil, 0, errors.New("Unpack Long String Failed")
	}

	offset := 0

	l := binary.BigEndian.Uint32(data[offset:])
	offset += translayer.Uint32Size

	if len(data) < offset+int(l) {
		return nil, 0, errors.New("Unpack Long String Failed")
	}

	var s []byte

	if l > 0 {
		s = data[offset : offset+int(l)]
	}

	offset += int(l)

	return s, offset, nil
}

func PackShortStringArray(arrs []string) ([]byte, error) {
	if len(arrs) == 0 {
		return nil, errors.New("Pack string Array Failed")
	}

	var r []byte

	bufl := translayer.UInt32ToBuf(uint32(len(arrs)))

	r = append(r, bufl...)

	for i := 0; i < len(arrs); i++ {
		s := arrs[i]

		rs, _ := PackShortString(s)
		if rs != nil {
			r = append(r, rs...)
		}
	}

	return r, nil
}

func UnPackShortStringArray(data []byte) ([]string, int, error) {
	if len(data) < translayer.Uint32Size {
		return nil, 0, errors.New("Unpack Short String Array Failed")
	}

	var rs []string

	offset := 0

	cnt := binary.BigEndian.Uint32(data[offset:])
	offset += translayer.Uint32Size
	if cnt == 0 {
		return rs, offset, nil
	}

	for i := 0; i < int(cnt); i++ {
		s, of1, e := UnPackShortString(data[offset:])
		if e != nil {
			return nil, 0, errors.New("Unpack short string array :" + strconv.Itoa(i) + " Failed")
		}

		offset += of1

		rs = append(rs, s)
	}

	return rs, offset, nil
}

func PackLongStringArray(arrs []string) ([]byte, error) {
	if len(arrs) == 0 {
		return nil, errors.New("Pack string Array Failed")
	}

	var r []byte

	bufl := translayer.UInt16ToBuf(uint16(len(arrs)))

	r = append(r, bufl...)

	for i := 0; i < len(arrs); i++ {
		s := arrs[i]

		rs, _ := PackLongString(s)
		if rs != nil {
			r = append(r, rs...)
		}
	}

	return r, nil
}

func UnPackLongStringArray(data []byte) ([]string, int, error) {
	if len(data) < translayer.Uint32Size {
		return nil, 0, errors.New("Unpack Long String Array Failed")
	}

	var rs []string

	offset := 0

	cnt := binary.BigEndian.Uint32(data[offset:])
	offset += translayer.Uint32Size
	if cnt == 0 {
		return rs, offset, nil
	}

	for i := 0; i < int(cnt); i++ {
		s, of1, e := UnPackLongString(data[offset:])
		if e != nil {
			return nil, 0, errors.New("Unpack long string array :" + strconv.Itoa(i) + " Failed")
		}

		offset += of1

		rs = append(rs, s)
	}

	return rs, offset, nil
}

func (ec *EnvelopeContent) Pack() ([]byte, error) {

	if len(ec.To) == 0 || len(ec.Subject) == 0 {
		return nil, errors.New("Envelope Must have TO Address and Subject")
	}

	var r []byte

	rto, err := PackShortStringArray(ec.To)
	if err != nil {
		return nil, errors.New("Pack Envelope Content TO addres error")
	}

	r = append(r, rto...)

	rcc, err := PackShortStringArray(ec.CC)
	if err != nil {
		return nil, errors.New("Pack Envelope Content CC addres error")
	}

	r = append(r, rcc...)

	rbc, err := PackShortStringArray(ec.BC)
	if err != nil {
		return nil, errors.New("Pack Envelope Content BC addres error")
	}

	r = append(r, rbc...)

	rsubj, err := PackLongString(ec.Subject)
	if err != nil {
		return nil, errors.New("Pack Envelope Content Subject error")
	}

	r = append(r, rsubj...)

	rdata, err := PackLongString(ec.Data)
	if err != nil {
		return nil, errors.New("Pack Envelope Content Data error")
	}

	r = append(r, rdata...)

	return r, nil
}

func (ec *EnvelopeContent) UnPack(data []byte) (int, error) {
	offset := 0
	to, of, err := UnPackShortStringArray(data[offset:])
	if err != nil {
		return 0, err
	}

	ec.To = to
	offset += of

	var cc []string
	cc, of, err = UnPackShortStringArray(data[offset:])
	if err != nil {
		return 0, err
	}

	ec.CC = cc
	offset += of

	var bc []string
	bc, of, err = UnPackShortStringArray(data[offset:])
	if err != nil {
		return 0, err
	}

	ec.BC = bc
	offset += of

	var sbj string
	sbj, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	ec.Subject = sbj
	offset += of

	var ed string
	ed, of, err = UnPackShortString(data[offset:])

	ec.Data = ed
	offset += of

	return offset, nil
}

type EnvelopeTail struct {
	IV  []byte //sn from Bhello
	Sig []byte //signature
}

func (ee *EnvelopeTail) Pack() ([]byte, error) {
	if len(ee.IV) == 0 || len(ee.Sig) == 0 {
		return nil, errors.New("Not a Correct Envelope Tail")
	}

	var r []byte

	riv, err := PackShortBytes(ee.IV)
	if err != nil {
		return nil, errors.New("Pack Envelope Tail IV error")
	}
	r = append(r, riv...)

	var rsig []byte
	rsig, err = PackShortBytes(ee.Sig)
	if err != nil {
		return nil, errors.New("Pack Envelope Tail Signature error")
	}

	r = append(r, rsig...)

	return r, nil
}

func (ee *EnvelopeTail) UnPack(data []byte) (int, error) {
	offset := 0
	var of int
	var err error
	ee.IV, of, err = UnPackShortBytes(data[offset:])

	if err != nil {
		return 0, err
	}

	offset += of

	ee.Sig, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil
}

//client -> server
type SendEnvelope struct {
	translayer.BMTransLayer
	EnvelopeHead
	CipherTxt []byte //crypt from EnvelopeContent
	EnvelopeTail
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

	var rc []byte
	rc, err = PackLongBytes(se.CipherTxt)
	if err != nil {
		return nil, err
	}

	r = append(r, rc...)

	var ret []byte
	ret, err = se.EnvelopeTail.Pack()
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

	var lc int
	se.CipherTxt, lc, err = UnPackLongBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += lc

	et := &se.EnvelopeTail
	var let int

	let, err = et.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += let

	return offset, nil
}

//server -> client
type RespSendEnvelope struct {
	translayer.BMTransLayer
	From     string
	RecpAddr string //recipient
	LAddr    []byte //local public key
	IV       []byte //same as SendEnvelope
}

func NewRespSendEnvelope() *RespSendEnvelope {
	rse := &RespSendEnvelope{}

	bmact := translayer.NewBMTL(translayer.RESP_ENVELOPE, nil)

	rse.BMTransLayer = *bmact

	return rse
}

func (rse *RespSendEnvelope) Pack() ([]byte, error) {

	if rse.From == "" || rse.RecpAddr == "" || len(rse.LAddr) == 0 || len(rse.IV) == 0 {
		return nil, errors.New("Response to Send Envelope Member no value")
	}

	var r []byte
	from, err := PackShortString(rse.From)
	if err != nil {
		return nil, err
	}
	r = append(r, from...)

	var recp []byte
	recp, err = PackShortString(rse.RecpAddr)
	if err != nil {
		return nil, err
	}
	r = append(r, recp...)

	var pub []byte
	pub, err = PackShortBytes(rse.LAddr)
	if err != nil {
		return nil, err
	}
	r = append(r, pub...)

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
	rse.From, l, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	rse.RecpAddr, l, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	rse.LAddr, l, err = UnPackShortBytes(data[offset:])
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
	EnvelopeHead
	CipherTxt []byte //encrypt from EnvelopeContent
	EnvelopeTail
	ErrorCode int
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

	eh := &sef.EnvelopeHead

	tmp, err = eh.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackLongBytes(sef.CipherTxt)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	et := &sef.EnvelopeTail

	tmp, err = et.Pack()
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

	eh := &sef.EnvelopeHead
	l, err = eh.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	sef.CipherTxt, l, err = UnPackLongBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += l

	et := &sef.EnvelopeTail
	l, err = et.UnPack(data[offset:])
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
