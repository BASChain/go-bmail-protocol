package test

import (
	"github.com/BASChain/go-bmail-protocol/bmprotocol"
	"testing"
)

func Test_BHello(t *testing.T) {
	bmh := bmprotocol.NewBMHello()

	data, _ := bmh.Pack()

}
