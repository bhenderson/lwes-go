package lwes

import (
    "net"
    "time"
)

type Emitter struct {
    Raddr *net.UDPAddr
    Heartbeat int8
    TTL int8
    Iface *net.Interface
    socket *net.UDPConn
}

var (
    startupEvent   = NewEvent("System::Startup")
    shutdownEvent  = NewEvent("System::Shutdown")
    heartbeatEvent = NewEvent("System::Heartbeat")
)

func NewEmitter(ip string, port int) (*Emitter, error) {
    raddr := &net.UDPAddr{
        IP: net.ParseIP("224.2.2.22"),
        Port: 12345,
    }
    e := &Emitter{Raddr: raddr, Heartbeat: 2, TTL: 3}

    soc, err := net.ListenUDP("udp", e.Raddr)
    if err != nil {
        return nil, err
    }

    e.socket = soc
    e.Emit(startupEvent)
    go e.emitHeartbeats()
    return e, nil
}

func (e *Emitter) Emit(event *Event) error {
    b, err := event.toBytes()

    if err != nil {
        return err
    }

    _, err = e.socket.WriteToUDP(b, e.Raddr)

    return err
}

// Close the emitter. Usually not needed.
func (e *Emitter) Close() {
    e.Emit(shutdownEvent)
    // test if open
    e.socket.Close()
}

func (e *Emitter) emitHeartbeats() {
    ticker := time.Tick(time.Duration(e.Heartbeat) * time.Second)
    for _ = range ticker {
        e.Emit(heartbeatEvent)
    }
}
