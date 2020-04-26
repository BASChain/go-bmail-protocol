package bmprotocol

import (
	"encoding/binary"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/pkg/errors"
)

type BMHello struct {
	translayer.BMTransLayer
}

func NewBMHello() *BMHello {
	bmtl := translayer.NewBMTL(translayer.HELLO, nil)
	bmh := &BMHello{}
	bmh.BMTransLayer = *bmtl

	return bmh
}

func (bmh *BMHello) Pack() ([]byte, error) {
	return bmh.BMTransLayer.Pack()
}

func (bmh *BMHello) UnPack(data []byte) (int, error) {
	//nothing todo
	return 0, nil
}

func (bmh *BMHello) String() string {
	return bmh.BMTransLayer.String()
}

type BMHelloACK struct {
	translayer.BMTransLayer
	sn []byte
}

func NewBMHelloACK(sn []byte) *BMHelloACK {
	bmact := translayer.NewBMTL(translayer.HELLO_ACK, nil)

	bmhack := &BMHelloACK{}

	bmhack.BMTransLayer = *bmact

	bmhack.sn = sn

	return bmhack
}

func (bmha *BMHelloACK) Pack() ([]byte, error) {
	var barr []byte

	bufl := translayer.UInt16ToBuf(uint16(len(bmha.sn)))

	barr = append(barr, bufl...)
	barr = append(barr, bmha.sn...)

	bmha.BMTransLayer.SetData(barr)

	return bmha.BMTransLayer.Pack()
}

func (bmha *BMHelloACK) String() string {
	s := bmha.BMTransLayer.HeadString()

	s += fmt.Sprintf("sn: %s", base58.Encode(bmha.sn))

	return s
}

func (bmha *BMHelloACK) UnPack(data []byte) (int, error) {

	if len(data) < translayer.Uint16Size {
		return 0, errors.New("Not a HELLO ACK data")
	}

	l := binary.BigEndian.Uint16(data)

	bmha.sn = data[translayer.Uint16Size:]

	if l != uint16(len(bmha.sn)) {
		return 0, errors.New("Serial Nunber Error")
	}

	return translayer.Uint16Size + len(bmha.sn), nil
}

/*
*

pack:
1.NewSendSignature
2.ForSigBuf
3.Calculate signature
4.SetSig
5.Pack
6.Send

unpack:
bmtl:=&BMTransLayer{}
bmtl.UnPack(buf)
if bmtl.typ == SEND_SIGNATURE{
	ss:=&SendSignature{}
	ss.UnPack(buf[bmtl.GetData()])
}

*/

type SendSignature struct {
	translayer.BMTransLayer
	sn            []byte
	localMailAddr string
	currentTime   int64 //Millisecond
	sig           []byte
}

func NewSendSignature(sn []byte, localMailAddr string) *SendSignature {
	ss := &SendSignature{}
	ss.sn = sn
	ss.localMailAddr = localMailAddr
	ss.currentTime = tools.GetNowMsTime()

	bmtl := translayer.NewBMTL(translayer.SEND_SIGNATURE, nil)

	ss.BMTransLayer = *bmtl

	return ss
}

func (ss *SendSignature) String() string {
	s := ss.BMTransLayer.HeadString()

	s += fmt.Sprintf("sn:%s\r\n", base58.Encode(ss.sn))

	s += fmt.Sprintf("mail:%s\r\n", ss.localMailAddr)

	s += fmt.Sprintf("time:%d\r\n", ss.currentTime)

	s += fmt.Sprintf("sig:%s\r\n", base58.Encode(ss.sig))

	return s
}

func (ss *SendSignature) GetSig() []byte {
	return ss.sig
}

func (ss *SendSignature) ForSigBuf() []byte {
	var r []byte

	bufl := translayer.UInt16ToBuf(uint16(len(ss.sn)))

	r = append(r, bufl...)

	r = append(r, ss.sn...)

	bufl = translayer.UInt16ToBuf(uint16(len(ss.localMailAddr)))

	r = append(r, bufl...)

	r = append(r, []byte(ss.localMailAddr)...)

	bufl = translayer.UInt64ToBuf(uint64(ss.currentTime))

	r = append(r, bufl...)

	return r
}

func (ss *SendSignature) SetSig(sig []byte) {
	ss.sig = sig
}

func (ss *SendSignature) Pack() ([]byte, error) {
	r := ss.ForSigBuf()

	bufl := translayer.UInt16ToBuf(uint16(len(ss.sig)))

	r = append(r, bufl...)

	r = append(r, ss.sig...)

	ss.BMTransLayer.SetData(r)

	return ss.BMTransLayer.Pack()
}

func (ss *SendSignature) UnPack(buf []byte) (int, error) {
	offset := 0

	if len(buf[offset:]) < translayer.Uint16Size {
		return 0, errors.New("Not a SendSignature data")
	}

	lsn := binary.BigEndian.Uint16(buf[offset:])
	offset += translayer.Uint16Size
	if len(buf) < offset+int(lsn) {
		return 0, errors.New("Not a SendSignature data")
	}
	ss.sn = buf[offset : offset+int(lsn)]
	offset += int(lsn)

	if len(buf[offset:]) < translayer.Uint16Size {
		return 0, errors.New("Not a SendSignature data")
	}
	laddr := binary.BigEndian.Uint16(buf[offset:])
	offset += translayer.Uint16Size
	if len(buf) < offset+int(laddr) {
		return 0, errors.New("Not a SendSignature data")
	}
	ss.localMailAddr = string(buf[offset : offset+int(laddr)])
	offset += int(laddr)
	if len(buf[offset:]) < translayer.Uint64Size {
		return 0, errors.New("Not a SendSignature data")
	}
	ss.currentTime = int64(binary.BigEndian.Uint64(buf[offset:]))
	offset += translayer.Uint64Size

	if len(buf[offset:]) < translayer.Uint16Size {
		return 0, errors.New("Not a SendSignature data")
	}

	lsig := binary.BigEndian.Uint16(buf[offset:])
	offset += translayer.Uint16Size
	if len(buf) < offset+int(lsig) {
		return 0, errors.New("Not a SendSignature data")
	}

	ss.sig = buf[offset:]

	return offset + int(lsig), nil

}

type ValidateSignature struct {
	translayer.BMTransLayer
	sn []byte
}

func NewValidSign(sn []byte) *ValidateSignature {
	vs := &ValidateSignature{}
	bmact := translayer.NewBMTL(translayer.VALIDATE_SIGNATURE, nil)

	vs.BMTransLayer = *bmact

	vs.sn = sn

	return vs
}

func (vs *ValidateSignature) String() string {
	s := vs.BMTransLayer.HeadString()

	s += fmt.Sprintf("sn: %s", base58.Encode(vs.sn))

	return s
}

func (vs *ValidateSignature) Pack() ([]byte, error) {
	var barr []byte

	bufl := translayer.UInt16ToBuf(uint16(len(vs.sn)))

	barr = append(barr, bufl...)
	barr = append(barr, vs.sn...)

	vs.BMTransLayer.SetData(barr)

	return vs.BMTransLayer.Pack()
}

func (vs *ValidateSignature) UnPack(data []byte) (int, error) {

	if len(data) < translayer.Uint16Size {
		return 0, errors.New("Not a Validate Signature data")
	}

	l := binary.BigEndian.Uint16(data)

	vs.sn = data[translayer.Uint16Size:]

	if l != uint16(len(vs.sn)) {
		return 0, errors.New("Serial Nunber Error")
	}

	return translayer.Uint16Size + len(vs.sn), nil
}
