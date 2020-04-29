package test

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"math/rand"
	"testing"
)

func Test_FileProperty(t *testing.T) {
	hash1 := make([]byte, 16)

	for {
		n, _ := rand.Read(hash1)
		if n != len(hash1) {
			continue
		}
		break
	}

	fp := &bmprotocol.FileProperty{hash1, "hash1.txt", 0, 12221}

	data, _ := fp.Pack()
	fmt.Println(fp.String())

	//fmt.Println(hex.EncodeToString(data))

	fpunpack := &bmprotocol.FileProperty{}
	offset, err := fpunpack.UnPack(data)
	if err != nil {
		fmt.Println(offset, err)
	}

	fmt.Println(fpunpack.String())

	if fp.String() == fpunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}

func Test_Attachment(t *testing.T) {

	hash1 := make([]byte, 16)

	for {
		n, _ := rand.Read(hash1)
		if n != len(hash1) {
			continue
		}
		break
	}

	//hash2 := make([]byte, 16)
	//
	//for {
	//	n, _ := rand.Read(hash2)
	//	if n != len(hash2) {
	//		continue
	//	}
	//	break
	//}

	a := &bmprotocol.Attachment{"/smallfile/",
		bmprotocol.FileProperty{hash1, "hash1.doc", 0, 1002000}}
	data, _ := a.Pack()

	fmt.Println(a.String())

	aunpack := &bmprotocol.Attachment{}

	//bmtl:=translayer.BMTransLayer{}
	//offset,_:=bmtl.UnPack(data)

	aunpack.UnPack(data)
	fmt.Println(aunpack.String())

	if a.String() == aunpack.String() {
		t.Log("pass")
	} else {
		t.Fatal("failed")
	}

}
