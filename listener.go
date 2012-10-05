package lwes

import (
    "net"
    "log"
    "time"
    "bytes"
    "encoding/binary"
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
    name string
    // http://golang.org/doc/articles/json_and_go.html
    attributes map[string]interface{}
}

// an action is a listener callback
type listenerAction func(event *Event)

//Listener starts listening on ip_addr and port
func Listener(laddr *net.UDPAddr, callback listenerAction) {
    // pointless if no callback func
    if callback == nil {
        return
    }

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
        buf := make([]byte, MAX_MSG_SIZE)
        read, raddr, err := socket.ReadFromUDP(buf)

        if err != nil {
            log.Fatal(err)
        }

        time := time.Now()

        event := NewEvent()
        event.attributes["receiptTime"] = time
        event.attributes["senderIp"]    = raddr.IP
        event.attributes["senderPort"]  = raddr.Port

        deserializeEvent(&event, buf[:read])

        callback(&event)
    }
}

// NewEvent returns an initialized Event
func NewEvent() Event {
    return Event{attributes: make(map[string]interface{})}
}

func deserializeEvent(event *Event, buf []byte) {
    p := bytes.NewBuffer(buf)

    var nameSize byte
    binary.Read(p, binary.BigEndian, &nameSize)

    event.name = string(p.Next(int(nameSize)))

    var attrSize uint16
    binary.Read(p, binary.BigEndian, &attrSize)

    for i:=0; i < int(attrSize); i++ {
        var attrNameSize byte
        var attrName string
        var attrType byte

        binary.Read(p, binary.BigEndian, &attrNameSize)
        attrName = string(p.Next(int(attrNameSize)))

        binary.Read(p, binary.BigEndian, &attrType)

        log.Println(attrName, attrType)

        switch int(attrType) {
        case 1: // LWES_U_INT_16_TOKEN
        case 2: // LWES_INT_16_TOKEN
        case 3: // LWES_U_INT_32_TOKEN
        case 4: // LWES_INT_32_TOKEN
        case 5: // LWES_STRING_TOKEN
        case 6: // LWES_IP_ADDR_TOKEN
        case 7: // LWES_INT_64_TOKEN
        case 8: // LWES_U_INT_64_TOKEN
        case 9: // LWES_BOOLEAN_TOKEN
        }
    }
}
