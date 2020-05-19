package test

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"math/rand"
	"testing"
)

func fillEH(eh *bmprotocol.EnvelopeRoute) {
	eh.From = "a@bas"
	eh.RecpAddr = "b@bas"

	pubkey := make([]byte, 32)

	for {
		n, _ := rand.Read(pubkey)
		if n != len(pubkey) {
			continue
		}
		break
	}
	eh.LPubKey = pubkey
}

func fillET(et *bmprotocol.EnvelopeSig) {
	iv := make([]byte, 16)

	for {
		n, _ := rand.Read(iv)
		if n != len(iv) {
			continue
		}
		break
	}

	sig := make([]byte, 32)

	for {
		n, _ := rand.Read(sig)
		if n != len(sig) {
			continue
		}
		break
	}

	et.Sn = iv
	et.Sig = sig
}

func fillEC(ec *bmprotocol.EnvelopeContent) {
	ec.To = []string{"toa@bas", "tob@bas", "toc@bas"}
	ec.CC = []string{"cca@bas", "ccb@bas"}
	ec.BC = []string{"bca@bas"}

	ec.Subject = "test a ec"
	ec.Data = "test e content"

	hash1 := make([]byte, 16)

	for {
		n, _ := rand.Read(hash1)
		if n != len(hash1) {
			continue
		}
		break
	}

	hash2 := make([]byte, 16)

	for {
		n, _ := rand.Read(hash2)
		if n != len(hash2) {
			continue
		}
		break
	}

	ec.Files = []bmprotocol.Attachment{{"", bmprotocol.FileProperty{hash1, "name.doc", 0, 10200}},
		{"", bmprotocol.FileProperty{hash2, "name2.xls", 1, 20400}}}
}

func Test_SendEnvelope(t *testing.T) {
	se := bmprotocol.NewSendEnvelope()

	eh := &se.Envelope.EnvelopeHead

	fillEH(eh)

	et := &se.Envelope.EnvelopeSig

	fillET(et)

	ec := &se.Envelope.EnvelopeContent

	fillEC(ec)

	data, _ := se.Pack()

	se.BMTransLayer.SetDataLen(uint32(len(data) - translayer.BMHeadSize()))

	fmt.Println(se.String())

	seUnpack := &bmprotocol.SendEnvelope{}

	bmtl := &translayer.BMTransLayer{}
	n, _ := bmtl.UnPack(data)

	seUnpack.BMTransLayer = *bmtl

	seUnpack.UnPack(data[n:])

	fmt.Println(seUnpack.String())

	if se.String() == seUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_RespSendEnvelope(t *testing.T) {
	rse := bmprotocol.NewRespSendEnvelope()

	sn := make([]byte, 16)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	newsn := make([]byte, 16)

	for {
		n, _ := rand.Read(newsn)
		if n != len(newsn) {
			continue
		}
		break
	}

	bid := make([]byte, 16)

	for {
		n, _ := rand.Read(bid)
		if n != len(bid) {
			continue
		}
		break
	}

	rse.Sn = sn
	copy(rse.EId[:], bid)
	rse.NewSn = newsn
	rse.ErrId = 1

	data, _ := rse.Pack()
	rse.BMTransLayer.SetDataLen(uint32(len(data) - translayer.BMHeadSize()))

	fmt.Println(rse.String())

	rseUnpack := &bmprotocol.RespSendEnvelope{}

	bmtl := &translayer.BMTransLayer{}
	n, _ := bmtl.UnPack(data)

	rseUnpack.BMTransLayer = *bmtl

	rseUnpack.UnPack(data[n:])

	fmt.Println(rseUnpack.String())

	if rse.String() == rseUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_SendCryptEnvelope(t *testing.T) {
	sce := bmprotocol.NewSendCryptEnvelope()

	es := &sce.CryptEnvelope.EnvelopeSig
	fillET(es)
	eh := &sce.CryptEnvelope.EnvelopeHead
	fillEH(eh)

	newsn := make([]byte, 64)

	for {
		n, _ := rand.Read(newsn)
		if n != len(newsn) {
			continue
		}
		break
	}

	sce.CryptEnvelope.CipherTxt = newsn

	data, _ := sce.Pack()

	sce.BMTransLayer.SetDataLen(uint32(len(data) - translayer.BMHeadSize()))

	fmt.Println(sce.String())

	sceUnpack := &bmprotocol.SendCryptEnvelope{}
	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPack(data)

	sceUnpack.BMTransLayer = *bmtl

	sceUnpack.UnPack(data[offset:])

	fmt.Println(sceUnpack.String())

	if sce.String() == sceUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("error")
	}

}
