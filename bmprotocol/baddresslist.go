package bmprotocol

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"math/rand"
)

//client hello ===IV,cipher text{self mail address,sn from client}==> server
//server response hello ===IV,sig{sn(from client)},sn(from server) ==>client
//client remote add===>IV,sig{sn(from server)},cipher text{mail list,group list} ==>server
//server response remote add ===> IV,cipher text{sn(old),sn(new from server)},error code===>client
//client pull mail list or group list ===>sig{sn(from server)},IV,cipher text{command(pull mail list or pull all or pull group)} ==server
//server response pull ===>IV,cipher text{mail list,group list,sn(old),sn(new from server)}
//client remote del===>sig{sn},IV,mcipher text{mail list,group list} ==>server
//server response remote del ===> IV,cipher text{sn(old),sn(new from server)},error code===client

type IV [16]byte
type SN [16]byte


type ContactHello struct {
	SelfMA string
	sn SN
}

func NewContactHello() *ContactHello  {
	m:=	&ContactHello{}

	for{
		n,_:=rand.Read(m.sn[:])
		if n != len(m.sn){
			continue
		}
		break
	}
	return m
}

func (mabh *ContactHello)GetSn() SN  {
	return mabh.sn
}

func (mabh *ContactHello)String() string  {
	s:= fmt.Sprintf("%-20s",mabh.SelfMA)
	s+= " sn:"+base58.Encode(mabh.sn[:])

	return s
}

func (mabh *ContactHello)Pack() ([]byte,error)  {
	var(
		r,tmp []byte
		err error
	)

	tmp,err = PackShortString(mabh.SelfMA)
	if err!=nil{
		return nil,err
	}

	r = append(r,tmp...)

	tmp,err = PackShortBytes(mabh.sn[:])
	if err!=nil{
		return nil,err
	}

	r = append(r,tmp...)

	return r,err

}

func (mabh *ContactHello)UnPack(data []byte) (int,error)  {
	var(
		offset,of int
		err error
	)

	mabh.SelfMA,of,err = UnPackShortString(data[offset:])
	if err!=nil{
		return 0,err
	}
	offset += of
	var sn []byte

	sn,	of ,err = UnPackShortBytes(data[offset:])
	if err!=nil{
		return 0,err
	}
	offset += of

	copy(mabh.sn[:],sn)

	return offset,nil

}


//client hello
type CryptContactHello struct {
	translayer.BMTransLayer
	iv IV
	cipherText []byte  //-->MAoBHello
}

func NewCryptContactHello()  *CryptContactHello {
	cch:=&CryptContactHello{}

	bmtl:=translayer.NewBMTL(translayer.CONTACT_HELLO,nil)

	cch.BMTransLayer = *bmtl

	for{
		n,err:=rand.Read(cch.iv[:])
		if err!=nil || n!=len(cch.iv){
			continue
		}
		break
	}

	return cch
}

func (cm *CryptContactHello)GetIV() IV  {
	return cm.iv
}

func (cm *CryptContactHello)SetCipherText(ct []byte)  {
	cm.cipherText = ct
}

func (cm *CryptContactHello)String() string  {
	s:=cm.BMTransLayer.HeadString()
	s+=base58.Encode(cm.iv[:])
	s+="\r\n"
	s+=base58.Encode(cm.cipherText)

	return s
}

func (cm *CryptContactHello)Pack() ([]byte,error){
	var(
		r,tmp []byte
		err error
	)

	tmp,err=PackShortBytes(cm.iv[:])
	if err != nil{
		return nil,err
	}
	r = append(r,tmp...)

	tmp,err = PackLongBytes(cm.cipherText)
	if err!=nil{
		return nil,err
	}
	r = append(r,tmp...)

	cm.BMTransLayer.SetData(r)

	return cm.BMTransLayer.Pack()

}

func (cm *CryptContactHello)UnPack(data []byte) (int,error){
	var(
		offset,of int
		err error
	)

	var iv []byte
	iv,of,err = UnPackShortBytes(data[offset:])
	if err!=nil{
		return 0,nil
	}
	copy(cm.iv[:],iv)

	offset += of

	cm.cipherText,of,err = UnPackLongBytes(data[offset:])

	offset += of
	return offset,nil

}

