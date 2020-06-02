package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
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
	c := bpopclient.NewClient2(net.ParseIP("34.92.157.168"), 100)
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

	//if len(os.Args)>1{
	//	t:=strconv.Atoi(os.Args[1])
	//}

	fmt.Println("get hello ack:", base58.Encode(c.GetSn()))
	se := NewDownloadCmd(c, c.GetSn(), tools.GetNowMsTime())

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
		cep := cmdack.CryptEps[i]
		jstr, _ := json.Marshal(*cep)
		fmt.Println(string(jstr))

		aesk, _ := bmailcrypt.GenerateAesKey(bmail.Address(cep.FromAddr).ToPubKey(), c.Priv)
		var uniqk []byte
		for j := 0; j < len(cep.RCPTs); j++ {
			rcpt := cep.RCPTs[j]
			if rcpt.ToName == se.Cmd.(*bpop.CmdDownload).MailAddr {
				uniqk = rcpt.AESKey
			}
		}

		if len(uniqk) == 0 {
			fmt.Println("Decrypt break")
			continue
		}

		//if se.Cmd.(*bpop.CmdDownload).MailAddr

		plainkey, _ := bmailcrypt.Decrypt(aesk, uniqk)

		subject, _ := base64.StdEncoding.DecodeString(cep.Subject)
		contxt, _ := base64.StdEncoding.DecodeString(cep.MailBody)
		plainsub, _ := bmailcrypt.Decrypt(plainkey, subject)
		plainbody, _ := bmailcrypt.Decrypt(plainkey, contxt)
		fmt.Println(string(plainsub))
		fmt.Println(string(plainbody))
		fmt.Println(int64(cep.DateSince1970))
	}

}

func NewDownloadCmd(c *bpopclient.BMClient2, sn []byte, t int64) *bpop.CommandSyn {

	csyn := &bpop.CommandSyn{}

	cdl := &bpop.CmdDownload{}

	csyn.Cmd = cdl

	copy(csyn.SN[:], sn)
	csyn.Sig = ed25519.Sign(c.Priv, sn)

	cdl.MailCnt = 20

	cdl.TimePivot = t
	cdl.Direction = bpop.DirectionToLeft
	cdl.Owner = bmail.ToAddress(c.PK)
	cdl.MailAddr = "testb@eth"

	return csyn
}
