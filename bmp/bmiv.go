package bmp

import (
	"crypto/aes"
	"crypto/rand"
	"io"
)

type BMailIV [aes.BlockSize]byte

func NewIV() (*BMailIV, error) {
	iv := new(BMailIV) //make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		return nil, err
	}
	return iv, nil
}

func (iv *BMailIV) Bytes() []byte {
	return iv[:]
}

const BMailSNSize = 16

type BMailSN [BMailSNSize]byte

func (sn *BMailSN) Bytes() []byte {
	return sn[:]
}
