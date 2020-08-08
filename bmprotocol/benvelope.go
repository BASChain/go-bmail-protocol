package bmprotocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/realbmail/go-bmail-protocol/translayer"
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
	aesDecrypt func(cipherText, key []byte) (iv []byte, plainBytes []byte)
)

func RegAesEncrypt(enc func(plainBytes, iv, key []byte) []byte) {
	aesEncrypt = enc
}

func RegAesDecrypt(dec func(cipherText, key []byte) (iv []byte, plainBytes []byte)) {
	aesDecrypt = dec
}

type EnvelopeRoute struct {
	From         string
	RecpAddr     string                //recipient
	RecpAddrType int                   //0 to,1 cc,2 bc
	EId          translayer.EnveUniqID //envelope unique id
}

func (eh *EnvelopeRoute) CopyTo(to *EnvelopeRoute) *EnvelopeRoute {
	if to == nil {
		to = &EnvelopeRoute{}
	}

	to.From = eh.From
	to.RecpAddr = eh.RecpAddr
	to.RecpAddrType = eh.RecpAddrType
	to.EId = eh.EId

	return to

}

func (eh *EnvelopeRoute) String() string {
	s := fmt.Sprintf("%-20s", eh.From)
	s += fmt.Sprintf("%-20s", eh.RecpAddr)
	s += fmt.Sprintf("%-4d", eh.RecpAddrType)

	s += fmt.Sprintf("%-30s", base58.Encode(eh.EId[:]))

	return s
}

func (eh *EnvelopeRoute) Pack() ([]byte, error) {
	if eh.From == "" || eh.RecpAddr == "" {
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

	tmp, err = PackShortBytes(eh.EId[:])
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return r, nil
}

func (eh *EnvelopeRoute) UnPack(data []byte) (int, error) {
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

	var tmp []byte
	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	copy(eh.EId[:], tmp)

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
	Sn  []byte `json:"sn"`
	Sig []byte `json:"sig"`
}

func (es *EnvelopeSig) CopyTo(to *EnvelopeSig) *EnvelopeSig {
	if to == nil {
		to = &EnvelopeSig{}
	}

	to.Sn = append(to.Sn, es.Sn...)
	to.Sig = append(to.Sig, es.Sig...)

	return to
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

	ee.Sig, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil
}

//mode:
// CryptModePS uint16 = 3
// CryptModePP uint16 = 9
// CryptModePSP uint16 = 11
// CryptModePSSP uint16 = 15

type EnvelopeCryptDesc struct {
	Mode    int //ps, pp, psp, pssp
	Pubkeys [][]byte
}

func (ecd *EnvelopeCryptDesc) String() string {
	s := fmt.Sprintf("mode: %d", ecd.Mode)
	s += fmt.Sprintf("     pubkey count:%d\r\n", len(ecd.Pubkeys))

	for i := 0; i < len(ecd.Pubkeys); i++ {
		s += fmt.Sprintf("%s\r\n", base58.Encode(ecd.Pubkeys[i]))
	}

	return s

}

func (ecd *EnvelopeCryptDesc) CopyTo(ecd1 *EnvelopeCryptDesc) *EnvelopeCryptDesc {
	if ecd1 == nil {
		ecd1 = &EnvelopeCryptDesc{}
	}

	ecd1.Mode = ecd.Mode

	for i := 0; i < len(ecd.Pubkeys); i++ {
		buf := make([]byte, len(ecd.Pubkeys[i]))
		copy(buf, ecd.Pubkeys[i])

		ecd1.Pubkeys = append(ecd1.Pubkeys, buf)

	}
	return ecd1

}

func (ecd *EnvelopeCryptDesc) Pack() ([]byte, error) {
	var (
		tmp, r []byte
		err    error
	)
	tmp = translayer.UInt32ToBuf(uint32(ecd.Mode))
	r = append(r, tmp...)

	tmp, err = PackShortBytesArray(ecd.Pubkeys)
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return r, nil
}

func (ecd *EnvelopeCryptDesc) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)

	if len(data) < translayer.Uint32Size {
		return 0, errors.New("unpack mode error")
	}
	ecd.Mode = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	ecd.Pubkeys, of, err = UnPackShortBytesArray(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil

}

type Envelope struct {
	EnvelopeSig
	EnvelopeRoute
	EnvelopeCryptDesc
	EnvelopeContent
}

