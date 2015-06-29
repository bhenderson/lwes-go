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
