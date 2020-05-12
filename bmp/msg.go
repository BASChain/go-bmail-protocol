package bmp

import (
	"encoding/binary"
	"fmt"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/translayer"
)

type Header struct {
	Ver     uint16
	MsgTyp  uint16
	DataLen uint32
}

func (h *Header) Pack() []byte {
	var r []byte

	bufl := translayer.UInt16ToBuf(uint16(h.Ver))
	r = append(r, bufl...)

	bufl = translayer.UInt16ToBuf(uint16(h.MsgTyp))
	r = append(r, bufl...)

	bufl = translayer.UInt32ToBuf(h.DataLen)

	r = append(r, bufl...)

	return r
}

func (h *Header) GetLen() int {
	if translayer.BMAILVER1 == 1 {
		return 8
	}
	return -1
}

func (h *Header) Unpack(data []byte) error {
	if len(data) < h.GetLen() {
		return fmt.Errorf("not a BMail Action Data")
	}

	offset := 0
	h.Ver = binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	h.MsgTyp = binary.BigEndian.Uint16(data[offset:])
	offset += translayer.Uint16Size

	if h.MsgTyp <= translayer.MIN_TYP || h.MsgTyp >= translayer.MAX_TYP {
		return fmt.Errorf("BMail Action Type Error")
	}

	l := binary.BigEndian.Uint32(data[offset:])
	offset += translayer.Uint32Size

	h.DataLen = l
	return nil
}

type PlainHELO struct {
	SrvAddr []bmail.Address
}
