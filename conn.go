package lwes

import (
    "net"
    "fmt"
)

type Conn struct {
    IP *net.IP
    Port int
    Iface *net.Interface
    Socket *net.UDPConn
    addr *net.UDPAddr
}

//bind starts listening on ip and port
func (c *Conn) Bind() error {
    var socket *net.UDPConn
    var err error

    addr := &net.UDPAddr{
        IP: *c.IP,
        Port: c.Port,
    }
    c.addr = addr

    if c.IP.IsMulticast() {
        socket, err = net.ListenMulticastUDP("udp4", c.Iface, addr)
    } else {
        socket, err = net.ListenUDP("udp4", addr)
    }

    if err != nil {
        return err
    }

    c.Socket = socket
    return nil
}

// Close closes the socket. Make sure to call this if calling bind explicitely.
func (c *Conn) Close() {
    if c.Socket != nil {
        c.Socket.Close()
    }
}

func (c *Conn) Read(b []byte) (int, *net.UDPAddr, error) {
    return c.Socket.ReadFromUDP(b)
}

func (c *Conn) Write(b []byte) (int, error) {
    return c.Socket.WriteToUDP(b, c.addr)
}

func NewConn(ip interface{}, port int, iface ...*net.Interface) (*Conn, error) {
    var ifi *net.Interface

    addr, err := toIP(ip)

    if err != nil {
        return nil, err
    }

    if iface != nil {
        ifi = iface[0]
    }

    c := &Conn{IP: addr, Port: port, Iface: ifi}

    err = c.Bind()

    if err != nil {
        return nil, err
    }

    return c, nil
}

func toIP(ip interface{}) (addr *net.IP, err error) {
    switch t := ip.(type) {
    default:
        return nil, fmt.Errorf("ip is invalid type %T", t)
    case string:
        i := net.ParseIP(t)
        if ip != nil {
            addr = &i
        } else {
            return nil, fmt.Errorf("ip is invalid")
        }
    case *net.IP:
        addr = t
    }
    return
}
