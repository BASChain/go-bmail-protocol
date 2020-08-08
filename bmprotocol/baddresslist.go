package bmprotocol

import (
	"fmt"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"github.com/btcsuite/btcutil/base58"
	"math/rand"

	"encoding/binary"
	"errors"
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
	sn     SN
}

func NewContactHello() *ContactHello {
	m := &ContactHello{}

	for {
		n, _ := rand.Read(m.sn[:])
		if n != len(m.sn) {
			continue
		}
		break
	}
	return m
}

func (mabh *ContactHello) GetSn() SN {
	return mabh.sn
}

func (mabh *ContactHello) String() string {
	s := fmt.Sprintf("%-20s", mabh.SelfMA)
	s += " sn:" + base58.Encode(mabh.sn[:])

	return s
}

func (mabh *ContactHello) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortString(mabh.SelfMA)
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	tmp, err = PackShortBytes(mabh.sn[:])
	if err != nil {
		return nil, err
	}

	r = append(r, tmp...)

	return r, err

}

func (mabh *ContactHello) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)

	mabh.SelfMA, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	var sn []byte

	sn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	copy(mabh.sn[:], sn)

	return offset, nil

}

//client hello
type CryptContactHello struct {
	translayer.BMTransLayer
	iv         IV
	cipherText []byte //-->MAoBHello
}

func NewCryptContactHello() *CryptContactHello {
	cch := &CryptContactHello{}

	bmtl := translayer.NewBMTL(translayer.CONTACT_HELLO)

	cch.BMTransLayer = *bmtl

	for {
		n, err := rand.Read(cch.iv[:])
		if err != nil || n != len(cch.iv) {
			continue
		}
		break
	}

	return cch
}

func (cm *CryptContactHello) GetIV() IV {
	return cm.iv
}

func (cm *CryptContactHello) SetCipherText(ct []byte) {
	cm.cipherText = ct
}

func (cm *CryptContactHello) String() string {
	s := cm.BMTransLayer.String()
	s += base58.Encode(cm.iv[:])
	s += "\r\n"
	s += base58.Encode(cm.cipherText)

	return s
}

func (cm *CryptContactHello) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortBytes(cm.iv[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackLongBytes(cm.cipherText)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	//cm.BMTransLayer.SetData(r)

	return cm.BMTransLayer.Pack()

}

func (cm *CryptContactHello) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)

	var iv []byte
	iv, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	copy(cm.iv[:], iv)

	offset += of

	cm.cipherText, of, err = UnPackLongBytes(data[offset:])

	offset += of
	return offset, nil

}

//server response hello ===IV,sig{sn(from client)},sn(from server) ==>client
type ContactHelloResp struct {
	translayer.BMTransLayer
	iv          IV
	clientSN    SN
	sigClientSN []byte //sig from clientSN
	serverSN    SN
}

func NewContactHelloResp() *ContactHelloResp {
	chr := &ContactHelloResp{}

	bmtl := translayer.NewBMTL(translayer.CONTACT_HELLO_RESP)
	chr.BMTransLayer = *bmtl

	for {
		n, err := rand.Read(chr.iv[:])
		if err != nil || n != len(chr.iv) {
			continue
		}
		break
	}

	for {
		n, err := rand.Read(chr.serverSN[:])
		if err != nil || n != len(chr.iv) {
			continue
		}
		break
	}

	return chr
}

func (mr *ContactHelloResp) SetClientSN(sn SN) {
	mr.clientSN = sn
}

func (mr *ContactHelloResp) SetSigClientSn(sig []byte) {
	mr.sigClientSN = sig
}

func (mr *ContactHelloResp) GetIV() IV {
	return mr.iv
}

func (mr *ContactHelloResp) String() string {
	s := mr.BMTransLayer.String()
	s += fmt.Sprintf("iv:%-30s", base58.Encode(mr.iv[:]))
	s += fmt.Sprintf("clientSn:%-30s\r\n", base58.Encode(mr.clientSN[:]))
	s += fmt.Sprintf("SigClientSN:%s\r\n", base58.Encode(mr.sigClientSN))
	s += fmt.Sprintf("ServerSn:%-30s\r\n", base58.Encode(mr.serverSN[:]))

	return s
}

