package bmprotocol

import (
	"encoding/binary"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/pkg/errors"
	"fmt"
)

type BPOPStat struct {
	translayer.BMTransLayer
}

func NewBPOPStat() *BPOPStat {
	bmtl := translayer.NewBMTL(translayer.STAT, nil)

	bps := &BPOPStat{}

	bps.BMTransLayer = *bmtl

	return bps
}

func (bps *BPOPStat) Pack() ([]byte, error) {
	return bps.BMTransLayer.Pack()
}

func (bps *BPOPStat) UnPack(data []byte) (int, error) {
	return 0, nil
}

type BPOPStatResp struct {
	translayer.BMTransLayer
	Total            int
	Received         int
	TotalStoredBytes int64
	TotalSpaceBytes  int64
}

func (bs *BPOPStatResp)String() string  {
	s:=bs.BMTransLayer.HeadString()
	s+=fmt.Sprintf("%-10d",bs.Total)
	s+=fmt.Sprintf("%-10d",bs.Received)
	s+=fmt.Sprintf("%-20d",bs.TotalStoredBytes)
	s+=fmt.Sprintf("%-20d",bs.TotalSpaceBytes)

	return s
}

func NewBPOPStatResp() *BPOPStatResp {
	bmtl := translayer.NewBMTL(translayer.STAT_RESP, nil)

	bpsr := &BPOPStatResp{}

	bpsr.BMTransLayer = *bmtl

	return bpsr
}

func (br *BPOPStatResp) Pack() ([]byte, error) {
	var (
		r, tmp []byte
	)
	tmp = translayer.UInt32ToBuf(uint32(br.Total))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(br.Received))
	r = append(r, tmp...)

	tmp = translayer.UInt64ToBuf(uint64((br.TotalStoredBytes)))
	r = append(r, tmp...)

	tmp = translayer.UInt64ToBuf(uint64(br.TotalSpaceBytes))
	r = append(r, tmp...)

	br.BMTransLayer.SetData(r)

	return br.BMTransLayer.Pack()
}

func (br *BPOPStatResp) UnPack(data []byte) (int, error) {
	var (
		offset int
	)

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack total error")
	}

	br.Total = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack Received error")
	}

	br.Received = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint64Size {
		return 0, errors.New("unpack TotalStoredBytes error")
	}

	br.TotalStoredBytes = int64(binary.BigEndian.Uint64(data[offset:]))
	offset += translayer.Uint64Size

	if len(data) < offset+translayer.Uint64Size {
		return 0, errors.New("unpack TotalSpaceBytes error")
	}
	br.TotalSpaceBytes = int64(binary.BigEndian.Uint64(data[offset:]))

	offset += translayer.Uint64Size

	return offset, nil

}

type BPOPList struct {
	translayer.BMTransLayer
	BeginID   int
	ListCount int
}

func NewBPOPList() *BPOPList {
	bpl := &BPOPList{}

	bmtl := translayer.NewBMTL(translayer.LIST, nil)

	bpl.BMTransLayer = *bmtl

	return bpl
}

func (bl *BPOPList)String() string  {
	s:=bl.BMTransLayer.HeadString()
	s += fmt.Sprintf("%-10d",bl.BeginID)
	s += fmt.Sprintf("%-10d",bl.ListCount)

	return s
}

func (bl *BPOPList) Pack() ([]byte, error) {
	var (
		r, tmp []byte
	)

	tmp = translayer.UInt32ToBuf(uint32(bl.BeginID))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(bl.ListCount))
	r = append(r, tmp...)

	bl.BMTransLayer.SetData(r)

	return bl.BMTransLayer.Pack()
}

func (bl *BPOPList) UnPack(data []byte) (int, error) {
	var (
		offset int
	)

	if len(data) < translayer.Uint32Size {
		return 0, errors.New("unpack BeginID error")
	}

	bl.BeginID = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack ListCount error")
	}
	bl.ListCount = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil
}

type ListNode struct {
	ID          int
	SizeOfBytes int
}

func (ln *ListNode) Pack() ([]byte, error) {
	var (
		r, tmp []byte
	)

	tmp = translayer.UInt32ToBuf(uint32(ln.ID))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(ln.SizeOfBytes))

	r = append(r, tmp...)

	return r, nil
}

