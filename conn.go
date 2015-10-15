package lwes

import (
	"fmt"
	"net"
	"net/url"
)

type Conn interface {
	Close()
	Read(b []byte) (int, net.Addr, error)
	Write(b []byte) (int, error)
}

func NewConn(laddr string, emitter bool, iface ...*net.Interface) (Conn, error) {
	u, err := url.Parse(laddr)
	if err != nil {
		return nil, err
	}
	nt := u.Scheme
	if nt == "" {
		nt = "udp"
	}
	laddr = u.Host + u.Path

	var ifi *net.Interface

	if len(iface) > 0 {
		ifi = iface[0]
	}

	var conn Conn

	switch nt {
	case "udp", "udp4", "udp6":
		conn, err = bindUDP(laddr, ifi, emitter)
	case "tcp", "tcp4", "tcp6":
	case "unix", "unixgram", "unixpacket":
	case "ip", "ip4", "ip6":
	default:
		err = fmt.Errorf("%q is not supported")
	}

	return conn, err
}

type udpConn struct {
	addr   *net.UDPAddr
	socket *net.UDPConn
}

func (c *udpConn) Close() {
	if c.socket != nil {
		c.socket.Close()
	}
}

func (c *udpConn) Read(b []byte) (int, net.Addr, error) {
	return c.socket.ReadFromUDP(b)
}

func (c *udpConn) Write(b []byte) (int, error) {
	return c.socket.WriteToUDP(b, c.addr)
}

func bindUDP(laddr string, iface *net.Interface, emitter bool) (*udpConn, error) {
	addr, err := net.ResolveUDPAddr("udp", laddr)

	if err != nil {
		return nil, err
	}

	c := &udpConn{addr: addr}

	if emitter {
		addr, err = net.ResolveUDPAddr("udp", ":0")
	}

	if err != nil {
		return nil, err
	}

	var conn *net.UDPConn

	if addr.IP != nil && addr.IP.IsMulticast() {
		conn, err = net.ListenMulticastUDP("udp", iface, addr)
	} else {
		conn, err = net.ListenUDP("udp", addr)
	}

	c.socket = conn
	return c, err
}
