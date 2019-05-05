package transport

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	fecMsgFlag   flags = 1 << 15
	moreDataFlag flags = 1 << 14
)

const (
	_               = iota
	headerMsg uint8 = iota
	dataMsg
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

const msgMinSz = 11

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
	if len(b) < msgMinSz {
		return fmt.Errorf("message: message must be at least %v bytes long", msgMinSz)
	}

	m.msgType = b[0]
	m.msgID = binary.BigEndian.Uint64(b[1:9])
	m.flags = flags(binary.BigEndian.Uint16(b[9:11]))

	return nil
}

func (m *message) reset() {
	m.flags, m.msgID, m.msgType = 0, 0, 0
}

// message + messageHeader length = 14
const msgHeaderMinSz = msgMinSz + 3
const msgFecHeaderMinSz = msgHeaderMinSz + 3

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
	if m.flags.isSet(fecMsgFlag) {
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

	if len(b) < msgDataMinSz {
		return fmt.Errorf("messageHeader: header must be at least %v bytes long", msgHeaderMinSz)
	}

	m.stripeL = binary.BigEndian.Uint16(b[11:13])
	m.contentT = b[13]

	if m.flags.isSet(fecMsgFlag) {
		if len(b) < msgFecHeaderMinSz {
			return fmt.Errorf("messageHeader: fec header must be at least %v bytes long", msgFecHeaderMinSz)
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

// message + messageData length = 19
const msgDataMinSz = msgMinSz + 8
const msgFecDataMinSz = msgDataMinSz + 3

type messageData struct {
	message
	stripeNr uint64
	fecPad   uint16
	fecNr    uint8
	data     []byte
}

func (m messageData) toBytes(b *bytes.Buffer) error {
	if err := m.message.toBytes(b); err != nil {
		return err
	}

	if err := binary.Write(b, binary.BigEndian, m.stripeNr); err != nil {
		return err
	}

	if m.flags.isSet(fecMsgFlag) {
		if err := binary.Write(b, binary.BigEndian, m.fecPad); err != nil {
			return err
		}
		b.WriteByte(m.fecNr)
	}

	if _, err := b.Write(m.data[:]); err != nil {
		return err
	}

	return nil
}

func (m *messageData) fromBytes(b []byte) error {
	if err := m.message.fromBytes(b); err != nil {
		return err
	}

	if len(b) < msgDataMinSz {
		return fmt.Errorf("messageData: must be at least %v bytes long", msgDataMinSz)
	}

	m.stripeNr = binary.BigEndian.Uint64(b[11:19])
	if m.flags.isSet(fecMsgFlag) {
		if len(b) < msgFecDataMinSz {
			return fmt.Errorf("messageData: fec data must be at least %v bytes long", msgFecDataMinSz)
		}
		m.fecPad = binary.BigEndian.Uint16(b[19:21])
		m.fecNr = b[21]
		m.data = append(m.data, b[22:]...)
	} else {
		m.data = append(m.data, b[19:]...)
	}

	return nil
}

func (m *messageData) reset() {
	m.message.reset()
	m.stripeNr, m.fecPad, m.fecNr = 0, 0, 0
	m.data = m.data[:0]
}
