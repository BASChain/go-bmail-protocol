package test

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"math/rand"
	"testing"
)

func Test_EnvelopeHead(t *testing.T) {
	eh := &bmprotocol.EnvelopeHead{}

	eh.From = "a@bas"
	eh.RecpAddr = "b@bas"
	eh.RecpAddrType = 1

	pubkey := make([]byte, 32)

	for {
		n, _ := rand.Read(pubkey)
		if n != len(pubkey) {
			continue
		}
		break
	}
	eh.LPubKey = pubkey

	data, _ := eh.Pack()

	fmt.Println(eh.String())

	ehUnpack := &bmprotocol.EnvelopeHead{}
	ehUnpack.UnPack(data)

	fmt.Println(ehUnpack.String())

	if eh.String() == ehUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_EnvelopeContent(t *testing.T) {
	ec := &bmprotocol.EnvelopeContent{}

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

	data, _ := ec.Pack()

	fmt.Println(ec.String())

	ecUnpack := &bmprotocol.EnvelopeContent{}
	ecUnpack.UnPack(data)
	fmt.Println(ecUnpack.String())

	if ec.String() == ecUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}
}

func Test_EnvelopeTail(t *testing.T) {
	et := &bmprotocol.EnvelopeTail{}

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

	et.IV = iv
	et.Sig = sig

	data, _ := et.Pack()
	fmt.Println(et.String())

	etUnpack := &bmprotocol.EnvelopeTail{}
	etUnpack.UnPack(data)

	fmt.Println(etUnpack.String())

	if et.String() == etUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_Envelop(t *testing.T) {
	e := &bmprotocol.Envelope{}

	eh := &e.EnvelopeHead

	eh.From = "a@bas"
	eh.RecpAddr = "b@bas"
	eh.RecpAddrType = 1

	pubkey := make([]byte, 32)

	for {
		n, _ := rand.Read(pubkey)
		if n != len(pubkey) {
			continue
		}
		break
	}
	eh.LPubKey = pubkey

	ec := &e.EnvelopeContent
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

	et := &e.EnvelopeTail
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

	et.IV = iv
	et.Sig = sig

	data, _ := e.Pack()
	fmt.Println(e.String())

	eunPack := &bmprotocol.Envelope{}
	eunPack.UnPack(data)

	fmt.Println(eunPack.String())

	if e.String() == eunPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_CryptEnvelope(t *testing.T) {
	crypttxt := make([]byte, 64)

	for {
		n, _ := rand.Read(crypttxt)
		if n != len(crypttxt) {
			continue
		}
		break
	}

	ce := &bmprotocol.CryptEnvelope{}
	ce.CipherTxt = crypttxt
	eh := &ce.EnvelopeHead
	eh.From = "a@bas"
	eh.RecpAddr = "b@bas"
	eh.RecpAddrType = 1

	pubkey := make([]byte, 32)

	for {
		n, _ := rand.Read(pubkey)
		if n != len(pubkey) {
			continue
		}
		break
	}
	eh.LPubKey = pubkey

	et := &ce.EnvelopeTail

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

	et.IV = iv
	et.Sig = sig

	fmt.Println(ce.String())

	data, _ := ce.Pack()

	ceUnpack := &bmprotocol.CryptEnvelope{}
	ceUnpack.UnPack(data)

	fmt.Println(ceUnpack.String())

	if ce.String() == ceUnpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}
