package bmprotocol

import (
	"encoding/binary"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
)

type BPOPStat struct {
	translayer.BMTransLayer
}

func NewBPOPStat() *BPOPStat {
	bmtl := translayer.NewBMTL(translayer.STAT)

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

func (bs *BPOPStatResp) String() string {
	s := bs.BMTransLayer.String()
	s += fmt.Sprintf("%-10d", bs.Total)
	s += fmt.Sprintf("%-10d", bs.Received)
	s += fmt.Sprintf("%-20d", bs.TotalStoredBytes)
	s += fmt.Sprintf("%-20d", bs.TotalSpaceBytes)

	return s
}

func NewBPOPStatResp() *BPOPStatResp {
	bmtl := translayer.NewBMTL(translayer.STAT_RESP)

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

	//br.BMTransLayer.SetData(r)

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

	bmtl := translayer.NewBMTL(translayer.LIST)

	bpl.BMTransLayer = *bmtl

	return bpl
}

func (bl *BPOPList) String() string {
	s := bl.BMTransLayer.String()
	s += fmt.Sprintf("%-10d", bl.BeginID)
	s += fmt.Sprintf("%-10d", bl.ListCount)

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

	//bl.BMTransLayer.SetData(r)

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

func (ln *ListNode) String() string {
	s := fmt.Sprintf("ID: %-10d", ln.ID)
	s += fmt.Sprintf("SizeOfBytes: %-20d", ln.SizeOfBytes)

	return s
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
		return 0, errors.New("unpack sizeofbytes error")
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
	bmtl := translayer.NewBMTL(translayer.LIST_RESP)
	bl := &BPOPListResp{}
	bl.BMTransLayer = *bmtl

	return bl
}

func (bl *BPOPListResp) String() string {
	s := bl.BMTransLayer.String()
	s += fmt.Sprintf("BeginId %-10d", bl.BeginID)
	s += fmt.Sprintf("ListCount: %-10d", bl.ListCount)
	s += fmt.Sprintf("RealCount: %-10d\r\n", len(bl.Nodes))
	for i := 0; i < len(bl.Nodes); i++ {
		s += bl.Nodes[i].String()
		s += "\r\n"
	}

	return s
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

	//bl.BMTransLayer.SetData(r)

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
		bl.Nodes = append(bl.Nodes, n)
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
	bmtl := translayer.NewBMTL(translayer.RETR)

	br := &BPOPRetr{}

	br.BMTransLayer = *bmtl

	return br
}

func (br *BPOPRetr) String() string {
	s := br.BMTransLayer.String()
	s += fmt.Sprintf("BeginId: %-10d", br.BeginID)
	s += fmt.Sprintf("RetrCount: %-10d", br.RetrCount)

	return s
}

func (br *BPOPRetr) Pack() ([]byte, error) {
	var (
		r, tmp []byte
	)
	tmp = translayer.UInt32ToBuf(uint32(br.BeginID))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(br.RetrCount))
	r = append(r, tmp...)

	//br.BMTransLayer.SetData(r)

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
	bmtl := translayer.NewBMTL(translayer.RETR_RESP)
	br := &BPOPRetrResp{}

	br.BMTransLayer = *bmtl

	return br
}

func (br *BPOPRetrResp) String() string {
	s := br.BMTransLayer.String()
	s += fmt.Sprintf("mails count:%-10d", len(br.Mails))
	for i := 0; i < len(br.Mails); i++ {
		ce := &br.Mails[i]

		s += ce.String()
		s += "\r\n"
	}
	s += fmt.Sprintf("beginid: %-10d", br.BeginID)
	s += fmt.Sprintf("RetrCount: %-10d", br.RetrCount)
	s += fmt.Sprintf("TotalCount: %-10d", br.TotalCount)

	return s
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

	//br.BMTransLayer.SetData(r)

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

type DelSection struct {
	Begin int
	End   int
}

func (ds *DelSection) String() string {
	s := fmt.Sprintf("Begin %-10d", ds.Begin)
	s += fmt.Sprintf("End %-10d", ds.End)

	return s
}

func (ds *DelSection) Pack() ([]byte, error) {
	var (
		r, tmp []byte
	)

	tmp = translayer.UInt32ToBuf(uint32(ds.Begin))
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(ds.End))
	r = append(r, tmp...)

	return r, nil

}

func (ds *DelSection) UnPack(data []byte) (int, error) {
	var (
		offset int
	)

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack begin error")
	}
	ds.Begin = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack begin error")
	}

	ds.End = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil

}

