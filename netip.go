package lwes

import (
	"encoding/json"
	"net"
)

type IP4 interface {
	To4() net.IP
}

func NewNetIP(p []byte) *NetIP {
	var ip net.IP
	if len(p) > 3 {
		ip = net.IPv4(p[3], p[2], p[1], p[0])
	}
	return &NetIP{&net.IPAddr{IP: ip}}
}

type NetIP struct {
	net.Addr
}

func (ip *NetIP) MarshalJSON() ([]byte, error) {
	return json.Marshal(ip.String())
}

func addrIP(addr net.Addr) *NetIP {
	switch x := addr.(type) {
	case *net.UDPAddr:
		addr = &net.IPAddr{IP: x.IP.To4()}
	}
	return &NetIP{addr}
}

func addrPort(addr net.Addr) int {
	switch x := addr.(type) {
	case *net.UDPAddr:
		return x.Port
	}
	return 0
}