func (mr *ContactHelloResp) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)
	tmp, err = PackShortBytes(mr.iv[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(mr.clientSN[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(mr.sigClientSN)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(mr.serverSN[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	//mr.BMTransLayer.SetData(r)

	return mr.BMTransLayer.Pack()

}

func (mr *ContactHelloResp) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		buf        []byte
	)
	buf, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	copy(mr.iv[:], buf)

	buf, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	copy(mr.clientSN[:], buf)

	mr.sigClientSN, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	buf, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	copy(mr.serverSN[:], buf)

	return offset, nil

}

type Gid [32]byte

type Cell struct {
	PhoneNum  string
	PhoneType string
}

func (c *Cell) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortString(c.PhoneNum)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortString(c.PhoneType)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return r, nil
}

func (c *Cell) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
	)

	c.PhoneNum, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	c.PhoneType, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}

	offset += of

	return offset, nil

}

type BMailAddrss struct {
	MailAddress string
	Alias       string
	Desc        string
	Phone       Cell
	GroupId     Gid
}

func (bma *BMailAddrss) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortString(bma.MailAddress)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortString(bma.Alias)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortString(bma.Desc)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	c := &bma.Phone
	tmp, err = c.Pack()
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(bma.GroupId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return r, nil
}

func (bma *BMailAddrss) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	bma.MailAddress, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	bma.Alias, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	bma.Desc, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	c := &bma.Phone
	of, err = c.UnPack(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	copy(bma.GroupId[:], tmp)

	return offset, nil

}

type GroupDesc struct {
	GroupId   Gid
	GroupType int //0 private,1 public in group, 2 public in mail domain, 3 public to all but no mail list
	GroupName string
}

func (gd *GroupDesc) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortBytes(gd.GroupId[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(gd.GroupType))
	r = append(r, tmp...)

	tmp, err = PackShortString(gd.GroupName)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	return r, nil

}

func (gd *GroupDesc) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}

	copy(gd.GroupId[:], tmp)
	offset += of

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack group type error")
	}
	gd.GroupType = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	gd.GroupName, of, err = UnPackShortString(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	return offset, nil

}

//client remote add===>IV,sig{sn(from server)},sn(from server),cipher text{mail list,group list} ==>server
type ContactAdd struct {
	iv          IV
	sigServerSn []byte
	serverSN    SN
	mailAddrs   []BMailAddrss
	groups      []GroupDesc
}

func (ca *ContactAdd) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)
	tmp, err = PackShortBytes(ca.iv[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(ca.sigServerSn[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)
	tmp, err = PackShortBytes(ca.serverSN[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(len(ca.mailAddrs)))
	r = append(r, tmp...)
	for i := 0; i < len(ca.mailAddrs); i++ {
		ma := &ca.mailAddrs[i]

		tmp, err = ma.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)

	}

	tmp = translayer.UInt32ToBuf(uint32(len(ca.groups)))
	r = append(r, tmp...)
	for i := 0; i < len(ca.groups); i++ {
		g := &ca.groups[i]

		tmp, err = g.Pack()
		if err != nil {
			return nil, err
		}
		r = append(r, tmp...)
	}
	return r, nil
}

func (ca *ContactAdd) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)
	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	copy(ca.iv[:], tmp)

	ca.sigServerSn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, err
	}
	offset += of
	copy(ca.serverSN[:], tmp)

	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack mail address error")
	}
	l := int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	for i := 0; i < l; i++ {
		ma := &BMailAddrss{}
		of, err = ma.UnPack(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += of
		ca.mailAddrs = append(ca.mailAddrs, *ma)
	}
	if len(data) < offset+translayer.Uint32Size {
		return 0, errors.New("unpack groups error")
	}
	l = int(binary.BigEndian.Uint32(data[offset:]))
	offset += translayer.Uint32Size

	for i := 0; i < l; i++ {
		g := &GroupDesc{}
		of, err = g.UnPack(data[offset:])
		if err != nil {
			return 0, err
		}
		offset += of
		ca.groups = append(ca.groups, *g)
	}

	return offset, nil
}