func (e *Envelope) String() string {
	s := e.EnvelopeSig.String()
	s += "\r\n"
	s += e.EnvelopeRoute.String()
	s += "\r\n"
	s += e.EnvelopeCryptDesc.String()
	s += "\r\n"
	s += e.EnvelopeContent.String()

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

	eh := &e.EnvelopeRoute

	b, err = eh.Pack()
	if err != nil {
		return nil, err
	}

	r = append(r, b...)

	ecd := &e.EnvelopeCryptDesc
	b, err = ecd.Pack()
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

	eh := &e.EnvelopeRoute

	of, err = eh.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	ecd := &e.EnvelopeCryptDesc

	of, err = ecd.UnPack(data[offset:])
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
	EnvelopeRoute
	EnvelopeCryptDesc
	CipherTxt []byte
}

func (ce *CryptEnvelope) String() string {

	s := ce.EnvelopeSig.String()
	s += "\r\n"
	s += ce.EnvelopeRoute.String()
	s += "\r\n"
	s += ce.EnvelopeCryptDesc.String()
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

	eh := &ce.EnvelopeRoute

	b, err = eh.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, b...)

	ecd := &ce.EnvelopeCryptDesc

	b, err = ecd.Pack()
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

	eh := &ce.EnvelopeRoute
	of, err = eh.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	ecd := &ce.EnvelopeCryptDesc
	of, err = ecd.UnPack(data[offset:])
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

func EncodeEnvelope(e *Envelope, key []byte) *CryptEnvelope {

	if e == nil {
		return nil
	}

	ce := &CryptEnvelope{}
	ceh := &ce.EnvelopeRoute
	(&e.EnvelopeRoute).CopyTo(ceh)

	csig := &ce.EnvelopeSig
	(&e.EnvelopeSig).CopyTo(csig)

	ecd := &ce.EnvelopeCryptDesc

	(&e.EnvelopeCryptDesc).CopyTo(ecd)

	data, err := e.ForCrypt()
	if err != nil {
		return nil
	}

	ce.CipherTxt = aesEncrypt(data, e.EnvelopeSig.Sn, key)

	if ce.CipherTxt == nil {
		return nil
	}

	return ce

}

func DeCodeEnvelope(ce *CryptEnvelope, key []byte) *Envelope {
	if ce == nil || len(ce.CipherTxt) == 0 {
		return nil
	}

	e := &Envelope{}

	eh := &e.EnvelopeRoute
	(&ce.EnvelopeRoute).CopyTo(eh)
	es := &e.EnvelopeSig
	(&ce.EnvelopeSig).CopyTo(es)

	ecd := &e.EnvelopeCryptDesc

	(&ce.EnvelopeCryptDesc).CopyTo(ecd)

	iv, plaintxt := aesDecrypt(ce.CipherTxt, key)
	if len(plaintxt) == 0 || bytes.Compare(iv, ce.EnvelopeSig.Sn) == 0 {
		return nil
	}

	ec := &e.EnvelopeContent

	_, err := ec.UnPack(plaintxt)
	if err != nil {
		return nil
	}

	return e

}

type ConfirmEnvelope struct {
	Sn         []byte
	NewSn      []byte
	EId        translayer.EnveUniqID
	CxtHashSig []byte
	ErrId      int
}

func (ce *ConfirmEnvelope) String() string {
	s := fmt.Sprintf("sn:%-30s", base58.Encode(ce.Sn))
	s += fmt.Sprintf("newsn:%-30s\r\n", base58.Encode(ce.NewSn))
	s += fmt.Sprintf("eid:%-30s", base58.Encode(ce.EId[:]))
	s += fmt.Sprintf("cxthashsig:%-30s", base58.Encode(ce.CxtHashSig))
	s += fmt.Sprintf("errid %d", ce.ErrId)

	return s
}

func (ce *ConfirmEnvelope) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortBytes(ce.Sn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(ce.NewSn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(ce.EId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(ce.CxtHashSig)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(ce.ErrId))
	r = append(r, tmp...)

	return r, nil

}

func (ce *ConfirmEnvelope) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	ce.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	ce.NewSn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	copy(ce.EId[:], tmp)

	ce.CxtHashSig, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack errid error")
	}
	ce.ErrId = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil

}
