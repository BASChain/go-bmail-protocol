package bmprotocol

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"time"
)

//client --helo--> server
//server -->helo_resp{sn}-->client

func GetNowMsTime() int64 {
	return time.Now().UnixNano() / 1e6
}

type BMHello struct {
	translayer.BMTransLayer
}

func NewBMHello() *BMHello {
	bmtl := translayer.NewBMTL(translayer.HELLO)
	bmh := &BMHello{}
	bmh.BMTransLayer = *bmtl

	return bmh
}

func (bmh *BMHello) Pack() ([]byte, error) {
	return bmh.BMTransLayer.Pack()
}

func (bmh *BMHello) UnPack(data []byte) (int, error) {
	//nothing todo
	return 0, nil
}

func (bmh *BMHello) String() string {
	return bmh.BMTransLayer.String()
}

type BMHelloACK struct {
	translayer.BMTransLayer
	sn []byte
}

func NewBMHelloACK(sn []byte) *BMHelloACK {
	bmact := translayer.NewBMTL(translayer.HELLO_ACK)

	bmhack := &BMHelloACK{}

	bmhack.BMTransLayer = *bmact

	bmhack.sn = sn

	return bmhack
}

func (bmha *BMHelloACK) Pack() ([]byte, error) {

	var (
		tmp []byte
		err error
	)

	r := NewHeadBuf()

	tmp, err = PackShortBytes(bmha.sn)
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return AddPackHead(&(bmha.BMTransLayer), r)
}

func (bmha *BMHelloACK) String() string {
	s := bmha.BMTransLayer.String()

	s += fmt.Sprintf("sn: %s", base58.Encode(bmha.sn))

	return s
}

func (bmha *BMHelloACK) UnPack(data []byte) (int, error) {

	var (
		of  int
		err error
	)

	bmha.sn, of, err = UnPackShortBytes(data)
	if err != nil {
		//fmt.Println("error is",err)
		return 0, err
	}

	return of, nil
}
