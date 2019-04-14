package transport

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	fecMsg flags = 1 << 15
)

type flags uint16

func (f flags) isSet(fl ...flags) bool {
	for i := range fl {
		if f&fl[i] == fl[i] {
			continue
		}
		return false
	}
	return true
}

type message struct {
	msgID   uint64
	msgType uint8
	flags   flags
}

func (m message) toBytes(b *bytes.Buffer) error {
	// first serialize message type
	b.WriteByte(m.msgType)

	// second serialize the message id
	if err := binary.Write(b, binary.BigEndian, m.msgID); err != nil {
		return err
	}

	// last serialize flags
	if err := binary.Write(b, binary.BigEndian, m.flags); err != nil {
		return err
	}

	return nil
}

func (m *message) fromBytes(b []byte) error {
	if len(b) < 11 {
		return errors.New("message: message must be at least 11 bytes long")
	}

	m.msgType = b[0]
	m.msgID = binary.BigEndian.Uint64(b[1:9])
	m.flags = flags(binary.BigEndian.Uint16(b[9:11]))

	return nil
}

func (m *message) reset() {
	m.flags, m.msgID, m.msgType = 0, 0, 0
}

// messageHeader indicates the start of a new transmission.
type messageHeader struct {
	message
	stripeL  uint16
	contentT uint8
	fecTot   uint8
	fecRec   uint8
	fecInt   uint8
}

func (m messageHeader) toBytes(b *bytes.Buffer) error {
	if err := m.message.toBytes(b); err != nil {
		return err
	}

	if err := binary.Write(b, binary.BigEndian, m.stripeL); err != nil {
		return err
	}

	b.WriteByte(m.contentT)
	if m.flags.isSet(fecMsg) {
		b.WriteByte(m.fecTot)
		b.WriteByte(m.fecRec)
		b.WriteByte(m.fecInt)
	}

	return nil
}

func (m *messageHeader) fromBytes(b []byte) error {
	if err := m.message.fromBytes(b); err != nil {
		return err
	}

	if len(b) < 14 {
		return errors.New("messageHeader: header must be at least 14 bytes long")
	}

	m.stripeL = binary.BigEndian.Uint16(b[11:13])
	m.contentT = b[13]

	if m.flags.isSet(fecMsg) {
		if len(b) < 17 {
			return errors.New("messageHeader: fec header must be at least 17 bytes long")
		}
		m.fecTot = b[14]
		m.fecRec = b[15]
		m.fecInt = b[16]
	}

	return nil
}

func (m *messageHeader) reset() {
	m.message.reset()
	m.fecInt, m.fecRec, m.fecTot = 0, 0, 0
	m.stripeL, m.contentT = 0, 0
}
