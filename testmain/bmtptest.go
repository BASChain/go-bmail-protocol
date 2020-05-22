package main

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas-mail-server/bmailcrypt"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-protocol/bmpclient2"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"math/rand"
	"net"
)

func main() {
	c := bmpclient2.NewClient2(net.ParseIP("39.99.198.143"), 100)
	if c == nil {
		fmt.Println("connect to peer error")
		return
	}
	defer c.Close()

	err := c.Helo()
	if err != nil {
		fmt.Println(err)
		//fmt.Println(err)
		return
	}

	fmt.Println("get hello ack:", base58.Encode(c.GetSn()))
	se := NewEnv(c, c.GetSn())
	//copy(se.SN[:],c.GetSn())

	resp, err1 := c.SendEnvelope(se)

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	if !bmailcrypt.Verify(c.SrvPk, c.Hash, resp.Sig) {
		fmt.Println("not a correct server")
	} else {
		fmt.Println("you bmail have send to a correct server")
	}

	if resp != nil {
		jstr, err := json.Marshal(*resp)
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Println(string(jstr))
		}

	}

}

func NewAddr(cnt int) []byte {
	sn := make([]byte, cnt)

	for {
		n, _ := rand.Read(sn)
		if n != len(sn) {
			continue
		}
		break
	}

	return sn
}

func fillEH(c *bmpclient2.BMClient2, eh *bmp.EnvelopeHead) {
	eh.From = "testa@eth"
	eh.To = "testb@eth"

	eh.FromAddr = bmail.ToAddress(c.PK[:])
	fmt.Println("from addr", eh.FromAddr)

	eh.Eid, _ = uuid.FromBytes(NewAddr(16))
	eh.ToAddr = "BM7JNBrt8SQX4AGc5fvkjJ9p2bwTt5Wyxnz6af22iHgh2p"
	fmt.Println("to addr", eh.ToAddr)
	copy(eh.IV[:], NewAddr(16))

}

//
//func fillET(et *bmprotocol.EnvelopeSig) {
//	iv := make([]byte, 16)
//
//	for {
//		n, _ := rand.Read(iv)
//		if n != len(iv) {
//			continue
//		}
//		break
//	}
//
//	sig := make([]byte, 32)
//
//	for {
//		n, _ := rand.Read(sig)
//		if n != len(sig) {
//			continue
//		}
//		break
//	}
//
//	et.Sn = iv
//	et.Sig = sig
//}

func fillEC(ec *bmp.EnvelopeBody) {

	ec.Subject = "test a ec"
	ec.MsgBody = "test e content"

}

func NewEnv(c *bmpclient2.BMClient2, sn []byte) *bmp.EnvelopeSyn {
	se := bmp.CryptEnvelope{}

	eh := &se.EnvelopeHead

	fillEH(c, eh)

	ec := &bmp.EnvelopeBody{}
	fillEC(ec)

	ecdata, _ := json.Marshal(*ec)

	aesk, _ := bmailcrypt.GenerateAesKey(bmail.Address(eh.ToAddr).ToPubKey(), c.Priv)

	cdata, _ := bmailcrypt.EncryptWithIV(aesk, eh.IV[:], ecdata)

	//todo...
	se.CryptData = cdata

	es := &bmp.EnvelopeSyn{}
	es.Env = &se

	es.Mode = 1
	copy(es.SN[:], sn)

	//osn:=tools.NewSn(128)
	es.Sig = ed25519.Sign(c.Priv, sn)
	//fmt.Println("envsyn sn:",base58.Encode(es.SN[:]))
	//fmt.Println("env sig:",base58.Encode(es.Sig))

	//r:=ed25519.Verify(c.PK,sn,es.Sig)
	//
	//
	//
	//fmt.Println("env sig verify:",r)

	es.Hash = se.Hash()

	c.Hash = es.Hash

	return es
}
