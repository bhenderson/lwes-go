package lwes

import (
    "net"
    "log"
    "time"
)

const (
    /*
        maximum datagram size for UDP is 64K minus IP layer overhead which is
        20 bytes for IP header, and 8 bytes for UDP header, so this value
        should be

        65535 - 28 = 65507
     */
    MAX_MSG_SIZE = 65507
)

type Event interface {
    Name() string
}

func Listener(ip_addr net.IP, port int) {
    laddr := net.UDPAddr{ip_addr, port}

    var socket *net.UDPConn
    var err error

    if ip_addr.IsMulticast() {
        socket, err = net.ListenMulticastUDP("udp4", nil, &laddr)
    } else {
        socket, err = net.ListenUDP("udp4", &laddr)
    }

    if err != nil {
        log.Fatal(err)
    }
    defer socket.Close()

    for {
        buff := make([]byte, MAX_MSG_SIZE)
        read, raddr, err := socket.ReadFromUDP(buff)

        if err != nil {
            log.Fatal(err)
        }

        time := time.Now()

        log.Println(raddr)
        log.Println(time)
        log.Println(buff[:read])
    }
}