func (ln *ListNode) UnPack(data []byte) (int, error) {
	var (
		offset int
	)
	if len(data) < translayer.Uint32Size {
		return 0, errors.New("unpack ID error")
	}

	ln.ID = int(binary.BigEndian.Uint32(data))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack ID error")
	}

	ln.SizeOfBytes = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil
}

type BPOPListResp struct {
	translayer.BMTransLayer
	BeginID   int
	ListCount int
	Nodes     []*ListNode
}

func NewBPOPListResp() *BPOPListResp {
	bmtl := translayer.NewBMTL(translayer.LIST_RESP, nil)
	bl := &BPOPListResp{}
	bl.BMTransLayer = *bmtl

	return bl
}

func (bl *BPOPListResp) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp = translayer.UInt32ToBuf(uint32(bl.BeginID))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(bl.ListCount))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(len(bl.Nodes)))
	r = append(r, tmp...)

	for i := 0; i < len(bl.Nodes); i++ {
		n := bl.Nodes[i]

		tmp, err = n.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)

	}

	bl.BMTransLayer.SetData(r)

	return bl.BMTransLayer.Pack()

}

func (bl *BPOPListResp) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)
	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack BeginID error")
	}

	bl.BeginID = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack ListCount error")
	}
	bl.ListCount = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack Node count error")
	}
	l := int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	for i := 0; i < l; i++ {
		n := &ListNode{}
		of, err = n.UnPack(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += of
	}

	return offset, nil
}

type BPOPRetr struct {
	translayer.BMTransLayer
	BeginID   int
	RetrCount int
}

func NewBPOPRetr() *BPOPRetr {
	bmtl := translayer.NewBMTL(translayer.RETR, nil)

	br := &BPOPRetr{}

	br.BMTransLayer = *bmtl

	return br
}

func (br *BPOPRetr) Pack() ([]byte, error) {
	var (
		r, tmp []byte
	)
	tmp = translayer.UInt32ToBuf(uint32(br.BeginID))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(br.RetrCount))
	r = append(r, tmp...)

	br.BMTransLayer.SetData(r)

	return br.BMTransLayer.Pack()
}

func (br *BPOPRetr) UnPack(data []byte) (int, error) {
	var (
		offset int
	)
	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack BeginID error")
	}

	br.BeginID = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack ListCount error")
	}
	br.RetrCount = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil
}

type BPOPRetrResp struct {
	translayer.BMTransLayer
	Mails      []CryptEnvelope
	BeginID    int
	RetrCount  int
	TotalCount int
}

func NewBPOPRetrResp() *BPOPRetrResp {
	bmtl := translayer.NewBMTL(translayer.RETR_RESP, nil)
	br := &BPOPRetrResp{}

	br.BMTransLayer = *bmtl

	return br
}

func (br *BPOPRetrResp) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)
	tmp = translayer.UInt32ToBuf(uint32(len(br.Mails)))

	r = append(r, tmp...)

	for i := 0; i < len(br.Mails); i++ {
		ce := &br.Mails[i]

		tmp, err = ce.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)
	}
	tmp = translayer.UInt32ToBuf(uint32(br.BeginID))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(br.RetrCount))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(br.TotalCount))
	r = append(r, tmp...)

	br.BMTransLayer.SetData(r)

	return br.BMTransLayer.Pack()

}

func (br *BPOPRetrResp) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)

	if offset+translayer.Uint32Size > len(data) {
		return 0, errors.New("unpack mails error")
	}
	mailcnt := int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	for i := 0; i < mailcnt; i++ {
		ce := &CryptEnvelope{}
		of, err = ce.UnPack(data[offset:])
		if err != nil {
			return 0, nil
		}
		offset += of
		br.Mails = append(br.Mails, *ce)
	}

	if offset+translayer.Uint32Size > len(data) {
		return 0, errors.New("unpack begin id error")
	}
	br.BeginID = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if offset+translayer.Uint32Size > len(data) {
		return 0, errors.New("unpack retr count error")
	}
	br.RetrCount = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if offset+translayer.Uint32Size > len(data) {
		return 0, errors.New("unpack total count error")
	}
	br.TotalCount = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil

}
