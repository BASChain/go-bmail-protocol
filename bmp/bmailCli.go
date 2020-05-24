package bmp

import (
	"fmt"
	"github.com/BASChain/go-bmail-account"
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

func (bmc *BMailClient) SendP2sMail(re *RawEnvelope) error {
	return nil
}

func (bmc *BMailClient) SendP2pMail(re *RawEnvelope) error {

	conn, err := NewBMConn(bmc.SrvIP)

	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Helo()
	if err != nil {
		return err
	}

	ack := &HELOACK{}
	if err := conn.ReadWithHeader(ack); err != nil {
		return err
	}
	fmt.Println("===hel ack===>", ack)
	if bmc.SrvBcas[ack.SrvBca] == false {
		return fmt.Errorf("invalid bmail server block chain address:[%s]", ack.SrvBca)
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
	msg := &EnvelopeSyn{
		Mode: BMailModeP2P,
		SN:   ack.SN,
		Sig:  bmc.Wallet.Sign(ack.SN.Bytes()),
		Hash: synHash,
		Env:  cryptEnv,
	}
	if err := conn.SendWithHeader(msg); err != nil {
		return err
	}

	msgAck := &EnvelopeAck{}
	if err := conn.ReadWithHeader(msgAck); err != nil {
		return err
	}
	fmt.Println("===envelop ack===>", msgAck, hexutil.Encode(synHash), hexutil.Encode(msgAck.Hash))
	if !bmail.Verify(ack.SrvBca, synHash, msgAck.Sig) {
		return fmt.Errorf("verify header ack failed:[%s]", ack.SrvBca)
	}

	return nil
}
