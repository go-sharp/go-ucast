package ucast

type msgType int8

const (
	msgHello msgType = iota + 1
	msgData
	msgFecData
)

type ucastMessage struct {
	msgID uint32
}
