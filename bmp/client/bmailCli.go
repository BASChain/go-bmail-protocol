package client

import (
	"fmt"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-protocol/bpop"
	resolver "github.com/BASChain/go-bmail-resolver"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	fmt.Println("======> mx:", len(ips), len(bcas))
	if len(ips) == 0 || len(bcas) == 0 {
		return nil, fmt.Errorf("no valid mx record")
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
		fmt.Println("===got mx bca===>", bca)
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

func (bmc *BMailClient) SendP2sMail(re *bmp.RawEnvelope) error {
	return nil
}

func (bmc *BMailClient) SendP2pMail(re *bmp.RawEnvelope) error {

	conn, err := bmp.NewBMConn(bmc.SrvIP)
	if err != nil {
		return err
	}
	defer conn.Close()
	ack, err := bmc.HandShake(conn)
	if err != nil {
		return err
	}
	aesKey, err := bmc.Wallet.AeskeyOf(re.ToAddr.ToPubKey())
	if err != nil {
		return err
	}
	cryptEnv, err := re.Seal(aesKey)
	if err != nil {
		return err
	}
	synHash := cryptEnv.Hash()
	msg := &bmp.EnvelopeSyn{
		Mode: bmp.BMailModeP2P,
		SN:   ack.SN,
		Sig:  bmc.Wallet.Sign(ack.SN.Bytes()),
		Hash: synHash,
		Env:  cryptEnv,
	}
	if err := conn.SendWithHeader(msg); err != nil {
		return err
	}

	msgAck := &bmp.EnvelopeAck{}
	if err := conn.ReadWithHeader(msgAck); err != nil {
		return err
	}
	fmt.Println("===envelop ack===>", msgAck, hexutil.Encode(synHash), hexutil.Encode(msgAck.Hash))
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
	fmt.Println("===hel ack===>", ack)
	if bmc.SrvBcas[ack.SrvBca] == false {
		return nil, fmt.Errorf("invalid bmail server block chain address:[%s]", ack.SrvBca)
	}

	fmt.Println("get hello ack:", ack)

	return ack, nil
}

func (bmc *BMailClient) ReceiveEnv(timeSince1970 int64) ([]bmp.CryptEnvelope, error) {
	conn, err := bmp.NewBMConn(bmc.SrvIP)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ack, err := bmc.HandShake(conn)
	if err != nil {
		return nil, err
	}
	sig := bmc.Wallet.Sign(ack.SN[:])
	cmd := &bpop.CommandSyn{
		Sig: sig,
		SN:  ack.SN,
		Cmd: &bpop.CmdDownload{
			MailCnt:    20,
			BeforeTime: timeSince1970,
			Owner:      bmc.Wallet.Address(),
			MailAddr:   bmc.Wallet.MailAddress(),
		},
	}

	if err := conn.SendWithHeader(cmd); err != nil {
		return nil, err
	}

	cmdAck := &bpop.CommandAck{}
	cmdAck.CmdCxt = &bpop.CmdDownloadAck{}
	if err := conn.ReadWithHeader(cmdAck); err != nil {
		return nil, err
	}
	//hash := resp.CmdCxt.Hash()
	//if bytes.Compare(hash[:], resp.Hash) != 0 {
	//	fmt.Println("hash error")
	//	return
	//}
	//
	//if !bmailcrypt.Verify(c.SrvPk, hash, resp.Sig) {
	//	fmt.Println("not a correct server")
	//} else {
	//	fmt.Println("you bmail have send to a correct server")
	//}

	fmt.Println("======> bpop ack data=>:", cmdAck, cmdAck.ErrorCode)
	envs := cmdAck.CmdCxt.(*bpop.CmdDownloadAck)
	fmt.Println("======> CmdDownloadAck:", envs)
	return envs.CryptEps, nil
}
