package translayer

import (
	"github.com/pkg/errors"
	"encoding/binary"
)

const BMAILVER1 uint16 = 1
const ED25519 uint16 = 1

type BMTransLayer struct {
	ver uint16
	cryptType uint16
	typ uint16
	data []byte
}

func (bmtl *BMTransLayer)GetData() []byte  {
	return bmtl.data
}

func NewBMTL(typ uint16,data []byte) *BMTransLayer {
	bmtl:=&BMTransLayer{}

	bmtl.ver = BMAILVER1
	bmtl.cryptType = ED25519

	bmtl.typ = typ
	bmtl.data = data

	return bmtl
}

func UInt16ToBuf(ui16 uint16) []byte {
	bufl:=make([]byte,2)

	binary.BigEndian.PutUint16(bufl,ui16)

	return bufl
}

func UInt32ToBuf(ui32 uint32) []byte {
	bufl:=make([]byte,4)

	binary.BigEndian.PutUint32(bufl,ui32)

	return bufl
}

func UInt64ToBuf(ui64 uint64 ) []byte {
	bufl:=make([]byte,8)

	binary.BigEndian.PutUint64(bufl,ui64)

	return bufl
}



func (bmtl *BMTransLayer)Pack() ([]byte,error) {

	if bmtl.typ <= MIN_TYP || bmtl.typ>MAX_TYP{
		return nil,errors.New("BMail Action Type Error")
	}
	var r []byte

	bufl:=UInt16ToBuf(uint16(bmtl.ver))
	r = append(r,bufl...)

	bufl=UInt16ToBuf(uint16(bmtl.cryptType))
	r = append(r,bufl...)

	bufl=UInt16ToBuf(uint16(bmtl.typ))
	r = append(r,bufl...)

	if len(bmtl.data)>0{
		l:=uint32(len(bmtl.data))


		bufl := UInt32ToBuf(uint32(l))

		r = append(r,bufl...)

		r = append(r,bmtl.data...)
	}

	return r,nil
}

func (bmtl *BMTransLayer)UnPack(data []byte) (int,error) {

	if len(data)<10{
		return 0,errors.New("Not a BMail Action Data")
	}

	offset := 0
	bmtl.ver = binary.BigEndian.Uint16(data[offset:])
	offset += 2

	bmtl.cryptType = binary.BigEndian.Uint16(data[offset:])
	offset += 2

	bmtl.typ = binary.BigEndian.Uint16(data[offset:])
	offset += 2

	if bmtl.typ <= MIN_TYP || bmtl.typ>=MAX_TYP{
		return 0,errors.New("BMail Action Type Error")
	}

	l:=binary.BigEndian.Uint32(data[offset:])
	offset += 4

	bmtl.data = data[offset:]

	if l != uint32(len(bmtl.data)){
		return 0,errors.New("Data Length Error")
	}

	return offset + len(bmtl.data),nil
}

