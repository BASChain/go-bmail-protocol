package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-bas-mail-server/bmailcrypt"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bpop"
	"github.com/BASChain/go-bmail-protocol/bpopclient"
	"github.com/btcsuite/btcutil/base58"
	"github.com/kprc/nbsnetwork/tools"
	"net"
)

func main() {
	c := bpopclient.NewClient2(net.ParseIP("39.99.198.143"), 100)
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
	se := NewDownloadCmd(c, c.GetSn())

	resp, err1 := c.SendCommand(se)

	if err1 != nil {
		fmt.Println(err1)
		return
	}

	hash := resp.CmdCxt.Hash()
	if bytes.Compare(hash[:], resp.Hash) != 0 {
		fmt.Println("hash error")
		return
	}

	if !bmailcrypt.Verify(c.SrvPk, hash, resp.Sig) {
		fmt.Println("not a correct server")
	} else {
		fmt.Println("you bmail have send to a correct server")
	}

	cmdack := resp.CmdCxt.(*bpop.CmdDownloadAck)

	for i := 0; i < len(cmdack.CryptEps); i++ {
		cep := &cmdack.CryptEps[i]
		jstr, _ := json.Marshal(cep.EnvelopeHead)
		fmt.Println(string(jstr))

		aesk, _ := bmailcrypt.GenerateAesKey(bmail.Address(cep.FromAddr).ToPubKey(), c.Priv)
		plainsub, _ := bmailcrypt.Decrypt(aesk, cep.CryptSub)
		plainbody, _ := bmailcrypt.Decrypt(aesk, cep.CryptBody)
		fmt.Println(string(plainsub))
		fmt.Println(string(plainbody))
	}

}

func NewDownloadCmd(c *bpopclient.BMClient2, sn []byte) *bpop.CommandSyn {

	csyn := &bpop.CommandSyn{}

	cdl := &bpop.CmdDownload{}

	csyn.Cmd = cdl

	copy(csyn.SN[:], sn)
	csyn.Sig = ed25519.Sign(c.Priv, sn)

	cdl.MailCnt = 5

	cdl.BeforeTime = tools.GetNowMsTime()
	cdl.Owner = bmail.ToAddress(c.PK)
	cdl.MailAddr = "testb@eth"

	return csyn
}
