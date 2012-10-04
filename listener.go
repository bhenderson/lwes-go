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

type Event struct {
    Name string
    // http://golang.org/doc/articles/json_and_go.html
    attributes map[string]interface{}
}

// an action is a listener callback
type listenerAction func(event *Event)

//Listener starts listening on ip_addr and port
func Listener(laddr *net.UDPAddr, callback listenerAction) {
    var socket *net.UDPConn
    var err error

    if laddr.IP.IsMulticast() {
        socket, err = net.ListenMulticastUDP("udp4", nil, laddr)
    } else {
        socket, err = net.ListenUDP("udp4", laddr)
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

        event := NewEvent()
        event.attributes["receiptTime"] = time
        event.attributes["senderIp"]    = raddr.IP
        event.attributes["senderPort"]  = raddr.Port

        deserializeEvent(&event, buff[:read])

        callback(&event)
    }
}

// NewEvent returns an initialized Event
func NewEvent() Event {
    return Event{attributes: make(map[string]interface{})}
}

func deserializeEvent(event *Event, buff []byte) {

    log.Println(buff)
}
