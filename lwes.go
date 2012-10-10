package lwes

import (
    "fmt"
    "net"
)

func toIP(ip interface{}) (laddr *net.IP, err error) {
    switch t := ip.(type) {
    default:
        return nil, fmt.Errorf("ip is invalid type %T", t)
    case string:
        i := net.ParseIP(t)
        if ip != nil {
            laddr = &i
        } else {
            return nil, fmt.Errorf("ip is invalid")
        }
    case *net.IP:
        laddr = t
    }
    return
}