//server response hello ===IV,sig{sn(from client)},sn(from server) ==>client
type ContactHelloResp struct {
	translayer.BMTransLayer
	iv IV
	clientSN SN
	sigClientSN []byte  //sig from clientSN
	serverSN SN
}

func NewContactHelloResp() *ContactHelloResp {
	chr := &ContactHelloResp{}

	bmtl := translayer.NewBMTL(translayer.CONTACT_HELLO_RESP,nil)
	chr.BMTransLayer = *bmtl

	for{
		n,err:=rand.Read(chr.iv[:])
		if err!=nil || n!=len(chr.iv){
			continue
		}
		break
	}

	for{
		n,err:=rand.Read(chr.serverSN[:])
		if err!=nil || n!=len(chr.iv){
			continue
		}
		break
	}


	return chr
}

func (mr *ContactHelloResp)SetClientSN(sn SN){
	mr.clientSN = sn
}

func (mr *ContactHelloResp)SetSigClientSn(sig []byte)  {
	mr.sigClientSN = sig
}

func (mr *ContactHelloResp)GetIV() IV  {
	return mr.iv
}


func (mr *ContactHelloResp)String() string  {
	s:=mr.BMTransLayer.String()
	s+=fmt.Sprintf("iv:%-30s",base58.Encode(mr.iv[:]))
	s+=fmt.Sprintf("clientSn:%-30s\r\n",base58.Encode(mr.clientSN[:]))
	s+=fmt.Sprintf("SigClientSN:%s\r\n",base58.Encode(mr.sigClientSN))
	s+=fmt.Sprintf("ServerSn:%-30s\r\n",base58.Encode(mr.serverSN[:]))

	return s
}

func (mr *ContactHelloResp)Pack() ([]byte,error){
	var(
		r,tmp []byte
		err error
	)
	tmp,err = PackShortBytes(mr.iv[:])
	if err!=nil{
		return nil,err
	}
	r = append(r,tmp...)

	tmp,err = PackShortBytes(mr.clientSN[:])
	if err!=nil{
		return nil,err
	}
	r = append(r,tmp...)

	tmp,err = PackShortBytes(mr.sigClientSN)
	if err!=nil{
		return nil,err
	}
	r = append(r,tmp...)

	tmp,err = PackShortBytes(mr.serverSN[:])
	if err!=nil{
		return nil,err
	}
	r = append(r,tmp...)

	mr.BMTransLayer.SetData(r)

	return mr.BMTransLayer.Pack()

}

func (mr *ContactHelloResp)UnPack(data []byte) (int,error) {
	var(
		offset,of int
		err error
		buf []byte
	)
	buf,of,err = UnPackShortBytes(data[offset:])
	if err!=nil{
		return 0,err
	}
	offset += of
	copy(mr.iv[:],buf)

	buf,of,err = UnPackShortBytes(data[offset:])
	if err!=nil{
		return 0,err
	}
	offset += of
	copy(mr.clientSN[:],buf)

	mr.sigClientSN,of,err = UnPackShortBytes(data[offset:])
	if err!=nil{
		return 0,err
	}
	offset += of


	buf,of,err = UnPackShortBytes(data[offset:])
	if err!=nil{
		return 0,err
	}
	offset += of
	copy(mr.serverSN[:],buf)

	return offset,nil

}


//client remote add===>IV,sig{sn(from server)},sn(from server),cipher text{mail list,group list} ==>server




type Gid [32]byte

type Cell struct {
	PhoneNum string

}

type BMailAddrss struct {
	MailAddress string
	Alias       string
	Desc        string
	GroupId     Gid
}

type BMailAddressList struct {

}


type GroupDesc struct {
	GroupId Gid
	GroupType int //0 private,1 public in group, 2 public in mail domain, 3 public to all but no mail list
	GroupName string
}

type BGroupAddressList struct {

}

