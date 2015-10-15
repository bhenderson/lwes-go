package lwes

import (
	"fmt"
	"net"
	"net/url"
	"os"
)

type Conn interface {
	Close() error
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
	case "unixgram":
		conn, err = bindUnix(laddr, ifi, emitter)
	case "ip", "ip4", "ip6":
	default:
		err = fmt.Errorf("%q is not supported")
	}

	return conn, err
}

type udpConn struct {
	addr *net.UDPAddr
	*net.UDPConn
}

func (c *udpConn) Read(b []byte) (int, net.Addr, error) {
	return c.UDPConn.ReadFromUDP(b)
}

func (c *udpConn) Write(b []byte) (int, error) {
	return c.UDPConn.WriteToUDP(b, c.addr)
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

	c.UDPConn = conn
	return c, err
}

type unixConn struct {
	addr *net.UnixAddr
	*net.UnixConn
}

func (c *unixConn) Close() error {
	defer os.Remove(c.addr.String())

	return c.UnixConn.Close()
}

func (c *unixConn) Read(b []byte) (int, net.Addr, error) {
	return c.UnixConn.ReadFromUnix(b)
}

func (c *unixConn) Write(b []byte) (int, error) {
	return c.UnixConn.Write(b)
}

func bindUnix(laddr string, iface *net.Interface, emitter bool) (*unixConn, error) {
	addr, err := net.ResolveUnixAddr("unixgram", laddr)

	if err != nil {
		return nil, err
	}

	c := &unixConn{addr: addr}

	if emitter {
		c.UnixConn, err = net.DialUnix("unixgram", nil, addr)
	} else {
		c.UnixConn, err = net.ListenUnixgram("unixgram", addr)
	}

	return c, err
}
