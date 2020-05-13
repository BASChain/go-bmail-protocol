package bmp

type CryptEnvelope struct {
	EnvelopeHead
	CryptData []byte `json:"cryptBody"`
}
