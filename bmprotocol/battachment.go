package bmprotocol

import (
	"encoding/binary"
	"fmt"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"io"
	"os"
)

type FileProperty struct {
	Hash      []byte
	FileName  string
	FileType  int
	IsEnCrypt bool
	FileSize  int
}

func (fp *FileProperty) String() string {
	s := fmt.Sprintf("Hash: %-50s", base58.Encode(fp.Hash))
	s += fmt.Sprintf("FileName: %-20s", fp.FileName)
	s += fmt.Sprintf("FileType: %-8d", fp.FileType)
	s += fmt.Sprintf("IsEnCrypt: %-8t", fp.IsEnCrypt)
	s += fmt.Sprintf("FileSize: %-12d", fp.FileSize)

	return s
}

func (fp *FileProperty) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortBytes(fp.Hash)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortString(fp.FileName)
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(fp.FileType))
	r = append(r, tmp...)

	var isEnc uint8

	if fp.IsEnCrypt {
		isEnc = 1
	}

	r = append(r, isEnc)

	tmp = translayer.UInt32ToBuf(uint32(fp.FileSize))
	r = append(r, tmp...)

	return r, nil
}

func (fp *FileProperty) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)
	fp.Hash, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	fp.FileName, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack FileType error")
	}

	fp.FileType = int(binary.BigEndian.Uint32(data[offset:]))

	offset += translayer.Uint32Size

	if len(data) < offset+translayer.Uin8Size {
		return 0, errors.New("unpack IsEncrypt error")
	}

	if data[offset] == 1 {
		fp.IsEnCrypt = true
	} else {
		fp.IsEnCrypt = false
	}
	offset += translayer.Uin8Size

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack FileSize error")
	}

	fp.FileSize = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	return offset, nil

}

type Attachment struct {
	Path string //if less 30MB, path is none,
	FileProperty
}

func (a *Attachment) String() string {
	s := "Path: " + a.Path + "\r\n"
	s += a.FileProperty.String()

	return s
}

func (a *Attachment) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortString(a.Path)
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	fp := &a.FileProperty

	tmp, err = fp.Pack()

	r = append(r, tmp...)

	return r, nil
}

func (a *Attachment) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)
	a.Path, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	fp := &a.FileProperty

	of, err = fp.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil

}

type SendAttachment struct {
	translayer.BMTransLayer
	FileProperty
	EnvelopeSig
	EId  translayer.EnveUniqID
	File *os.File
}

type SAReader struct {
	head      []byte
	Reader    io.Reader
	readSize  int
	totalSize int
}

func NewSendAttachment() *SendAttachment {
	bmtl := translayer.NewBMTL(translayer.SEND_ATTACHMNENT)

	sa := &SendAttachment{}
	sa.BMTransLayer = *bmtl

	return sa
}

func (sa *SendAttachment) String() string {
	s := sa.BMTransLayer.String()
	s += sa.FileProperty.String()
	s += sa.EnvelopeSig.String()
	s += fmt.Sprintf("%-30s", base58.Encode(sa.EId[:]))

	return s
}

func (sar *SAReader) GetTotalSize() int {
	return sar.totalSize
}

func (sar *SAReader) Read(p []byte) (n int, err error) {

	var toread int

	if sar.readSize < len(sar.head) {
		toread = len(sar.head) - sar.readSize
		if toread > len(p) {
			toread = len(p)
			copy(p, sar.head)
			sar.readSize += toread
			return toread, nil
		} else {
			copy(p, sar.head)
		}

	}

	if sar.totalSize <= sar.readSize {
		return toread, nil
	}

	nn, err := sar.Reader.Read(p[toread:])
	n = nn + toread

	sar.readSize += n

	return

}

func (sa *SendAttachment) GetReader() (*SAReader, error) {
	var (
		r, tmp []byte
		err    error
	)

	r = NewHeadBuf()

	fp := &sa.FileProperty

	tmp, err = fp.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	es := &sa.EnvelopeSig
	tmp, err = es.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(sa.EId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	AddPackHead(&(sa.BMTransLayer), r)

	sar := &SAReader{}

	sar.head = r
	sar.totalSize = len(r) + sa.FileSize
	sar.Reader = sa.File

	return sar, nil
}

type RespSendAttachment struct {
	translayer.BMTransLayer
	FileProperty
	Sn    []byte
	NewSn []byte
	EId   translayer.EnveUniqID
	ErrId int
}

func NewRespSendAttachment() *RespSendAttachment {
	bmtl := translayer.NewBMTL(translayer.RESP_ATTACHMENT)

	rsa := &RespSendAttachment{}

	rsa.BMTransLayer = *bmtl

	return rsa
}

func (rsa *RespSendAttachment) String() string {
	s := rsa.BMTransLayer.String()
	s += rsa.FileProperty.String()
	s += fmt.Sprintf("%-30s", base58.Encode(rsa.Sn))
	s += fmt.Sprintf("%-30s", base58.Encode(rsa.NewSn))
	s += fmt.Sprintf("%-30s", base58.Encode(rsa.EId[:]))
	s += fmt.Sprintf("%d", rsa.ErrId)

	return s

}

func (rsa *RespSendAttachment) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)
	r = NewHeadBuf()

	fp := &rsa.FileProperty

	tmp, err = fp.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rsa.Sn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rsa.NewSn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(rsa.EId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(rsa.ErrId))
	r = append(r, tmp...)

	AddPackHead(&(rsa.BMTransLayer), r)

	return r, nil
}

func (rsa *RespSendAttachment) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	fp := &rsa.FileProperty
	of, err = fp.UnPack(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	rsa.Sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	rsa.NewSn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}

	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	copy(rsa.EId[:], tmp)

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack errid error")
	}
	rsa.ErrId = int(binary.BigEndian.Uint32(data[offset:]))

	offset += of

	return offset, nil

}
