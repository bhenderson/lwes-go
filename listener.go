package lwes

import (
	"fmt"
	"net"
	"time"
)

type Listener struct {
	socket Conn
}

// NewListener creates a new Listener and binds to ip and port and iface
func NewListener(udp string, iface ...*net.Interface) (*Listener, error) {
	conn, err := NewConn(udp, false, iface...)
	l := &Listener{socket: conn}

	return l, err
}

// Recv receives an event
func (l *Listener) Recv() (*Event, error) {
	if l.socket == nil {
		return nil, fmt.Errorf("socket is not bound")
	}

	buf := make([]byte, MAX_MSG_SIZE)
	read, raddr, err := l.socket.Read(buf)

	if err != nil {
		return nil, err
	}

	time := time.Now()

	event := NewEvent()
	event.fromBytes(buf[:read])

	event.Attributes["ReceiptTime"] = time
	event.Attributes["SenderIP"] = addrIP(raddr)
	event.Attributes["SenderPort"] = addrPort(raddr)

	return event, nil
}
