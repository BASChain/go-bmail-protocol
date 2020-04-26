package bmprotocol

import "github.com/BASChain/go-bmail-protocol/translayer"


//client -> server
type SendEnvelope struct {
	translayer.BMTransLayer
	From string
	RecpAddr string     //recipient
	LAddr []byte		//local public key
	//===cipher text begin====
	To []string
	CC []string
	BC []string
	Subject string
	Data string
	//====cipher text end ===
	IV []byte			//sn from Bhello
	Sig []byte          //signature
}


//server -> client
type RespSendEnvelope struct {
	translayer.BMTransLayer
	IV []byte			//same as SendEnvelope
}


//server -> client
type SendEnvelopeFail struct {

}


//client -> server
type RespSendEnvelopeFail struct {
	
}