package test

import (
	crand "crypto/rand"
	"crypto/rsa"
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

	offset, _ := bmtl.UnPack(data)

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

	fmt.Println(bmha.String())

	data, _ := bmha.Pack()

	//fmt.Println(hex.EncodeToString(data))

	bmtl := &translayer.BMTransLayer{}

	offset, _ := bmtl.UnPack(data)

	fmt.Println(offset)

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

	fmt.Println(pub)

}

func Test_ValidateSignature(t *testing.T) {

}
