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

// messageHeader indicates the start of a new transmission.
type messageHeader struct {
	msgID   uint64
	msgType uint8
	flags   flags
	stripeL uint16
	fecTot  uint8
	fecRec  uint8
	fecInt  uint8
}

func (m messageHeader) toBytes(b *bytes.Buffer) error {

	if err := binary.Write(b, binary.BigEndian, m.msgID); err != nil {
		return err
	}

	b.WriteByte(m.msgType)

	if err := binary.Write(b, binary.BigEndian, m.flags); err != nil {
		return err
	}

	if err := binary.Write(b, binary.BigEndian, m.stripeL); err != nil {
		return err
	}

	if m.flags.isSet(fecMsg) {
		b.WriteByte(m.fecTot)
		b.WriteByte(m.fecRec)
		b.WriteByte(m.fecInt)
	}

	return nil
}

func (m *messageHeader) fromBytes(b []byte) error {
	m.reset()
	if len(b) < 13 {
		return errors.New("messageHeader: header must be at least 13 bytes long")
	}

	m.msgID = binary.BigEndian.Uint64(b[0:8])
	m.msgType = b[8]
	m.flags = flags(binary.BigEndian.Uint16(b[9:11]))
	m.stripeL = binary.BigEndian.Uint16(b[11:13])

	if m.flags.isSet(fecMsg) {
		if len(b) < 16 {
			return errors.New("messageHeader: fec header must be at least 16 bytes long")
		}
		m.fecTot = b[13]
		m.fecRec = b[14]
		m.fecInt = b[15]
	}

	return nil
}

func (m *messageHeader) reset() {
	m.fecInt, m.fecRec, m.fecTot = 0, 0, 0
	m.flags, m.msgID, m.stripeL = 0, 0, 0
	m.msgType = 0
}
