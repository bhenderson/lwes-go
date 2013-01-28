package lwes

import (
    "net"
    "fmt"
    "time"
)

type Listener struct {
    socket *Conn
}

type listenerAction func(*Event, error)

// NewListener creates a new Listener and binds to ip and port and iface
func NewListener(ip interface{}, port int, iface ...*net.Interface) (*Listener, error) {
    conn, err := NewConn(ip, port, iface...)
    l := &Listener{socket: conn}

    return l, err
}

// Each takes a listenerAction and gives it an *Event. See listenerAction.
func (l *Listener) Each(action listenerAction) {
    defer l.socket.Close()

    for { action(l.Recv()) }
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

    fmt.Printf("%#v", buf[:read])

    event := NewEvent()
    event.fromBytes(buf[:read])

    event.Attributes["receiptTime"] = time
    event.Attributes["senderIp"]    = raddr.IP.To16()
    event.Attributes["senderPort"]  = raddr.Port


    return event, nil
}
