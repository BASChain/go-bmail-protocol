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
