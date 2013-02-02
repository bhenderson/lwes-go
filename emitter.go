package lwes

import (
    "net"
    "time"
    "os"
    "os/signal"
)

type Emitter struct {
    Heartbeat int8
    TTL int8
    closer chan bool
    socket *Conn
}

var (
    startupEvent   = NewEvent("System::Startup")
    shutdownEvent  = NewEvent("System::Shutdown")
    heartbeatEvent = NewEvent("System::Heartbeat")
)

func NewEmitter(udp string, iface ...*net.Interface) (*Emitter, error) {
    conn, err := NewConn(udp, true, iface...)

    if err != nil {
        return nil, err
    }

    e := &Emitter{Heartbeat: 1, TTL: 3, closer: make(chan bool), socket: conn}
    go e.emitHeartbeats()
    return e, nil
}

func (e *Emitter) Emit(event *Event) error {
    b, err := event.toBytes()

    if err != nil {
        return err
    }

    // fmt.Printf("%s\n", b[0])
    // fmt.Printf("bytes: %#v\n", b)

    _, err = e.socket.Write(b)

    return err
}

// Close the emitter. Usually not needed.
func (e *Emitter) Close() {
    e.closer <- true
    // test if open
    e.socket.Close()
}

// Send a heartbeat event every Heartbeat seconds.
// if Heartbeat is 0, don't send any.
// shutdown event may not send (os.Exit)
func (e *Emitter) emitHeartbeats() {

    e.Emit(startupEvent)
    defer e.Emit(shutdownEvent)

    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, os.Kill)

    ticker := time.Tick(time.Duration(e.Heartbeat) * time.Second)
    for {
        select {
        case <- ticker:
            e.Emit(heartbeatEvent)
        case <- c:
            // received SIGINT
            return
        case <- e.closer:
            return
        }
    }
}
