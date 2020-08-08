package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/realbmail/go-bas-mail-server/bmailcrypt"
	"github.com/realbmail/go-bmail-account"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-bmail-protocol/bmpclient2"
	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"golang.org/x/crypto/curve25519"
	"math/rand"
	"net"
	"time"
)

func main() {
	c := bmpclient2.NewClient2(net.ParseIP("34.92.157.168"), 100)
	if c == nil {
		fmt.Println("connect to peer error")
		return
	}
	defer c.Close()

	err := c.Helo()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("get hello ack:", base58.Encode(c.GetSn()))
	se := NewEnv(c, c.GetSn())

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

	currentTime := time.Now().UnixNano()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(currentTime))

	copy(sn, buf)

	return sn
}

//
//func fillEH(c *bmpclient2.BMClient2, eh *bmp.EnvelopeHead) {
//	eh.From = "testa@eth"
//	eh.To = "testb@eth"
//
//	eh.FromAddr = bmail.ToAddress(c.PK[:])
//	fmt.Println("from addr", eh.FromAddr)
//
//	eh.Eid, _ = uuid.FromBytes(NewAddr(16))
//	eh.ToAddr = "BM7JNBrt8SQX4AGc5fvkjJ9p2bwTt5Wyxnz6af22iHgh2p"
//	fmt.Println("to addr", eh.ToAddr)
//	copy(eh.IV[:], NewAddr(16))
//	eh.Date = time.Duration(time.Now().UnixNano() / 1e6)
//
//}
//
//func fillEC(ec *bmp.EnvelopeBody) {
//
//	ec.Subject = "test a ec"
//	ec.MsgBody = "test e content"
//
//}

func NewEnv(c *bmpclient2.BMClient2, sn []byte) *bmp.EnvelopeSyn {
	se := bmp.BMailEnvelope{}

	se.FromName = "testa@eth"
	eid, _ := uuid.FromBytes(NewAddr(16))
	se.Eid = eid.String()
	se.FromAddr = bmail.ToAddress(c.PK[:])
	se.RCPTs = []*bmp.Recipient{&bmp.Recipient{ToName: "testb@eth",
		ToAddr: "BM7JNBrt8SQX4AGc5fvkjJ9p2bwTt5Wyxnz6af22iHgh2p", RcptType: bmp.RcpTypeTo}}

	se.SessionID = base64.StdEncoding.EncodeToString(NewAddr(16))

	se.Subject = "this is a test crypt bmail"
	se.MailBody = "hello bmail, it's a bmail body"

	//eh := &se.EnvelopeHead
	//
	//fillEH(c, eh)
	//
	//ec := &bmp.EnvelopeBody{}
	//fillEC(ec)

	//ecdata, _ := json.Marshal(*ec)

	aesk, _ := bmailcrypt.GenerateAesKey(bmail.Address(se.RCPTs[0].ToAddr).ToPubKey(), c.Priv)

	uniqKey := NewAddr(curve25519.ScalarSize)

	cdata, _ := bmailcrypt.Encrypt(aesk, uniqKey)

	se.RCPTs[0].AESKey = cdata

	sdata, _ := bmailcrypt.Encrypt(uniqKey, []byte(se.Subject))

	se.Subject = base64.StdEncoding.EncodeToString(sdata)

	cxtdata, _ := bmailcrypt.Encrypt(uniqKey, []byte(se.MailBody))
	se.MailBody = base64.StdEncoding.EncodeToString(cxtdata)

	es := &bmp.EnvelopeSyn{}
	es.Env = &se

	copy(es.SN[:], sn)

	es.Sig = ed25519.Sign(c.Priv, sn)

	es.Hash = se.Hash()

	c.Hash = es.Hash

	return es
}
