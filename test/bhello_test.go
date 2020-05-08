package test

import (
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
	sn := make([]byte, 16)

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

	offset, _ := bmtl.UnPack(data)

	bmhaUnPack := &bmprotocol.BMHelloACK{}
	bmhaUnPack.BMTransLayer = *bmtl

	//fmt.Println(offset)

	bmhaUnPack.UnPack(data[offset:])

	fmt.Println(bmhaUnPack.String())

	if bmha.String() == bmhaUnPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

