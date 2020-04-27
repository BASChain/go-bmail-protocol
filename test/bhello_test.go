package test

import (
	"crypto"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"math/rand"
	"testing"
)

func Test_BHello(t *testing.T) {
	bmh := bmprotocol.NewBMHello()

	fmt.Println(bmh.String())

	data, _ := bmh.Pack()

	bmhUnPack := &bmprotocol.BMHello{}

	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPackHead(data)

	bmhUnPack.BMTransLayer = *bmtl

	bmhUnPack.UnPack(data[offset:])

	fmt.Println(bmhUnPack.String())

	if bmh.String() == bmhUnPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_BMHelloACK(t *testing.T) {
	sn := make([]byte, 32)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	bmha := bmprotocol.NewBMHelloACK(sn)

	data, _ := bmha.Pack()

	fmt.Println(bmha.String())

	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPackHead(data)

	bmhaUnPack := &bmprotocol.BMHelloACK{}
	bmhaUnPack.BMTransLayer = *bmtl

	bmhaUnPack.UnPack(data[offset:])

	fmt.Println(bmhaUnPack.String())

	if bmha.String() == bmhaUnPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_SendSignature(t *testing.T) {
	rsapriv, _ := rsa.GenerateKey(crand.Reader, 2048)

	pub := &rsapriv.PublicKey

	sn := make([]byte, 32)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	ss := bmprotocol.NewSendSignature(sn, "admin@bas")

	fsig := ss.ForSigBuf()

	hashed := sha256.Sum256(fsig)

	sig, _ := rsa.SignPKCS1v15(crand.Reader, rsapriv, crypto.SHA256, hashed[:])

	ss.SetSig(sig)

	data, _ := ss.Pack()

	fmt.Println(ss.String())

	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPackHead(data)

	ssUnPack := &bmprotocol.SendSignature{}

	ssUnPack.BMTransLayer = *bmtl

	ssUnPack.UnPack(data[offset:])

	fmt.Println(ssUnPack.String())

	hashed = sha256.Sum256(ss.ForSigBuf())

	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], ssUnPack.GetSig()); err != nil {
		t.Fatal("failed")
	}

	if ss.String() == ssUnPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}
}

func Test_ValidateSignature(t *testing.T) {

	sn := make([]byte, 32)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	vs := bmprotocol.NewValidSign(sn, bmprotocol.Validate_Success)

	data, _ := vs.Pack()

	fmt.Println(vs.String())

	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPackHead(data)

	vsUnPack := &bmprotocol.ValidateSignature{}
	vsUnPack.BMTransLayer = *bmtl
	vsUnPack.UnPack(data[offset:])

	fmt.Println(vsUnPack.String())

	if vs.String() == vsUnPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}
}
