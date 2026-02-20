package p2p

//this hold the arbituary data  that is sent over each transport
//between the two nodes in the network
type Message struct {}

type RPC struct {
	From string
	Payload []byte
	Stream bool
}