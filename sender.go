package goucast

import (
	"io"
	"net"

	"golang.org/x/time/rate"
)

// Sender reads data from a writer and sends it over UDP to a receiver.
type Sender struct {
	nextMsgID uint32
	stripSize uint16
	limiter   rate.Limiter
	conn      net.UDPConn
}

func (s *Sender) Send(name string, r io.Reader) error {
	panic("not implemented yet")
}
