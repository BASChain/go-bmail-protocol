package test

import (
	"testing"
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
)

func Test_BPOPStat(t *testing.T)  {
	b:=bmprotocol.NewBPOPStat()

	fmt.Println(b.String())

	data,_:=b.Pack()

	bunPack := &bmprotocol.BPOPStat{}

	bmtl:=&translayer.BMTransLayer{}
	offset,_:=bmtl.UnPack(data)

	bunPack.BMTransLayer = *bmtl
	bunPack.UnPack(data[offset:])

	fmt.Println(bunPack.String())

	if b.String() == bunPack.String(){
		t.Log("pass")
	}else{
		t.Fatal("failed")
	}
}

func Test_BPOPStatResp(t *testing.T)  {
	bs:=bmprotocol.NewBPOPStatResp()

	bs.Total = 100
	bs.Received = 10
	bs.TotalSpaceBytes = 1000000000
	bs.TotalStoredBytes = 200000000


	data,_:=bs.Pack()
	fmt.Println(bs.String())

	bsUnPck:=&bmprotocol.BPOPStatResp{}
	bmtl:=&translayer.BMTransLayer{}



	offset,_:=bmtl.UnPackHead(data)

	bsUnPck.BMTransLayer= *bmtl


	bsUnPck.UnPack(data[offset:])

	fmt.Println(bsUnPck.String())

	if bs.String() == bsUnPck.String(){
		t.Log("pass")
	}else{
		t.Fatal("failed")
	}


}

func Test_BPOPList(t *testing.T)  {
	bl:=bmprotocol.NewBPOPList()

	bl.BeginID = 20
	bl.ListCount = 15

	data,_:=bl.Pack()

	fmt.Println(bl.String())

	blunpack := &bmprotocol.BPOPList{}
	bmtl:=&translayer.BMTransLayer{}
	offset,_:=bmtl.UnPack(data)
	blunpack.BMTransLayer = *bmtl

	blunpack.UnPack(data[offset:])

	fmt.Println(blunpack.String())
	if bl.String() == blunpack.String(){
		t.Log("pass")
	}else{
		t.Fatal("failed")
	}

}


