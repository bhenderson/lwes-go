package lwes

import (
    "bytes"
    "encoding/binary"
    "net"
    "fmt"
    "encoding/json"
)

// http://golang.org/doc/articles/json_and_go.html
type eventAttrs map[string]interface{}

type Event struct {
    // TODO should this be a normal struct?
    Name string
    Attributes eventAttrs
}

// NewEvent returns an initialized Event
func NewEvent(argv...string) *Event {
    e := &Event{Attributes: make(eventAttrs)}
    if argv != nil {
        e.Name = argv[0]
    }
    return e
}

// Iterator interface
func (e *Event) Iterator() eventAttrs {
    return e.Attributes
}

// Get an attribute
func (e *Event) Get(s string) interface{} {
    return e.Attributes[s]
}

// Originally this was meant to make setting Attributes a private function. But
// emitter uses json.Decode to set them and it looks nice.
func (e *Event) SetAttribute(name string, d interface{}) {
    // TODO validate types
    // Should we validate string length?
    switch v := d.(type) {
    default:
        e.Attributes[name] = v
    }
}

func (event *Event) toBytes() ([]byte, error) {
    buf := new(bytes.Buffer)
    var err error

    // TODO write errors
    write := func(d interface{}) bool {
        if d == nil { return true }
        err = binary.Write(buf, binary.BigEndian, d)
        return err == nil
    }

    writeRaw := func(d...interface{}) bool {
        if d == nil { return true }
        _, err = buf.Write(d[0].([]byte))
        return err == nil
    }

    writeKey := func(k string) bool {
        l := len(k)
        if l > MAX_SHORT_STRING_SIZE {
            err = fmt.Errorf("key length exceeds MAX_SHORT_STRING_SIZE(%d)", MAX_SHORT_STRING_SIZE)
            return false
        }
        if !write(byte(l)) { return false }
        _, err = buf.Write([]byte(k))
        return err == nil
    }

    writeAttr := func(k string, t byte, d interface{}, r...interface{}) bool {
        return writeKey(k) && write(t) && write(d) && writeRaw(r...)
    }

    // write name length
    // write name
    // write num Attributes
    if ! (writeKey(event.Name) && write(uint16(len(event.Attributes)))) {
        return nil, err
    }

    for key := range event.Attributes {
        switch v := event.Attributes[key].(type) {
        default:
            fmt.Printf("unknown key type: %T %#v\n", v,v)
        case uint8:
            writeAttr(key, LWES_U_INT_16_TOKEN, uint16(v))
        case *uint8:
            writeAttr(key, LWES_U_INT_16_TOKEN, uint16(*v))
        case uint16, *uint16:
            writeAttr(key, LWES_U_INT_16_TOKEN, v)
        case int8:
            writeAttr(key, LWES_INT_16_TOKEN, int16(v))
        case *int8:
            writeAttr(key, LWES_INT_16_TOKEN, int16(*v))
        case int16, *int16:
            writeAttr(key, LWES_INT_16_TOKEN, v)
        case uint32, *uint32:
            writeAttr(key, LWES_U_INT_32_TOKEN, v)
        case int32, *int32:
            writeAttr(key, LWES_INT_32_TOKEN, v)
        case string:
            writeAttr(key, LWES_STRING_TOKEN, uint16(len(v)), []byte(v))
        case *string:
            writeAttr(key, LWES_STRING_TOKEN, uint16(len(*v)), []byte(*v))
        case net.IP:
            if tmpIP := v.To4(); tmpIP != nil {
                b := []byte{tmpIP[3], tmpIP[2], tmpIP[1], tmpIP[0]}
                writeAttr(key, LWES_IP_ADDR_TOKEN, nil, b)
            }
        case *net.IP:
            if tmpIP := v.To4(); tmpIP != nil {
                b := []byte{tmpIP[3], tmpIP[2], tmpIP[1], tmpIP[0]}
                writeAttr(key, LWES_IP_ADDR_TOKEN, nil, b)
            }
        case int64, *int64, float64, *float64:
            writeAttr(key, LWES_INT_64_TOKEN, v)
        case uint64, *uint64:
            writeAttr(key, LWES_U_INT_64_TOKEN, v)
        case bool:
            var b byte
            if v { b = 1 } else { b = 0 }
            writeAttr(key, LWES_BOOLEAN_TOKEN, b)
        case *bool:
            var b byte
            if *v { b = 1 } else { b = 0 }
            writeAttr(key, LWES_BOOLEAN_TOKEN, b)
        // int and uint might be 32 or 64
        case int:
            writeAttr(key, LWES_INT_64_TOKEN, int64(v))
        case *int:
            writeAttr(key, LWES_INT_64_TOKEN, int64(*v))
        case uint:
            writeAttr(key, LWES_U_INT_64_TOKEN, uint64(v))
        case *uint:
            writeAttr(key, LWES_U_INT_64_TOKEN, uint64(*v))
        }
    }

    if err != nil { return nil, err }
    if buf.Len() > MAX_MSG_SIZE {
        return nil, fmt.Errorf("num bytes exceeds MAX_MSG_SIZE")
    }
    return buf.Bytes(), nil
}

