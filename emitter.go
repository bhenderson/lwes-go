package lwes

import (
    "net"
    "fmt"
)

type Emitter struct {
    Raddr *net.UDPAddr
    Heartbeat int8
    TTL int8
    Iface *net.Interface
    socket *net.UDPConn
}

func NewEmitter(ip string, port int, heartbeat int8, ttl int8, iface *net.Interface) (*Emitter, error) {
    raddr := &net.UDPAddr{
        IP: net.ParseIP(ip),
        Port: port,
    }
    e := &Emitter{Raddr: raddr, Heartbeat: heartbeat, TTL: ttl, Iface: iface}
    err := e.setSocket()

    if err != nil {
        return nil, err
    }

    return e, nil
}

func (e *Emitter) setSocket() error {
    soc, err := net.ListenUDP("udp", e.Raddr)

    if err != nil {
        return err
    }

    e.socket = soc

    return nil
}

func (e *Emitter) Emit(event *Event) error {
    b, err := event.toBytes()

    if err != nil {
        return err
    }

    fmt.Printf("sending: %v\n", b)
    _, err = e.socket.WriteToUDP(b, e.Raddr)

    return err
}
