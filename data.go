package ucast

import (
	"bytes"
	"encoding/binary"
	"errors"
	// Just to ensure package is offline available
	_ "github.com/vivint/infectious"
)

const (
	isFecFlag = 1 << 7
)

type serializer interface {
	toNetByteOrder() ([]byte, error)
}

type deserializer interface {
	fromNetByteOrder(data []byte) error
}

type msgType int8

const (
	msgHello msgType = iota + 1
	msgData
	msgFecData
)

type ucastMessage struct {
	msgID uint32
}

type ucastHelloMessage struct {
	ucastMessage
	isFecMsg      bool
	fecRequired   uint8
	fecPieces     uint8
	fecInterleave uint8
	stripeSize    uint16
	name          string
}

func (u ucastHelloMessage) toNetByteOrder() ([]byte, error) {
	if len(u.name) > 512 {
		return nil, errors.New("maximum length of name is 512 bytes")
	}

	// Hello header is at least 11 and at maximum 523 bytes long
	var buf = bytes.NewBuffer(make([]byte, 0, 523))

	// First always serialize type information
	buf.WriteByte(byte(msgHello))

	if err := binary.Write(buf, binary.BigEndian, u.msgID); err != nil {
		return nil, err
	}

	// Write flags
	var flags uint8
	if u.isFecMsg {
		flags |= isFecFlag
	}
	buf.WriteByte(byte(flags))

	// Write fec config
	buf.WriteByte(byte(u.fecRequired))
	buf.WriteByte(byte(u.fecPieces))
	buf.WriteByte(byte(u.fecInterleave))

	if err := binary.Write(buf, binary.BigEndian, u.stripeSize); err != nil {
		return nil, err
	}

	if _, err := buf.WriteString(u.name); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *ucastHelloMessage) fromNetByteOrder(data []byte) error {
	if len(data) < 11 {
		return errors.New("message size is too small hello message is at least 10 bytes long")
	}

	if msgType(data[0]) != msgHello {
		return errors.New("data is not a hello message")
	}

	u.msgID = binary.BigEndian.Uint32(data[1:5])
	u.isFecMsg = data[5]&isFecFlag == isFecFlag
	u.fecRequired = data[6]
	u.fecPieces = data[7]
	u.fecInterleave = data[8]
	u.stripeSize = binary.BigEndian.Uint16(data[9:11])
	if len(data) > 11 {
		u.name = string(data[11:len(data)])
	}

	return nil
}
