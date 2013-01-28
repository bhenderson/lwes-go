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
    closer chan bool
    socket *net.UDPConn
}

var (
    startupEvent   = NewEvent("System::Startup")
    shutdownEvent  = NewEvent("System::Shutdown")
    heartbeatEvent = NewEvent("System::Heartbeat")
)

func NewEmitter(ip string, port int) (*Emitter, error) {
    raddr := &net.UDPAddr{
        IP: net.ParseIP(ip),
        Port: port,
    }
    e := &Emitter{Raddr: raddr, Heartbeat: 1, TTL: 3}

    soc, err := net.ListenUDP("udp4", e.Raddr)
    if err != nil {
        return nil, err
    }

    e.closer = make(chan bool)
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

    // fmt.Printf("%s\n", b[0])
    // fmt.Printf("%#v\n", b)

    // TODO toBytes is working correctly but emitter is still broken.
    // that said, if I send eventSlice (from test) it works!
    _, err = e.socket.WriteToUDP(b, e.Raddr)

    return err
}

// Close the emitter. Usually not needed.
func (e *Emitter) Close() {
    e.closer <- true
    // test if open
    e.socket.Close()
}

// Send a heartbeat event every Heartbeat seconds.
func (e *Emitter) emitHeartbeats() {

    defer e.Emit(shutdownEvent)

    ticker := time.Tick(time.Duration(e.Heartbeat) * time.Second)
    for {
        select {
        case <- ticker:
            e.Emit(heartbeatEvent)
        case <- e.closer:
            return
        }
    }
}
