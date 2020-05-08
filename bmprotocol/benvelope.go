package bmprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
)

const (
	PeerUnreachable int = iota + 1
	AddressUnavailable
)

const (
	ErrMsg_PeerUnreachable    string = "Peer is unreachable"
	ErrMsg_AddressUnavailable string = "Recipient is not available"
)

var (
	aesEncrypt func(plainBytes, iv, key []byte) []byte
	aesDecrypt func(cipherText, key []byte) (iv []byte,plainBytes []byte)
)



type EnvelopeHead struct {
	From         string
	RecpAddr     string //recipient
	RecpAddrType int    //0 to,1 cc,2 bc
	LPubKey      []byte //local public key
	EId          translayer.EnveUniqID   //envelope unique id
}

func (eh *EnvelopeHead) String() string {
	s := fmt.Sprintf("%-20s", eh.From)
	s += fmt.Sprintf("%-20s", eh.RecpAddr)
	s += fmt.Sprintf("%-4d", eh.RecpAddrType)
	s += fmt.Sprintf("%-50s", base58.Encode(eh.LPubKey))
	s += fmt.Sprintf("%-30s",base58.Encode(eh.EId[:]))

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
	tmp, err = PackShortBytes(eh.EId[:])
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

	var tmp []byte
	tmp,of,err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	copy(eh.EId[:],tmp)

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

type EnvelopeSig struct {
	Sn  []byte //sn from Bhello
	Sig []byte //signature
}

func (et *EnvelopeSig) String() string {
	s := fmt.Sprintf("%-50s", base58.Encode(et.Sn))
	s += fmt.Sprintf("%s", base58.Encode(et.Sig))

	return s
}

func (ee *EnvelopeSig) Pack() ([]byte, error) {
	if len(ee.Sn) == 0 || len(ee.Sig) == 0 {
		return nil, errors.New("Not a Correct Envelope Tail")
	}

	var r []byte

	riv, err := PackShortBytes(ee.Sn)
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

func (ee *EnvelopeSig) UnPack(data []byte) (int, error) {
	offset := 0
	var of int
	var err error
	ee.Sn, of, err = UnPackShortBytes(data[offset:])

	if err != nil {
		return 0, err
	}

	offset += of

	ee.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil
}

type Envelope struct {
	EnvelopeSig
	EnvelopeHead
	EnvelopeContent

}

func (e *Envelope)String() string  {
	s:=e.EnvelopeSig.String()
	s+="\r\n"
	s+=e.EnvelopeHead.String()
	s+="\r\n"

	s+=e.EnvelopeContent.String()


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

	var r []byte

	es := &e.EnvelopeSig

	b, err := es.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, b...)

	eh := &e.EnvelopeHead



	b, err = eh.Pack()
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



	return r, nil
}

func (e *Envelope) UnPack(data []byte) (int, error) {


	var (
		of, offset int
		err        error
	)

	es := &e.EnvelopeSig

	of, err = es.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	eh := &e.EnvelopeHead

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



	return offset, nil

}

type CryptEnvelope struct {
	EnvelopeSig
	EnvelopeHead
	CipherTxt []byte

}

func (ce *CryptEnvelope) String() string {

	s := ce.EnvelopeSig.String()
	s += "\r\n"
	s += ce.EnvelopeHead.String()
	s += "\r\n"
	s += base58.Encode(ce.CipherTxt)


	return s
}

func (ce *CryptEnvelope) Pack() ([]byte, error) {
	var r []byte

	es := &ce.EnvelopeSig
	b, err := es.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, b...)

	eh := &ce.EnvelopeHead


	b, err = eh.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, b...)

	b, err = PackLongBytes(ce.CipherTxt)
	if err != nil {
		return nil, err
	}
	r = append(r, b...)



	return r, nil
}

func (ce *CryptEnvelope) UnPack(data []byte) (int, error) {

	var (
		of, offset int
		err        error
	)

	es := &ce.EnvelopeSig
	of, err = es.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	eh := &ce.EnvelopeHead
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



	return offset, nil
}
