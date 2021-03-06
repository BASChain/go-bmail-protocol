package client

import (
	"fmt"
	"github.com/realbmail/go-bmail-account"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-bmail-protocol/bpop"
	resolver "github.com/realbmail/go-bmail-resolver"
	"net"
	"strings"
)

type ClientConf struct {
	Resolver resolver.NameResolver
	Wallet   bmail.Wallet
}

type BMailClient struct {
	Wallet   bmail.Wallet
	SrvIP    net.IP
	SrvBcas  map[bmail.Address]bool
	resolver resolver.NameResolver
}

func NewClient(cc *ClientConf) (*BMailClient, error) {

	r := cc.Resolver
	mailName := cc.Wallet.MailAddress()
	if len(mailName) == 0 {
		return nil, fmt.Errorf("invalid sender account")
	}

	mailParts := strings.Split(mailName, "@")
	if len(mailParts) != 2 {
		return nil, fmt.Errorf("invalid mail name")
	}

	ips, bcas := r.DomainMX(mailParts[1])
	if len(ips) == 0 || len(bcas) == 0 {
		return nil, fmt.Errorf("no valid mx[%s] record", mailParts[1])
	}
	fmt.Println("====2==> mx:", ips[0], bcas[0])
	srvIP := choseBestServer(ips)

	obj := &BMailClient{
		Wallet:   cc.Wallet,
		SrvIP:    srvIP,
		SrvBcas:  make(map[bmail.Address]bool),
		resolver: r,
	}
	for _, bca := range bcas {
		obj.SrvBcas[bca] = true
		fmt.Println("===mx bca===>", bca)
	}
	return obj, nil
}

func choseBestServer(ips []net.IP) net.IP {
	fmt.Println("======>selected server ip:", ips[0].String())
	return ips[0]
}

func (bmc *BMailClient) Close() {
	bmc.Wallet = nil
	bmc.SrvIP = nil
	bmc.SrvBcas = nil
	bmc.resolver = nil
}

func (bmc *BMailClient) SendMail(bme *bmp.BMailEnvelope) error {

	conn, err := bmp.NewBMConn(bmc.SrvIP)
	if err != nil {
		return err
	}
	defer conn.Close()
	ack, err := bmc.HandShake(conn)
	if err != nil {
		return err
	}
	synHash := bme.Hash()
	signature := bmc.Wallet.Sign(ack.SN.Bytes())

	msg := &bmp.EnvelopeSyn{
		SN:   ack.SN,
		Sig:  signature,
		Hash: synHash,
		Env:  bme,
	}
	if err := conn.SendWithHeader(msg); err != nil {
		return err
	}

	msgAck := &bmp.EnvelopeAck{}
	if err := conn.ReadWithHeader(msgAck); err != nil {
		return err
	}
	if !bmail.Verify(ack.SrvBca, synHash, msgAck.Sig) {
		return fmt.Errorf("verify header ack failed:[%s]", ack.SrvBca)
	}

	return nil
}
func (bmc *BMailClient) HandShake(conn *bmp.BMailConn) (*bmp.HELOACK, error) {

	if err := conn.Helo(); err != nil {
		return nil, err
	}

	ack := &bmp.HELOACK{}
	if err := conn.ReadWithHeader(ack); err != nil {
		return nil, err
	}
	if bmc.SrvBcas[ack.SrvBca] == false {
		return nil, fmt.Errorf("invalid bmail server block chain address:[%s]", ack.SrvBca)
	}
	return ack, nil
}

func (bmc *BMailClient) ReceiveEnv(timeSince1970 int64, olderThanSince bool, maxCount int) ([]*bmp.BMailEnvelope, error) {
	conn, err := bmp.NewBMConn(bmc.SrvIP)
	if err != nil {
		fmt.Println("NewBMConn------>", err)
		return nil, err
	}
	defer conn.Close()

	ack, err := bmc.HandShake(conn)
	if err != nil {
		fmt.Println("HandShake------>", err)
		return nil, err
	}
	fmt.Println("HandShake------success>")
	sig := bmc.Wallet.Sign(ack.SN[:])
	cmd := &bpop.CommandSyn{
		Sig: sig,
		SN:  ack.SN,
		Cmd: &bpop.CmdDownload{
			MailCnt:   maxCount,
			TimePivot: timeSince1970,
			Direction: olderThanSince,
			Owner:     bmc.Wallet.Address(),
			MailAddr:  bmc.Wallet.MailAddress(),
		},
	}

	if err := conn.SendWithHeader(cmd); err != nil {
		fmt.Println("SendWithHeader------>", err)
		return nil, err
	}

	fmt.Println("======>:SendWithHeader success>", timeSince1970)
	cmdAck := &bpop.CommandAck{}
	cmdAck.CmdCxt = &bpop.CmdDownloadAck{}
	if err := conn.ReadWithHeader(cmdAck); err != nil {
		fmt.Println("ReadWithHeader------>", err)
		return nil, err
	}
	//hash := resp.CmdCxt.Hash()
	//if bytes.Compare(hash[:], resp.Hash) != 0 {
	//	fmt.Println("hash error")
	//	return
	//}

	if cmdAck.ErrorCode != 0 {
		if cmdAck.ErrorCode == 1 {
			return make([]*bmp.BMailEnvelope, 0), nil
		}
		return nil, fmt.Errorf("fetch data failed, server error:%d", ack.ErrCode)
	}

	//if !bmail.Verify(ack.SrvBca, cmdAck.Hash, cmdAck.Sig) {
	//	return nil, fmt.Errorf("verify header ack failed:[%s]", ack.SrvBca)
	//}

	envs := cmdAck.CmdCxt.(*bpop.CmdDownloadAck)
	fmt.Println("======>:envelope loaded success=>", len(envs.CryptEps))
	return envs.CryptEps, nil
}
