package bmprotocol

import (
	"encoding/binary"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"io"

	"fmt"
	"os"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
)

type FileProperty struct {
	Hash     []byte
	FileName string
	FileType int
	FileSize int
}

func (fp *FileProperty) String() string {
	s := fmt.Sprintf("Hash: %-50s", base58.Encode(fp.Hash))
	s += fmt.Sprintf("FileName: %-20s", fp.FileName)
	s += fmt.Sprintf("FileType: %-8d", fp.FileType)
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
	File *os.File
}

type SAReader struct {
	head      []byte
	Reader    io.Reader
	readSize  int
	totalSize int
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

	fp := &sa.FileProperty

	tmp, err = fp.Pack()
	if err != nil {
		return nil, err
	}

	sa.BMTransLayer.SetDataLen(uint32(len(tmp)))

	r, err = sa.BMTransLayer.Pack()

	sar := &SAReader{}

	sar.head = r
	sar.totalSize = len(r) + sa.FileSize
	sar.Reader = sa.File

	return sar, nil
}
