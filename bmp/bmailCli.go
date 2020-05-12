package bmp

import (
	"fmt"
	"github.com/BASChain/go-bmail-account"
	resolver "github.com/BASChain/go-bmail-resolver"
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
	SrvBcas  []bmail.Address
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
		return nil, fmt.Errorf("no valid mx record")
	}
	srvIP := choseBestServer(ips)

	return &BMailClient{
		Wallet:   cc.Wallet,
		SrvIP:    srvIP,
		SrvBcas:  bcas,
		resolver: r,
	}, nil
}

func choseBestServer(ips []net.IP) net.IP {
	return ips[0]
}

func (bmc *BMailClient) SendMail(env *Envelope) error {
	return nil
}