type CryptContactAdd struct {
	translayer.BMTransLayer
	iv          IV
	sigServerSn []byte
	serverSN    SN
	cipherTxt   []byte
}

func (cca *CryptContactAdd) SetIV(iv IV) {
	cca.iv = iv
}

func (cca *CryptContactAdd) SetSigServerSn(sig []byte) {
	cca.sigServerSn = sig
}
func (cca *CryptContactAdd) SetServerSn(sn SN) {
	cca.serverSN = sn
}

func (cca *CryptContactAdd) SetCipherTxt(ct []byte) {
	cca.cipherTxt = ct
}

func NewCryptContactAdd() *CryptContactAdd {
	bmtl := translayer.NewBMTL(translayer.CONTACT_ADD)
	cca := &CryptContactAdd{}
	cca.BMTransLayer = *bmtl

	//for{
	//	n,err:=rand.Read(cca.iv[:])
	//	if err!=nil || n!=len(cca.iv){
	//		continue
	//	}
	//	break
	//}

	return cca
}

func (cca *CryptContactAdd) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)
	tmp, err = PackShortBytes(cca.iv[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)
	tmp, err = PackShortBytes(cca.sigServerSn)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)
	tmp, err = PackShortBytes(cca.serverSN[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)
	tmp, err = PackLongBytes(cca.cipherTxt)
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	//cca.BMTransLayer.SetData(r)

	return cca.BMTransLayer.Pack()

}

func (cca *CryptContactAdd) UnPack(data []byte) (int, error) {
	var (
		offset, of int
		err        error
		tmp        []byte
	)

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of
	copy(cca.iv[:], tmp)

	cca.sigServerSn, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	tmp, of, err = UnPackShortBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of
	copy(cca.serverSN[:], tmp)

	cca.cipherTxt, of, err = UnPackLongBytes(data[offset:])
	if err != nil {
		return 0, nil
	}
	offset += of

	return offset, nil
}

//server response remote add ===> IV,cipher text{sn(old),sn(new from server)},error code===>client
type ContactAddResp struct {
	iv      IV
	sn      SN
	snNew   SN
	errCode int
}

func (car *ContactAddResp) Pack() ([]byte, error) {
	var (
		r, tmp []byte
		err    error
	)

	tmp, err = PackShortBytes(car.iv[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(car.sn[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp, err = PackShortBytes(car.snNew[:])
	if err != nil {
		return nil, err
	}
	r = append(r, tmp...)

	tmp = translayer.UInt32ToBuf(uint32(car.errCode))
	r = append(r, tmp...)

	return r, nil

}

type CryptContactAddResp struct {
	translayer.BMTransLayer
	iv        IV
	cipherTxt []byte
	errCode   int
}

const (
	PullMail int = iota + 1
	PullGroup
	s
)

//client pull mail list or group list ===>sig{sn(from server)},IV,cipher text{command(pull mail list or pull all or pull group)} ==server
type ContactPull struct {
	sn          SN
	sigSN       []byte
	iv          IV
	pullCommand int
}

type CryptContactPull struct {
	translayer.BMTransLayer
	iv          IV
	cipherTxt   []byte
	pullCommand int //PullMail,PullGroup,PullAll
}

//server response pull ===>IV,cipher text{mail list,group list,sn(old),sn(new from server)}
type ContactPullResp struct {
	iv     IV
	sn     SN
	snNew  SN
	mails  []BMailAddrss
	groups []GroupDesc
}

type CryptContactPullResp struct {
	iv        IV
	cipherTxt []byte
}

//client remote del===>sig{sn},IV,cipher text{mail list,group list} ==>server
type ContactDel struct {
	iv     IV
	sn     SN
	sigSn  []byte
	mails  []BMailAddrss
	groups []GroupDesc
}

type CryptContactDel struct {
	iv        IV
	sn        SN
	sigSn     []byte
	cipherTxt []byte
}

//server response remote del ===> IV,cipher text{sn(old),sn(new from server)},error code===client
type ContactDelResp struct {
	iv      IV
	sn      SN
	snNew   SN
	errCode int
}

type CryptContactDelResp struct {
	iv        IV
	cipherTxt []byte
	errCode   int
}
