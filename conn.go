package lwes

import (
    "net"
)

type Conn struct {
    addr *net.UDPAddr
    iface *net.Interface
    socket *net.UDPConn
}

// Bind starts listening on udp addr.
// if emitter is true, bind to :0 and write to c.addr
// else bind to c.addr and read from c.addr
func (c *Conn) Bind(emitter bool) error {
    var addr *net.UDPAddr
    var conn *net.UDPConn
    var err error

    if emitter {
        addr, err = net.ResolveUDPAddr("udp", ":0")
    } else {
        addr = c.addr
    }

    if err != nil {
        return err
    }

    if addr.IP != nil && addr.IP.IsMulticast() {
        conn, err = net.ListenMulticastUDP("udp", nil, addr)
    } else {
        conn, err = net.ListenUDP("udp", addr)
    }

    if err != nil {
        return err
    }

    c.socket = conn
    return nil
}

// Close closes the socket. Make sure to call this if calling bind explicitely.
func (c *Conn) Close() {
    if c.socket != nil {
        c.socket.Close()
    }
}

func (c *Conn) Read(b []byte) (int, *net.UDPAddr, error) {
    return c.socket.ReadFromUDP(b)
}

func (c *Conn) Write(b []byte) (int, error) {
    return c.socket.WriteToUDP(b, c.addr)
}

func NewConn(udp string, emitter bool, iface ...*net.Interface) (*Conn, error) {
    addr, err := net.ResolveUDPAddr("udp", udp)

    if err != nil {
        return nil, err
    }

    var ifi *net.Interface

    if iface != nil {
        ifi = iface[0]
    }

    c := &Conn{addr: addr, iface: ifi}

    err = c.Bind(emitter)

    if err != nil {
        return nil, err
    }

    return c, nil
}
