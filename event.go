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
    attributes eventAttrs
}

// NewEvent returns an initialized Event
func NewEvent(argv...string) *Event {
    e := &Event{attributes: make(eventAttrs)}
    if argv != nil {
        e.Name = argv[0]
    }
    return e
}

// Iterator interface
func (e *Event) Iterator() eventAttrs {
    return e.attributes
}

// Get an attribute
func (e *Event) Get(s string) interface{} {
    return e.attributes[s]
}

func (e *Event) SetAttribute(name string, d interface{}) {
    // TODO validate types
    // Should we validate string length?
    switch v := d.(type) {
    default:
        e.attributes[name] = v
    }
}

func (event *Event) toBytes() ([]byte, error) {
    buf := new(bytes.Buffer)
    var err error

    // TODO write errors
    write := func(d interface{}) bool {
        err = binary.Write(buf, binary.BigEndian, d)
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

    writeAttr := func(k string, t int, d interface{}) bool {
        return writeKey(k) && write(byte(t)) && write(d)
    }

    // write name length
    // write name
    // write num attributes
    if ! (writeKey(event.Name) && write(uint16(len(event.attributes)))) {
        return nil, err
    }

    for key := range event.attributes {
        switch v := event.attributes[key].(type) {
        default:
            // fmt.Printf("unknown key type: %T %#v\n", v,v)
        case uint8:
            writeAttr(key, 1, uint16(v))
        case *uint8:
            writeAttr(key, 1, uint16(*v))
        case uint16, *uint16:
            writeAttr(key, 1, v)
        case int8:
            writeAttr(key, 2, int16(v))
        case *int8:
            writeAttr(key, 2, int16(*v))
        case int16, *int16:
            writeAttr(key, 2, v)
        case uint32, *uint32:
            writeAttr(key, 3, v)
        case int32, *int32:
            writeAttr(key, 4, v)
        case string:
            if writeAttr(key, 5, uint16(len(v))) {
                buf.Write([]byte(v))
            }
        case *string:
            if writeAttr(key, 5, uint16(len(*v))) {
                buf.Write([]byte(*v))
            }
        case net.IP:
            if writeKey(key) && write(byte(6)) {
                tmpIP := v.To4()
                buf.Write([]byte{tmpIP[3], tmpIP[2], tmpIP[1], tmpIP[0]})
            }
        case *net.IP:
            if writeKey(key) && write(byte(6)) {
                tmpIP := v.To4()
                buf.Write([]byte{tmpIP[3], tmpIP[2], tmpIP[1], tmpIP[0]})
            }
        case int64, *int64:
            writeAttr(key, 7, v)
        case uint64, *uint64:
            writeAttr(key, 8, v)
        case bool:
            var b int
            if v { b = 1 } else { b = 0 }
            writeAttr(key, 9, byte(b))
        case *bool:
            var b int
            if *v { b = 1 } else { b = 0 }
            writeAttr(key, 9, byte(b))
        // int and uint might be 32 or 64
        case int:
            writeAttr(key, 7, int64(v))
        case *int:
            writeAttr(key, 7, int64(*v))
        case uint:
            writeAttr(key, 8, uint64(v))
        case *uint:
            writeAttr(key, 8, uint64(*v))
        }

        if err != nil { return nil, err }
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

        switch int(attrType) {
        case 1: // LWES_U_INT_16_TOKEN
            read(&tmpUint16)
            event.attributes[attrName] = tmpUint16
        case 2: // LWES_INT_16_TOKEN
            read(&tmpInt16)
            event.attributes[attrName] = tmpInt16
        case 3: // LWES_U_INT_32_TOKEN
            read(&tmpUint32)
            event.attributes[attrName] = tmpUint32
        case 4: // LWES_INT_32_TOKEN
            read(&tmpInt32)
            event.attributes[attrName] = tmpInt32
        case 5: // LWES_STRING_TOKEN
            read(&tmpUint16)
            event.attributes[attrName] = string(p.Next(int(tmpUint16)))
        case 6: // LWES_IP_ADDR_TOKEN
            tmpIp := p.Next(4)
            // not sure if this is completely accurate
            event.attributes[attrName] = net.IPv4(tmpIp[3], tmpIp[2], tmpIp[1], tmpIp[0])
        case 7: // LWES_INT_64_TOKEN
            read(&tmpInt64)
            event.attributes[attrName] = tmpInt64
        case 8: // LWES_U_INT_64_TOKEN
            read(&tmpUint64)
            event.attributes[attrName] = tmpUint64
        case 9: // LWES_BOOLEAN_TOKEN
            event.attributes[attrName] = 1 == p.Next(1)[0]
        }
    }
}

// PrettyString returns a "pretty" formatted string.
func (e *Event) PrettyString() string {
    var buf bytes.Buffer

    buf.WriteString(e.Name)
    buf.WriteString("\n")

    for key := range e.attributes {
        buf.WriteString(key)
        buf.WriteString(": ")
        // gah
        buf.WriteString(fmt.Sprintln(e.attributes[key]))
    }

    return buf.String()
}

// MarshalJSON returns a json byte array of name:attributes
// net.IP is base64 encoded
func (e *Event) MarshalJSON() (data []byte, err error) {
    m := make(eventAttrs)
    m[e.Name] = e.attributes
    return json.Marshal(m)
}
