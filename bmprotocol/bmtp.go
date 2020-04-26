package bmprotocol

import "github.com/BASChain/go-bmail-protocol/translayer"


const(
	PeerUnreachable int = iota + 1
	AddressUnavailable
)

const(
	ErrMsg_PeerUnreachable string = "Peer is unreachable"
	ErrMsg_AddressUnavailable string = "Recipient is not available"
)


type EnvelopeHead struct {
	From string
	RecpAddr string     //recipient
	LAddr []byte		//local public key
}




type EnvelopeContent struct {
	To []string
	CC []string
	BC []string
	Subject string
	Data string
}


type EnvelopeEnd struct {
	IV []byte			//sn from Bhello
	Sig []byte          //signature
}

//client -> server
type SendEnvelope struct {
	translayer.BMTransLayer
	EnvelopeHead
	CipherTxt []byte    //crypt from EnvelopeContent
	EnvelopeEnd
}


//server -> client
type RespSendEnvelope struct {
	translayer.BMTransLayer
	From string
	RecpAddr string     //recipient
	LAddr []byte		//local public key
	IV []byte			//same as SendEnvelope
}


//server -> client
type SendEnvelopeFail struct {
	translayer.BMTransLayer
	EnvelopeHead
	CipherTxt []byte    //crypt from EnvelopeContent
	EnvelopeEnd
	ErrorCode int
}


//client -> server
type RespSendEnvelopeFail struct {
	translayer.BMTransLayer
	EnvelopeHead
	IV []byte
}



