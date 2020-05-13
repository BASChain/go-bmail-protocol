package bmp

import (
	"crypto/sha256"
	"encoding/json"
)

type CryptEnvelope struct {
	EnvelopeHead
	CryptData []byte `json:"cryptBody"`
}

func (c *CryptEnvelope) Hash() []byte {
	data, _ := json.Marshal(c)
	hash := sha256.Sum256(data)
	return hash[:]
}
