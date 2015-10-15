package lwes

import (
	"encoding/json"
	"net"
)

type NetIP net.IP

func (ip NetIP) MarshalJSON() ([]byte, error) {
	return json.Marshal(ip.String())
}

func (ip NetIP) String() string {
	return net.IP(ip).String()
}

func addrIP(addr net.Addr) NetIP {
	switch x := addr.(type) {
	case *net.UDPAddr:
		return NetIP(x.IP.To16())
	case *net.IPAddr:
		return NetIP(x.IP.To16())
	}
	return NetIP{}
}

func addrPort(addr net.Addr) int {
	switch x := addr.(type) {
	case *net.UDPAddr:
		return x.Port
	}
	return 0
}
