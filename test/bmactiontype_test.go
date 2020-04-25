package test

import (
	"crypto/rand"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"testing"
)

func Test_BMTransLayer(t *testing.T) {
	data := make([]byte, 32)

	for {
		n, _ := rand.Read(data)
		if n != len(data) {
			continue
		}
		break
	}

	bmtl := translayer.NewBMTL(translayer.HELLO, data)

	fmt.Println(bmtl.String())

	packData, _ := bmtl.Pack()

	bmtlUnPack := &translayer.BMTransLayer{}

	bmtlUnPack.UnPack(packData)

	fmt.Println(bmtlUnPack.String())

	if bmtl.String() == bmtlUnPack.String() {
		t.Log("pass")
	} else {
		t.Fatal("Failed")
	}

}
