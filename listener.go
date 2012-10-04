package lwes

import (
    "net"
    "log"
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

func Listener() {
    // TODO change to ListenUDP if addr is zero
    socket, err := net.ListenMulticastUDP("udp4", nil, &net.UDPAddr{
        // IP:   net.IPv4zero,
        IP:   net.IPv4(224,2,2,22),
        Port: 12345,
    })
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

        log.Println(raddr)
        log.Println(buff[:read])
    }
}
