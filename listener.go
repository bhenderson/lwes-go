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

// http://golang.org/doc/articles/json_and_go.html
type eventAttrs map[string]interface{}

type Event struct {
    name string
    attributes eventAttrs
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
    return Event{attributes: make(eventAttrs)}
}

func deserializeEvent(event *Event, buf []byte) {
    // TODO should the interface for this func be different, ie, should it implement
    // Writer?
    p := bytes.NewBuffer(buf)

    var nameSize byte
    binary.Read(p, binary.BigEndian, &nameSize)

    event.name = string(p.Next(int(nameSize)))

    var attrSize uint16
    binary.Read(p, binary.BigEndian, &attrSize)

    // temporary types
    // var tmpByte   byte
    var tmpUint16 uint16
    var tmpInt16  int16
    var tmpUint32 uint32
    var tmpInt32  int32
    var tmpUint64 uint64
    var tmpInt64  int64
    // var tmpBool   byte
    // var tmpIpaddr net.IP
    // var tmpString string

    for i:=0; i < int(attrSize); i++ {
        var attrNameSize byte
        var attrName string
        var attrType byte

        binary.Read(p, binary.BigEndian, &attrNameSize)
        // TODO should we camelCase attrName?
        attrName = string(p.Next(int(attrNameSize)))

        binary.Read(p, binary.BigEndian, &attrType)

        // log.Println(attrName, attrType)

        switch int(attrType) {
        case 1: // LWES_U_INT_16_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint16)
            event.attributes[attrName] = tmpUint16
        case 2: // LWES_INT_16_TOKEN
            binary.Read(p, binary.BigEndian, &tmpInt16)
            event.attributes[attrName] = tmpInt16
        case 3: // LWES_U_INT_32_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint32)
            event.attributes[attrName] = tmpUint32
        case 4: // LWES_INT_32_TOKEN
            binary.Read(p, binary.BigEndian, &tmpInt32)
            event.attributes[attrName] = tmpInt32
        case 5: // LWES_STRING_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint16)
            event.attributes[attrName] = string(p.Next(int(tmpUint16)))
        case 6: // LWES_IP_ADDR_TOKEN
            tmpIp := p.Next(4)
            // not sure if this is completely accurate
            event.attributes[attrName] = net.IPv4(tmpIp[3], tmpIp[2], tmpIp[1], tmpIp[0])
        case 7: // LWES_INT_64_TOKEN
            binary.Read(p, binary.BigEndian, &tmpInt64)
            event.attributes[attrName] = tmpInt64
        case 8: // LWES_U_INT_64_TOKEN
            binary.Read(p, binary.BigEndian, &tmpUint64)
            event.attributes[attrName] = tmpUint64
        case 9: // LWES_BOOLEAN_TOKEN
            event.attributes[attrName] = 1 == p.Next(1)[0]
        }
    }
}

// Name returns the name or class of an event. This is separate from an attribute
func (e Event) Name() string {
    return e.name
}

// Iterator interface
func (e Event) Iter() eventAttrs {
    return e.attributes
}

// Get an attribute
func (e Event) Get(s string) interface{} {
    return e.attributes[s]
}
