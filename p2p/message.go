package p2p

//this hold the data
type RPC struct {
	From string
	Payload []byte
	Stream bool
}