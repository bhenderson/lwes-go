package lwes

import (
    "net"
    "fmt"
    "time"
)

type Listener struct {
    IP *net.IP
    Port int
    Iface *net.Interface
    socket *net.UDPConn
}

type listenerAction func(*Event, error)

// NewListener creates a new Listener and binds to ip and port and iface
func NewListener(ip interface{}, port int, iface ...*net.Interface) (*Listener, error) {
    var ifi *net.Interface

    laddr, err := toIP(ip)

    if err != nil {
        return nil, err
    }

    if iface != nil {
        ifi = iface[0]
    }

    l := &Listener{IP: laddr, Port: port, Iface: ifi}

    err = l.bind()

    if err != nil {
        return nil, err
    }

    return l, nil
}

// Close closes the socket. Make sure to call this if calling bind explicitely.
func (l *Listener) Close() {
    if l.socket != nil {
        l.socket.Close()
    }
}

// Each takes a listenerAction and gives it an *Event. See listenerAction.
func (l *Listener) Each(action listenerAction) {
    defer l.Close()

    for { action(l.Recv()) }
}

// Recv receives an event
func (l *Listener) Recv() (*Event, error) {
    if l.socket == nil {
        return nil, fmt.Errorf("socket is not bound")
    }

    buf := make([]byte, MAX_MSG_SIZE)
    read, raddr, err := l.socket.ReadFromUDP(buf)

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

//bind starts listening on ip and port
func (l *Listener) bind() error {
    var socket *net.UDPConn
    var err error

    laddr := &net.UDPAddr{
        IP: *l.IP,
        Port: l.Port,
    }

    if l.IP.IsMulticast() {
        socket, err = net.ListenMulticastUDP("udp4", l.Iface, laddr)
    } else {
        socket, err = net.ListenUDP("udp4", laddr)
    }

    if err != nil {
        return err
    }

    l.socket = socket
    return nil
}