type DelSectionResult struct {
	DelSection
	ErroCode int //0 success,1 failure
}

func (dsr *DelSectionResult) String() string {
	s := dsr.DelSection.String()
	s += fmt.Sprintf("ErrCode %-10d", dsr.ErroCode)

	return s
}

func (dsr *DelSectionResult) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = dsr.DelSection.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(dsr.ErroCode))
	r = append(r, tmp...)

	return r, nil
}

func (dsr *DelSectionResult) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)

	ds := &dsr.DelSection

	of, err = ds.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack errcode error")
	}

	dsr.ErroCode = int(binary.BigEndian.Uint32(data[offset:]))

	offset += translayer.Uint32Size

	return offset, nil
}

type BPOPDelete struct {
	translayer.BMTransLayer
	Section []DelSection
	Sn      []byte
	Sig     []byte
}

func NewBPOPDelete() *BPOPDelete {
	bd := &BPOPDelete{}

	bmtl := translayer.NewBMTL(translayer.DELETE)

	bd.BMTransLayer = *bmtl

	return bd
}

func (bd *BPOPDelete) String() string {
	s := bd.BMTransLayer.String()
	s += fmt.Sprintf("Section coubt: %-10d\r\n", len(bd.Section))
	for i := 0; i < len(bd.Section); i++ {
		s += bd.Section[i].String()
		s += "\r\n"
	}
	s += "Sn: " + base58.Encode(bd.Sn) + "\r\n"
	s += "Sig: " + base58.Encode(bd.Sig)

	return s
}

func (bd *BPOPDelete) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp = translayer.UInt32ToBuf(uint32(len(bd.Section)))
	r = append(r, tmp...)
	for i := 0; i < len(bd.Section); i++ {
		sec := &bd.Section[i]
		tmp, err = sec.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)
	}

	tmp, err = PackShortBytes(bd.Sn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(bd.Sig)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	//bd.BMTransLayer.SetData(r)

	return bd.BMTransLayer.Pack()

}

func (bd *BPOPDelete) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)
	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack section error")
	}
	cnt := int(binary.BigEndian.Uint32(data[offset:]))

	offset += translayer.Uint32Size

	for i := 0; i < cnt; i++ {
		sec := &DelSection{}
		of, err = sec.UnPack(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += of
		bd.Section = append(bd.Section, *sec)
	}
	bd.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	bd.Sig, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	return offset, nil
}

type BPOPDeleteResp struct {
	translayer.BMTransLayer
	Result []DelSectionResult
	Sn     []byte
}

func NewBPOPDeleteResp() *BPOPDeleteResp {
	bmtl := translayer.NewBMTL(translayer.DELETE_RESP)
	br := &BPOPDeleteResp{}
	br.BMTransLayer = *bmtl
	return br
}

func (br *BPOPDeleteResp) String() string {
	s := br.BMTransLayer.String()
	s += fmt.Sprintf("result count: %-10d\r\n", len(br.Result))
	for i := 0; i < len(br.Result); i++ {
		r := &br.Result[i]
		s += r.String()
		s += "\r\n"
	}
	s += "sn: " + base58.Encode(br.Sn)
	return s
}

func (bd *BPOPDeleteResp) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp = translayer.UInt32ToBuf(uint32(len(bd.Result)))
	r = append(r, tmp...)
	for i := 0; i < len(bd.Result); i++ {
		result := &bd.Result[i]
		tmp, err = result.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)
	}

	tmp, err = PackShortBytes(bd.Sn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	//bd.BMTransLayer.SetData(r)

	return bd.BMTransLayer.Pack()

}

func (bd *BPOPDeleteResp) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)
	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack result section error")
	}
	cnt := int(binary.BigEndian.Uint32(data[offset:]))

	offset += translayer.Uint32Size

	for i := 0; i < cnt; i++ {
		sec := &DelSectionResult{}
		of, err = sec.UnPack(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += of
		bd.Result = append(bd.Result, *sec)
	}
	bd.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	return offset, nil
}
