package main

import (
	"crypto/rand"
	"crypto/ed25519"
	"github.com/BASChain/go-account"
	"fmt"
	ed255192 "github.com/kprc/chatserver/ed25519"
	"github.com/btcsuite/btcutil/base58"
)

func DeriveKey(seed []byte) (pub ed25519.PublicKey,priv ed25519.PrivateKey)  {
	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privateKey[32:])

	return publicKey,privateKey
}

func main()  {

	pub1, pri1, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}


	pub2, pri2, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	pub3, pri3, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	aesgrp1,grpks,err:= ed255192.GenGroupAesKey(pri1,[][]byte{pub1,pub2,pub3})
	if err!=nil{
		fmt.Println(err)
		return
	}

	aesgrp2,err:=ed255192.DeriveGroupKey(pri2,grpks,[][]byte{pub1,pub2,pub3})
	if err!=nil{
		fmt.Println(err)
		return
	}

	aesgrp3,err:=ed255192.DeriveGroupKey(pri3,grpks,[][]byte{pub1,pub2,pub3})
	if err!=nil{
		fmt.Println(err)
		return
	}

	aesk12,err := account.GenerateAesKey(pub2,pri1)
	if err!=nil{
		fmt.Println(err)
		return
	}

	pub12,priv12:=DeriveKey(aesk12)


	aes123,err:=account.GenerateAesKey(pub3,priv12)
	if err!=nil{
		fmt.Println(err)
		return
	}

	ciphertxt,err:=account.Encrypt(aes123,[]byte("hello world"))
	if err!=nil{
		fmt.Println(err)
		return
	}

	ciphertxt1,err:=account.Encrypt(aes123,[]byte("hello world"))
	if err!=nil{
		fmt.Println(err)
		return
	}


	aes21,err:=account.GenerateAesKey(pub1,pri2)
	if err!=nil{
		fmt.Println(err)
		return
	}

	_,priv21:=DeriveKey(aes21)

	aes213,err:=account.GenerateAesKey(pub3,priv21)
	if err!=nil{
		fmt.Println(err)
		return
	}

	ciphertxt2,err:=account.Encrypt(aes213,[]byte("hello world"))
	if err!=nil{
		fmt.Println(err)
		return
	}

	plaintxt,err:=account.Decrypt(aes213,ciphertxt)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println(string(plaintxt))

	aes312,err:=account.GenerateAesKey(pub12,pri3)
	if err!=nil{
		fmt.Println(err)
		return
	}
	plaintxt1,err:=account.Decrypt(aes312,ciphertxt1)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println(string(plaintxt1))


	plaintxt2,err:=account.Decrypt(aes312,ciphertxt2)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Println(string(plaintxt2))

	fmt.Println(base58.Encode(aesgrp1))
	fmt.Println(base58.Encode(aesgrp2))
	fmt.Println(base58.Encode(aesgrp3))

}



