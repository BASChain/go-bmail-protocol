package bmp

import (
	"fmt"
	"github.com/BASChain/go-bmail-account"
	resolver "github.com/BASChain/go-bmail-resolver"
	"net"
	"strings"
)

type ClientConf struct {
	Debug  bool
	Wallet bmail.Wallet
}

type BMailClient struct {
	Wallet   bmail.Wallet
	conn     *BMailConn
	resolver resolver.NameResolver
}

func NewClient(cc *ClientConf) (*BMailClient, error) {

	r := resolver.NewEthResolver(cc.Debug)
	mailName := cc.Wallet.MailAddress()
	if len(mailName) == 0 {
		return nil, fmt.Errorf("invalid mail account")
	}
	mailParts := strings.Split(mailName, "@")
	if len(mailParts) != 2 {
		return nil, fmt.Errorf("invalid mail name")
	}

	ips := r.DomainMX(mailParts[1])
	if len(ips) == 0 {
		return nil, fmt.Errorf("no valid mx record")
	}

	srvIP := choseBestServer(ips)
	conn := NewBMConn(srvIP)
	if conn == nil {
		return nil, fmt.Errorf("connect to bmail server failed")
	}
	return &BMailClient{
		Wallet:   cc.Wallet,
		conn:     conn,
		resolver: r,
	}, nil
}

func choseBestServer(ips []net.IP) net.IP {
	return ips[0]
}

func (bc *BMailClient) Prepare() Envelop {
	return nil
}

func (bc *BMailClient) Close() {
	bc.conn.Close()
}