func (event *Event) fromBytes(buf []byte) {
    p := bytes.NewBuffer(buf)

    // TODO read errors
    read := func(d interface{}) {
        binary.Read(p, binary.BigEndian, d)
    }

    // temporary types
    var tmpUint16 uint16
    var tmpInt16  int16
    var tmpUint32 uint32
    var tmpInt32  int32
    var tmpUint64 uint64
    var tmpInt64  int64

    var nameSize byte
    read(&nameSize)
    event.Name = string(p.Next(int(nameSize)))

    var attrSize uint16
    read(&attrSize)

    for i:=0; i < int(attrSize); i++ {
        var attrNameSize byte
        var attrName string
        var attrType byte

        read(&attrNameSize)
        attrName = string(p.Next(int(attrNameSize)))

        read(&attrType)

        switch attrType {
        case LWES_U_INT_16_TOKEN:
            read(&tmpUint16)
            event.Attributes[attrName] = tmpUint16
        case LWES_INT_16_TOKEN:
            read(&tmpInt16)
            event.Attributes[attrName] = tmpInt16
        case LWES_U_INT_32_TOKEN:
            read(&tmpUint32)
            event.Attributes[attrName] = tmpUint32
        case LWES_INT_32_TOKEN:
            read(&tmpInt32)
            event.Attributes[attrName] = tmpInt32
        case LWES_STRING_TOKEN:
            read(&tmpUint16)
            event.Attributes[attrName] = string(p.Next(int(tmpUint16)))
        case LWES_IP_ADDR_TOKEN:
            tmpIp := p.Next(4)
            // not sure if this is completely accurate
            event.Attributes[attrName] = net.IPv4(tmpIp[3], tmpIp[2], tmpIp[1], tmpIp[0])
        case LWES_INT_64_TOKEN:
            read(&tmpInt64)
            event.Attributes[attrName] = tmpInt64
        case LWES_U_INT_64_TOKEN:
            read(&tmpUint64)
            event.Attributes[attrName] = tmpUint64
        case LWES_BOOLEAN_TOKEN:
            event.Attributes[attrName] = 1 == p.Next(1)[0]
        }
    }
}

// PrettyString returns a "pretty" formatted string.
func (e *Event) PrettyString() string {
    var buf bytes.Buffer

    buf.WriteString(e.Name)
    buf.WriteString("\n")

    for key := range e.Attributes {
        buf.WriteString(key)
        buf.WriteString(": ")
        // gah
        buf.WriteString(fmt.Sprintln(e.Attributes[key]))
    }

    return buf.String()
}

// MarshalJSON returns a json byte array of name:attributes
// net.IP is base64 encoded
func (e *Event) MarshalJSON() (data []byte, err error) {
    m := make(eventAttrs)
    m[e.Name] = e.Attributes
    return json.Marshal(m)
}
