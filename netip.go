package lwes

import (
    "encoding/json"
    "net"
)

type netIP net.IP

func (ip netIP) MarshalJSON() ([]byte, error) {
    return json.Marshal(ip.String())
}

func (ip netIP) String() string {
    return net.IP(ip).String()
}
