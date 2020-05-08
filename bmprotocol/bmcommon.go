package bmprotocol

import (
	"github.com/BASChain/go-bmail-protocol/translayer"
	"encoding/binary"
	"strconv"
	"errors"
)

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
		return nil, 0, errors.New("Unpack Short bytes Failed")
	}

	offset := 0

	l := binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	if len(data) < offset+int(l) {
		return nil, 0, errors.New("Unpack Short bytes Failed")
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

func NewHeadBuf() []byte {
	return make([]byte,translayer.BMHeadSize())
}

func AddPackHead(bmtl *translayer.BMTransLayer,appendData []byte) ([]byte,error)  {

	datalen := len(appendData) - translayer.BMHeadSize()
	if bmtl == nil || datalen < 0{
		return nil,errors.New("head is nil")
	}

	bmtl.SetDataLen(uint32(datalen))

	r,err:=bmtl.Pack()
	if err!=nil{
		return nil,err
	}

	copy(appendData[:translayer.BMHeadSize()],r)

	return appendData,nil
}