package main

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmclient"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"math/rand"
	"net"
)

func main() {
	c := bmclient.NewClient(net.ParseIP("39.99.198.143"), 100)
	if c == nil {
		fmt.Println("connect to peer error")
		return
	}
	defer c.Close()

	err := c.HeloSendAndRcv()
	if err != nil {
		fmt.Println(err)
		return
	}

	se := NewEnv()
	se.Sn = c.GetSn()

	resp, err1 := c.SendEnvelope(se)

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	if resp != nil {
		fmt.Println(resp.String())
	}

}

func fillEH(eh *bmprotocol.EnvelopeHead) {
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

	bid := make([]byte, 16)

	for {
		n, _ := rand.Read(bid)
		if n != len(bid) {
			continue
		}
		break
	}

	copy(eh.EId[:], bid)

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

	ec.Files = []bmprotocol.Attachment{{"", bmprotocol.FileProperty{hash1, "name.doc", 0, true, 10200}},
		{"", bmprotocol.FileProperty{hash2, "name2.xls", 1, true, 20400}}}
}

func NewEnv() *bmprotocol.SendEnvelope {
	se := bmprotocol.NewSendEnvelope()

	eh := &se.Envelope.EnvelopeHead

	fillEH(eh)

	et := &se.Envelope.EnvelopeSig

	fillET(et)

	ec := &se.Envelope.EnvelopeContent

	fillEC(ec)

	return se
}
