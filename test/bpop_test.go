package test

import (
	"crypto/rand"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"testing"
)

func Test_BPOPStat(t *testing.T) {
	b := bmprotocol.NewBPOPStat()

	fmt.Println(b.String())

	data, _ := b.Pack()

	bunPack := &bmprotocol.BPOPStat{}

	bmtl := &translayer.BMTransLayer{}
	offset, _ := bmtl.UnPack(data)

	bunPack.BMTransLayer = *bmtl
	bunPack.UnPack(data[offset:])

	fmt.Println(bunPack.String())

	if b.String() == bunPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}
}

func Test_BPOPStatResp(t *testing.T) {
	bs := bmprotocol.NewBPOPStatResp()

	bs.Total = 100
	bs.Received = 10
	bs.TotalSpaceBytes = 1000000000
	bs.TotalStoredBytes = 200000000

	data, _ := bs.Pack()
	fmt.Println(bs.String())

	bsUnPck := &bmprotocol.BPOPStatResp{}
	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPackHead(data)

	bsUnPck.BMTransLayer = *bmtl

	bsUnPck.UnPack(data[offset:])

	fmt.Println(bsUnPck.String())

	if bs.String() == bsUnPck.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_BPOPList(t *testing.T) {
	bl := bmprotocol.NewBPOPList()

	bl.BeginID = 20
	bl.ListCount = 15

	data, _ := bl.Pack()

	fmt.Println(bl.String())

	blunpack := &bmprotocol.BPOPList{}
	bmtl := &translayer.BMTransLayer{}
	offset, _ := bmtl.UnPack(data)
	blunpack.BMTransLayer = *bmtl

	blunpack.UnPack(data[offset:])

	fmt.Println(blunpack.String())
	if bl.String() == blunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_ListNode(t *testing.T) {
	ln := &bmprotocol.ListNode{}
	ln.SizeOfBytes = 20000000
	ln.ID = 20

	data, _ := ln.Pack()

	fmt.Println(ln.String())

	lnunpack := &bmprotocol.ListNode{}

	lnunpack.UnPack(data)

	fmt.Println(lnunpack.String())

	if ln.String() == lnunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}
}

func Test_BPOPListResp(t *testing.T) {
	bl := bmprotocol.NewBPOPListResp()

	bl.BeginID = 120
	bl.ListCount = 20

	ln1 := &bmprotocol.ListNode{ID: 120, SizeOfBytes: 100200}
	ln2 := &bmprotocol.ListNode{ID: 121, SizeOfBytes: 10020022}

	bl.Nodes = append(bl.Nodes, ln1)
	bl.Nodes = append(bl.Nodes, ln2)

	data, _ := bl.Pack()
	fmt.Println(bl.String())

	//fmt.Println(hex.EncodeToString(data))

	blunpack := &bmprotocol.BPOPListResp{}
	bmtl := &translayer.BMTransLayer{}
	offset, _ := bmtl.UnPackHead(data)

	blunpack.BMTransLayer = *bmtl

	_, err := blunpack.UnPack(data[offset:])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(blunpack.String())

	if bl.String() == blunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_BPOPRetr(t *testing.T) {
	br := bmprotocol.NewBPOPRetr()
	br.BeginID = 110
	br.RetrCount = 20

	data, _ := br.Pack()
	fmt.Println(br.String())

	brunpack := &bmprotocol.BPOPRetr{}
	bmtl := &translayer.BMTransLayer{}
	offset, _ := bmtl.UnPack(data)
	brunpack.BMTransLayer = *bmtl
	brunpack.UnPack(data[offset:])

	fmt.Println(brunpack.String())
	if br.String() == brunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func newCE() *bmprotocol.CryptEnvelope {
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

	return ce
}

func Test_BPOPRetrResp(t *testing.T) {
	br := bmprotocol.NewBPOPRetrResp()

	br.Mails = append(br.Mails, *newCE())
	br.Mails = append(br.Mails, *newCE())
	br.Mails = append(br.Mails, *newCE())

	br.BeginID = 102
	br.RetrCount = 20

	br.TotalCount = 105

	data, _ := br.Pack()
	fmt.Println(br.String())

	brunpack := &bmprotocol.BPOPRetrResp{}
	bmtl := &translayer.BMTransLayer{}
	offset, _ := bmtl.UnPack(data)

	brunpack.BMTransLayer = *bmtl
	brunpack.UnPack(data[offset:])

	fmt.Println(brunpack.String())

	if br.String() == brunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}
