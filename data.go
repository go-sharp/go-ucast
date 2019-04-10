package goucast

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	// Just to ensure package is offline available
	_ "github.com/vivint/infectious"
)

const (
	fecFlag     = 1 << 7
	lastPkgFlag = 1 << 7
)

type msgType int8

const (
	msgTypeHello msgType = iota + 1
	msgTypeData
	msgTypeFecData
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

func (u ucastHelloMessage) toBytes() ([]byte, error) {
	if len(u.name) > 512 {
		return nil, errors.New("maximum length of name is 512 bytes")
	}

	// Hello header is at least 11 and at maximum 523 bytes long
	var buf = bytes.NewBuffer(make([]byte, 0, 523))

	// First always serialize type information
	buf.WriteByte(byte(msgTypeHello))

	if err := binary.Write(buf, binary.BigEndian, u.msgID); err != nil {
		return nil, err
	}

	// Write flags
	var flags uint8
	if u.isFecMsg {
		flags |= fecFlag
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

func (u *ucastHelloMessage) fromBytes(data []byte) error {
	if len(data) < 11 {
		return errors.New("message size is too small, hello message is at least 11 bytes long")
	}

	if msgType(data[0]) != msgTypeHello {
		return errors.New("data is not a hello message")
	}

	u.msgID = binary.BigEndian.Uint32(data[1:5])
	u.isFecMsg = data[5]&fecFlag == fecFlag
	u.fecRequired = data[6]
	u.fecPieces = data[7]
	u.fecInterleave = data[8]
	u.stripeSize = binary.BigEndian.Uint16(data[9:11])
	if len(data) > 11 {
		u.name = string(data[11:len(data)])
	}

	return nil
}

type ucastDataMessage struct {
	ucastMessage
	typ          msgType
	isLastPacket bool
	stripeNr     uint32
	pieceNr      uint8
	data         []byte
}

func (u ucastDataMessage) toBytes() ([]byte, error) {
	maxData := 1462
	if u.typ == msgTypeFecData {
		maxData = 1461
	}
	// Max length max udp size - ucastDataMessage header size
	if len(u.data) > maxData {
		return nil, fmt.Errorf("maximum length of data is %v bytes", maxData)
	}

	var buf = bytes.NewBuffer(make([]byte, 0, 1472))
	// Writing message type
	buf.WriteByte(byte(u.typ))
	// Write message id
	if err := binary.Write(buf, binary.BigEndian, u.msgID); err != nil {
		return nil, err
	}

	// Write flags
	var flags uint8
	if u.isLastPacket {
		flags |= lastPkgFlag
	}
	buf.WriteByte(flags)

	// Write stripe number
	if err := binary.Write(buf, binary.BigEndian, u.stripeNr); err != nil {
		return nil, err
	}

	// If fec data message write pieceNr
	if u.typ == msgTypeFecData {
		buf.WriteByte(u.pieceNr)
	}

	// Write data
	buf.Write(u.data)
	return buf.Bytes(), nil
}

func (u *ucastDataMessage) fromBytes(data []byte) error {
	if len(data) < 10 {
		return errors.New("message size is too small, data message is at least 10 bytes long")
	}

	typ := msgType(data[0])
	if typ != msgTypeData && typ != msgTypeFecData {
		return errors.New("data is neither a data nor a fec data message")
	}

	if typ == msgTypeFecData && len(data) < 11 {
		return errors.New("message size is too small, fec data message is at least 11 bytes long")
	}

	u.msgID = binary.BigEndian.Uint32(data[1:5])
	u.isLastPacket = data[5]&lastPkgFlag == lastPkgFlag
	u.stripeNr = binary.BigEndian.Uint32(data[5:9])

	if typ == msgTypeFecData {
		u.pieceNr = data[9]
		u.data = data[10:]
	} else {
		u.data = data[9:]
	}

	return nil
}
