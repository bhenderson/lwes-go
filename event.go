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
    name string
    attributes eventAttrs
}

// NewEvent returns an initialized Event
func NewEvent() *Event {
    return &Event{attributes: make(eventAttrs)}
}

// Name returns the name or class of an event. This is separate from an attribute
func (e *Event) Name() string {
    return e.name
}

// Iterator interface
func (e *Event) Iterator() eventAttrs {
    return e.attributes
}

// Get an attribute
func (e *Event) Get(s string) interface{} {
    return e.attributes[s]
}

func (event *Event) ToBytes() ([]byte, error) {
    buf := new(bytes.Buffer)
    var err error

    // TODO how do these functions effect memory?
    write := func(d interface{}) bool {
        err = binary.Write(buf, binary.BigEndian, d)
        return err == nil
    }

    writeAttr := func(i int, d interface{}) bool {
        return write(byte(i)) && write(d)
    }

    // write attribute name, attribute type (as an int), attribute value
    writePair := func(s string, i int, d interface{}) bool {
        return writeAttr( len(s), []byte(s) ) && writeAttr(i, d)
    }

    if ! (
        // length of event name
        write( byte(len(event.Name()))       ) &&
        // event name
        write( []byte(event.Name())          ) &&
        // num attributes
        write( uint16(len(event.attributes)) ) ) {
            return nil, err
    }

    for key := range event.attributes {
        switch v := event.attributes[key].(type) {
        default:
            // unknown type. skip it.
        case uint16:
            writePair( key, 1, v )
        case int16:
            writePair( key, 2, v )
        case uint32:
            writePair( key, 3, v )
        case int32:
            writePair( key, 4, v )
        case string:
            if writePair( key, 5, uint16(len(v)) ) {
                write( []byte(v) )
            }
        case net.IP:
            val := v[len(v) - 4:]
            writePair( key, 6, []byte{val[3],val[2],val[1],val[0]} )
        case uint64:
            writePair( key, 7, v )
        case int64:
            writePair( key, 8, v )
        case int:
            writePair( key, 8, int64(v) )
        case bool:
            var b byte
            if v { b = 1 } else { b = 0 }
            writePair( key, 9, b )
        }

        if err != nil { return nil, err }
    }

    return buf.Bytes(), nil
}

func (event *Event) FromBytes(buf []byte) {
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
    event.name = string(p.Next(int(nameSize)))

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

    buf.WriteString(e.Name())
    buf.WriteString("\n")

    for key := range e.attributes {
        buf.WriteString(key)
        buf.WriteString(": ")
        buf.WriteString(fmt.Sprintln(e.attributes[key]))
    }

    return buf.String()
}

// MarshalJSON returns a json byte array of name:attributes
// net.IP is base64 encoded
func (e *Event) MarshalJSON() (data []byte, err error) {
    m := make(eventAttrs)
    m[e.Name()] = e.attributes
    return json.Marshal(m)
}
