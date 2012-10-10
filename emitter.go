package lwes

import (
    "net"
    "fmt"
)

type Emitter struct {
    Address string
    Port int
    Heartbeat int8
    TTL int8
    Iface *net.Interface
    socket net.Conn
}

func NewEmitter(ip string, port int, heartbeat int8, ttl int8, iface *net.Interface) (*Emitter, error) {
    e := &Emitter{Address: ip, Port: port, Heartbeat: heartbeat, TTL: ttl, Iface: iface}
    err := e.dial()

    if err != nil {
        return nil, err
    }

    return e, nil
}

func (e *Emitter) dial() error {
    soc, err := net.Dial("udp", net.JoinHostPort(e.Address, fmt.Sprintf("%d", e.Port)))

    if err != nil {
        return err
    }

    e.socket = soc

    return nil
}

func (e *Emitter) Emit(event *Event) error {
    b, err := event.ToBytes()

    if err != nil {
        return err
    }

    i, err := e.socket.Write(b)

    fmt.Println(i)
    fmt.Println(string(b))

    return err
}
