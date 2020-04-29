package bmprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
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
	From         string
	RecpAddr     string //recipient
	RecpAddrType int    //0 to,1 cc,2 bc
	LPubKey      []byte //local public key
}

func (eh *EnvelopeHead) String() string {
	s := fmt.Sprintf("%-20s", eh.From)
	s += fmt.Sprintf("%-20s", eh.RecpAddr)
	s += fmt.Sprintf("%-4d", eh.RecpAddrType)
	s += fmt.Sprintf("%-50s", base58.Encode(eh.LPubKey))

	return s
}

func (eh *EnvelopeHead) Pack() ([]byte, error) {
	if eh.From == "" || eh.RecpAddr == "" || len(eh.LPubKey) == 0 {
		return nil, errors.New("")
	}

	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortString(eh.From)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortString(eh.RecpAddr)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(eh.RecpAddrType))
	r = append(r, tmp...)

	tmp, err = PackShortBytes(eh.LPubKey)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return r, nil
}

func (eh *EnvelopeHead) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)
	eh.From, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, nil
	}

	offset += of

	eh.RecpAddr, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack recp addr type error")
	}
	eh.RecpAddrType = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	eh.LPubKey, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	return offset, nil
}

type EnvelopeContent struct {
	To      []string
	CC      []string
	BC      []string
	Subject string
	Data    string
	Files   []Attachment
}

func (ec *EnvelopeContent) String() string {
	s := "to: "
	for i := 0; i < len(ec.To); i++ {
		s += fmt.Sprintf("%-20s", ec.To[i])
	}
	if len(ec.To) > 0 {
		s += "\r\n"
	}

	s += "cc: "
	for i := 0; i < len(ec.CC); i++ {
		s += fmt.Sprintf("%-20s", ec.CC[i])
	}

	if len(ec.CC) > 0 {
		s += "\r\n"
	}

	s += "bc: "
	for i := 0; i < len(ec.BC); i++ {
		s += fmt.Sprintf("%-20s", ec.BC[i])
	}

	if len(ec.BC) > 0 {
		s += "\r\n"
	}

	s += "subject: "
	s += fmt.Sprintf("%-50s\r\n", ec.Subject)

	s += "data: "

	s += fmt.Sprintf("%s\r\n", ec.Data)

	s += fmt.Sprintf("Attachment count: %-4d\r\n", len(ec.Files))
	for i := 0; i < len(ec.Files); i++ {
		a := &(ec.Files[i])
		s += a.String()
	}

	return s
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

	fcnt := len(ec.Files)
	tmp := translayer.UInt16ToBuf(uint16(fcnt))
	r = append(r, tmp...)

	for i := 0; i < len(ec.Files); i++ {
		a := &ec.Files[i]
		tmp, err = a.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)
	}

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
	sbj, of, err = UnPackLongString(data[offset:])
	if err != nil {
		return 0, err
	}

	ec.Subject = sbj
	offset += of

	var ed string
	ed, of, err = UnPackLongString(data[offset:])

	ec.Data = ed
	offset += of

	if len(data) < offset+translayer.Uint16Size {
		return 0, errors.New("unpack Files error")
	}

	cnt := binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	for i := 0; i < int(cnt); i++ {
		a := &Attachment{}
		of, err = a.UnPack(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += of

		ec.Files = append(ec.Files, *a)
	}

	return offset, nil
}

type EnvelopeTail struct {
	IV  []byte //sn from Bhello
	Sig []byte //signature
}

func (et *EnvelopeTail) String() string {
	s := fmt.Sprintf("%-50s", base58.Encode(et.IV))
	s += fmt.Sprintf("%s", base58.Encode(et.Sig))

	return s
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

type Envelope struct {
	EnvelopeHead
	EnvelopeContent

	EnvelopeTail
}

func (e *Envelope)String() string  {
	s:=e.EnvelopeHead.String()
	s+="\r\n"
	s+=e.EnvelopeContent.String()
	s+="\r\n"
	s+=e.EnvelopeTail.String()

	return s
}


func (e *Envelope) ForCrypt() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	ec := &e.EnvelopeContent
	tmp, err = ec.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return r, nil
}

func (e *Envelope) Pack() ([]byte, error) {
	eh := &e.EnvelopeHead

	var r []byte

	b, err := eh.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, b...)

	ec := &e.EnvelopeContent
	b, err = ec.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, b...)

	et := &e.EnvelopeTail

	b, err = et.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, b...)

	return r, nil
}

func (e *Envelope) UnPack(data []byte) (int, error) {
	eh := &e.EnvelopeHead

	var (
		of, offset int
		err        error
	)

	of, err = eh.UnPack(data)
	if err != nil {
		return 0, err
	}
	offset += of

	ec := &e.EnvelopeContent
	of, err = ec.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	et := &e.EnvelopeTail

	of, err = et.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil

}

type CryptEnvelope struct {
	EnvelopeHead
	CipherTxt []byte
	EnvelopeTail
}

func (ce *CryptEnvelope) String() string {
	s := ce.EnvelopeHead.String()
	s += "\r\n"
	s += base58.Encode(ce.CipherTxt)
	s += "\r\n"
	s += ce.EnvelopeTail.String()

	return s
}

func (ce *CryptEnvelope) Pack() ([]byte, error) {
	eh := &ce.EnvelopeHead

	var r []byte
	b, err := eh.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, b...)

	b, err = PackLongBytes(ce.CipherTxt)
	if err != nil {
		return nil, err
	}
	r = append(r, b...)

	et := &ce.EnvelopeTail
	b, err = et.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, b...)

	return r, nil
}

func (ce *CryptEnvelope) UnPack(data []byte) (int, error) {
	eh := &ce.EnvelopeHead
	var (
		of, offset int
		err        error
	)

	of, err = eh.UnPack(data)
	if err != nil {
		return 0, err
	}
	offset += of

	ce.CipherTxt, of, err = UnPackLongBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	et := &ce.EnvelopeTail
	of, err = et.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	return offset, nil
}
