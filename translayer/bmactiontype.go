package translayer

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

const BMAILVER1 uint16 = 1
const ED25519 uint16 = 1

type BMTransLayer struct {
	ver       uint16
	cryptType uint16
	typ       uint16
	dataLen   uint32
}

func BMHeadSize() int {
	rv := reflect.ValueOf(BMTransLayer{})

	cnt := rv.NumField()

	size := 0

	for i := 0; i < cnt; i++ {
		f := rv.Field(i)
		switch f.Kind() {
		case reflect.Uint16:
			size += Uint16Size
		case reflect.Uint8:
			size += Uin8Size
		case reflect.Uint32:
			size += Uint32Size
		case reflect.Uint64:
			size += Uint64Size
		case reflect.Slice:
			size += 0
		}
	}

	return size
}

func (bmtl *BMTransLayer) String() string {
	s := fmt.Sprintf("Version: %-4d", bmtl.ver)
	s += fmt.Sprintf("CryptType: %-4d", bmtl.cryptType)
	s += fmt.Sprintf("MsgType: %-4d", bmtl.typ)
	s += fmt.Sprintf("DataLength:%-8d\r\n", bmtl.dataLen)

	return s
}

func (bmtl *BMTransLayer) SetDataLen(l uint32) {
	bmtl.dataLen = l
}

func (bmtl *BMTransLayer) getDataLen() uint32 {
	return bmtl.dataLen
}

func NewBMTL(typ uint16) *BMTransLayer {
	bmtl := &BMTransLayer{}

	bmtl.ver = BMAILVER1
	bmtl.cryptType = ED25519

	bmtl.typ = typ

	return bmtl
}

func UInt16ToBuf(ui16 uint16) []byte {
	bufl := make([]byte, Uint16Size)

	binary.BigEndian.PutUint16(bufl, ui16)

	return bufl
}

func UInt32ToBuf(ui32 uint32) []byte {
	bufl := make([]byte, Uint32Size)

	binary.BigEndian.PutUint32(bufl, ui32)

	return bufl
}

func UInt64ToBuf(ui64 uint64) []byte {
	bufl := make([]byte, Uint64Size)

	binary.BigEndian.PutUint64(bufl, ui64)

	return bufl
}

func (bmtl *BMTransLayer) Pack() ([]byte, error) {

	if bmtl.typ <= MIN_TYP || bmtl.typ > MAX_TYP {
		return nil, errors.New("BMail Action Type Error")
	}

	var r []byte

	bufl := UInt16ToBuf(uint16(bmtl.ver))
	r = append(r, bufl...)

	bufl = UInt16ToBuf(uint16(bmtl.cryptType))
	r = append(r, bufl...)

	bufl = UInt16ToBuf(uint16(bmtl.typ))
	r = append(r, bufl...)

	bufl = UInt32ToBuf(bmtl.dataLen)

	r = append(r, bufl...)

	return r, nil
}

func (bmtl *BMTransLayer) UnPack(data []byte) (int, error) {
	if len(data) < BMHeadSize() {
		return 0, errors.New("Not a BMail Action Data")
	}

	offset := 0
	bmtl.ver = binary.BigEndian.Uint16(data[offset:])
	offset += Uint16Size

	bmtl.cryptType = binary.BigEndian.Uint16(data[offset:])
	offset += Uint16Size

	bmtl.typ = binary.BigEndian.Uint16(data[offset:])
	offset += Uint16Size

	if bmtl.typ <= MIN_TYP || bmtl.typ >= MAX_TYP {
		return 0, errors.New("BMail Action Type Error")
	}

	l := binary.BigEndian.Uint32(data[offset:])
	offset += Uint32Size

	bmtl.dataLen = l

	return offset, nil
}
