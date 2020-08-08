package bpopclient

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BASChain/go-account"
	"github.com/realbmail/go-bas-mail-server/bmailcrypt"
	"github.com/realbmail/go-bmail-account"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-bmail-protocol/bpop"
	"github.com/realbmail/go-bmail-protocol/translayer"
	"github.com/howeyc/gopass"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type BMClient2 struct {
	sn      []byte
	c       *net.TCPConn
	timeout int //second
	PK      ed25519.PublicKey
	Priv    ed25519.PrivateKey
	Hash    []byte
	SrvPk   ed25519.PublicKey
}

func (c *BMClient2) GetSn() []byte {
	return c.sn
}

const (
	keyfile string = "ed25519key.testb"
)

func GetKeyfile() string {
	fn, _ := filepath.Abs(os.Args[0])
	name := filepath.Base(fn)

	dir := string(fn[:len(fn)-len(name)])

	return path.Join(dir, keyfile)
}

func keyISGenerated() bool {
	if tools.FileExists(GetKeyfile()) {
		return true
	}

	return false
}

func LoadKey(password string) (ed25519.PublicKey, ed25519.PrivateKey) {
	data, err := tools.OpenAndReadAll(GetKeyfile())
	if err != nil {
		log.Fatal("Load From key file error")
		return nil, nil
	}

	kj := &bmailcrypt.KeyJson{}

	err = json.Unmarshal(data, kj)
	if err != nil {
		log.Fatal("Load From json error")
		return nil, nil
	}

	pk := bmail.Address(kj.PubKey).ToPubKey()
	var priv ed25519.PrivateKey
	priv, err = account.DecryptSubPriKey(pk, kj.CipherKey, password)
	if err != nil {
		log.Fatal("Decrypt PrivKey failed")
		return nil, nil
	}

	return pk, priv
}

func GenEd25519KeyAndSave(password string) error {

	var (
		priv ed25519.PrivateKey
		pub  ed25519.PublicKey
		err  error
	)
	cnt := 0
	for {
		cnt++
		pub, priv, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			if cnt > 10 {
				return err
			}
			continue
		} else {
			break
		}
	}

	var cipherTxt string
	cipherTxt, err = account.EncryptSubPriKey(priv, pub, password)
	if err != nil {
		return err
	}

	kj := &bmailcrypt.KeyJson{PubKey: bmail.ToAddress(pub[:]).String(), CipherKey: cipherTxt}

	fmt.Println("PubKey:", kj.PubKey)

	var data []byte
	data, err = json.Marshal(*kj)
	fmt.Println("save to: ", GetKeyfile())
	err = tools.Save2File(data, GetKeyfile())
	if err != nil {
		return err
	}

	return nil
}

func inputpassword() (password string, err error) {
	passwd, err := gopass.GetPasswdPrompt("Please Enter Password: ", true, os.Stdin, os.Stdout)
	if err != nil {
		return "", err
	}

	if len(passwd) < 1 {
		return "", errors.New("Please input valid password")
	}

	return string(passwd), nil
}

func NewClient2(serverIP net.IP, timeout int) *BMClient2 {
	laddr := &net.TCPAddr{}
	raddr := &net.TCPAddr{IP: serverIP, Port: 1110}
	conn, err := net.DialTCP("tcp4", laddr, raddr)
	if err != nil {
		return nil
	}

	var passwd string

	passwd, _ = inputpassword()

	if !keyISGenerated() {

		GenEd25519KeyAndSave(passwd)
	}

	pk, priv := LoadKey(passwd)

	c := &BMClient2{}
	c.c = conn
	c.timeout = timeout
	c.PK = pk
	c.Priv = priv

	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))

	return c
}

func (c *BMClient2) Close() {
	c.c.Close()
	c.c = nil
	c.sn = nil
}

func (c *BMClient2) Helo() (err error) {

	if c.c == nil {
		return errors.New("client is not initialized")
	}

	header := bmp.Header{
		Ver:    translayer.BMAILVER1,
		MsgTyp: translayer.HELLO,
		MsgLen: 0,
	}

	data := header.GetBytes()
	if _, err := c.c.Write(data); err != nil {
		return err
	}

	buf := make([]byte, translayer.BMHeadSize())
	var n int

	n, err = c.c.Read(buf)
	if n != translayer.BMHeadSize() || err != nil {
		fmt.Println(err, n)
		return errors.New("helo Read a bad bmail head")
	}

	bmtl := &translayer.BMTransLayer{}
	bmtl.UnPack(buf)

	if bmtl.GetMsgType() != translayer.HELLO_ACK || bmtl.GetDataLen() == 0 {
		return errors.New("Received a error message: " + strconv.Itoa(int(bmtl.GetMsgType())))
	}

	//read left
	buf = make([]byte, bmtl.GetDataLen())

	n, err = c.c.Read(buf)
	if n != int(bmtl.GetDataLen()) || err != nil {
		return errors.New("helo Read a bad bmail data")
	}

	ha := &bmp.HELOACK{}

	err = json.Unmarshal(buf, ha)
	if err != nil {
		return err
	}

	c.SrvPk = bmail.Address(ha.SrvBca).ToPubKey()

	c.sn = ha.SN[:]

	return nil
}

func (c *BMClient2) SendCommand(cmd *bpop.CommandSyn) (ca *bpop.CommandAck, err error) {
	if c.c == nil {
		return nil, errors.New("client is not initialized")
	}
	if bytes.Compare(cmd.SN[:], c.sn) != 0 {
		return nil, errors.New("Envelope not correct")
	}

	data, err := cmd.GetBytes()
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))

	header := bmp.Header{
		Ver:    translayer.BMAILVER1,
		MsgTyp: cmd.MsgType(),
		MsgLen: len(data),
	}

	if _, err := c.c.Write(header.GetBytes()); err != nil {
		return nil, err
	}

	var n int
	n, err = c.c.Write(data)
	if n != len(data) || err != nil {
		return nil, errors.New("Send envelope Failed")
	}
	buf := make([]byte, translayer.BMHeadSize())

	n, err = c.c.Read(buf)
	if n != translayer.BMHeadSize() || err != nil {
		return nil, errors.New("Read a bad bmail head")
	}

	bmtl := &translayer.BMTransLayer{}
	bmtl.UnPack(buf)

	if bmtl.GetMsgType() != translayer.RETR_RESP || bmtl.GetDataLen() == 0 {
		return nil, errors.New("Received a error message: " + strconv.Itoa(int(bmtl.GetMsgType())))
	}

	buf = make([]byte, bmtl.GetDataLen())

	//
	//n, err = c.c.Read(buf)
	//if n != int(bmtl.GetDataLen()) || err != nil {
	//	fmt.Println(err)
	//	return nil, errors.New("Read a bad bmail data")
	//}

	total := 0

	for {
		n, err := c.c.Read(buf[total:])
		if err != nil && total < int(bmtl.GetDataLen()) {
			return nil, err
		}

		total += n
		if total >= int(bmtl.GetDataLen()) {
			break
		}
	}

	//for {
	//	n, err := c.c.Read(buf[total:])
	//	if err != nil {
	//		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
	//			total += n
	//			continue
	//		} else if err != io.EOF {
	//			return nil, err
	//		}
	//		total += n
	//	} else {
	//		total += n
	//	}
	//	if n == 0 && err == io.EOF {
	//		return nil, errors.New("no data to read")
	//	}
	//	if total >= int(bmtl.GetDataLen()) {
	//		break
	//	}
	//}

	resp := &bpop.CommandAck{}
	resp.CmdCxt = &bpop.CmdDownloadAck{}

	err = json.Unmarshal(buf, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
