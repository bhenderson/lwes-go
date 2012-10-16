package lwes

import (
    "testing"
    "net"
)

func TestEventDeserializer(t *testing.T) {
    /*
        Event3
        boolean2: false
        boolean1: true
        time_sec: int32, 1350013760
        remote_addr: ip_addr, 192.168.0.1
        time_usec: int32, 410856
        field1: String value
    */

    eventSlice := []byte{0x6, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x33, 0x0, 0x6, 0x8, 0x74,
                         0x69, 0x6d, 0x65, 0x5f, 0x73, 0x65, 0x63, 0x4, 0x50, 0x77, 0x93, 0x40,
                         0x9, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x63, 0x4, 0x0,
                         0x6, 0x44, 0xe8, 0xb, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f, 0x61,
                         0x64, 0x64, 0x72, 0x6, 0x1, 0x0, 0xa8, 0xc0, 0x6, 0x66, 0x69, 0x65,
                         0x6c, 0x64, 0x31, 0x5, 0x0, 0xc, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
                         0x20, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x8, 0x62, 0x6f, 0x6f, 0x6c, 0x65,
                         0x61, 0x6e, 0x31, 0x9, 0x1, 0x8, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61,
                         0x6e, 0x32, 0x9, 0x0}
    e := NewEvent()
    e.fromBytes(eventSlice)

    if e.Name != "Event3" {
        t.Fatalf("e.Name got %#v, want %#v", e.Name, "Event1")
    }

    assertField := func(field string, d interface{}) {
        if v := e.Get(field); v != d {
            t.Fatalf("%v is\n\t%#v,\nwant\n\t%#v", field, v, d)
        }
    }

    if v := len(e.attributes); v != 6 {
        t.Fatalf("attributes length should be %v, but was %v", 6, v)
    }

    assertField("boolean1", true)
    assertField("boolean2", false)
    assertField("time_sec", int32(1350013760))
    assertField("time_usec", int32(410856))
    assertField("field1", "String value")
    // assertField("remote_addr", []byte{192,168,0,1})

    v := e.Get("remote_addr").(net.IP)

    if v[12] != 192 ||
       v[13] != 168 ||
       v[14] != 0   ||
       v[15] != 1 {
           t.Fatalf("remote_addr expected to be 192.168.0.1, got %v", v)
    }
}

func TestEventSerializer(t *testing.T) {
    f1  := uint8(15)
    f2  := int8(15)
    f3  := uint16(15)
    f4  := int16(15)
    f5  := uint32(15)
    f6  := int32(15)
    f7  := uint64(15)
    f8  := int64(15)
    f9  := "string"
    f10 := net.ParseIP("1.1.1.1")
    f11 := true // bool
    f12 := 15 // int
    f13 := 'h' // rune
    f14 := byte(16)

    e := NewEvent()
    e.Name = "Event4"
    e.SetAttribute("f1",  f1)
    e.SetAttribute("f2",  f2)
    e.SetAttribute("f3",  f3)
    e.SetAttribute("f4",  f4)
    e.SetAttribute("f5",  f5)
    e.SetAttribute("f6",  f6)
    e.SetAttribute("f7",  f7)
    e.SetAttribute("f8",  f8)
    e.SetAttribute("f9",  f9)
    e.SetAttribute("f10", f10)
    e.SetAttribute("f11", f11)
    e.SetAttribute("f12", f12)
    e.SetAttribute("f13", f13)
    e.SetAttribute("f14", f14)

    e.SetAttribute("f1p",  &f1)
    e.SetAttribute("f2p",  &f2)
    e.SetAttribute("f3p",  &f3)
    e.SetAttribute("f4p",  &f4)
    e.SetAttribute("f5p",  &f5)
    e.SetAttribute("f6p",  &f6)
    e.SetAttribute("f7p",  &f7)
    e.SetAttribute("f8p",  &f8)
    e.SetAttribute("f9p",  &f9)
    e.SetAttribute("f10p", &f10)
    e.SetAttribute("f11p", &f11)
    e.SetAttribute("f12p", &f12)
    e.SetAttribute("f13p", &f13)
    e.SetAttribute("f14p", &f14)

    ev := NewEvent()
    b, err := e.toBytes()

    if err != nil {
        t.Fatal("toBytes: ", err)
    }

    ev.fromBytes(b)

    assertEqual := func(expected interface{}, actual interface{}, msg string) {
        if expected != actual {
            t.Fatalf("%v\nexpected (%T)%#v to be equal to (%T)%#v", msg, expected,expected,actual,actual)
        }
    }

    assertField := func(field string, d interface{}) {
        assertEqual(d, ev.Get(field), field)
    }

    assertIPField := func(field string, d interface{}) {
        ip := ev.Get(field)
        assertEqual(ip.(net.IP)[12], d.(net.IP)[12], field)
        assertEqual(ip.(net.IP)[13], d.(net.IP)[13], field)
        assertEqual(ip.(net.IP)[14], d.(net.IP)[14], field)
        assertEqual(ip.(net.IP)[15], d.(net.IP)[15], field)
    }

    assertEqual("Event4", ev.Name, "ev.Name")
    // assertEqual(len(ev.attributes), len(e.attributes), "attribute length")

    assertField("f1",  uint16(f1))
    assertField("f2",  int16(f2))
    assertField("f3",  f3)
    assertField("f4",  f4)
    assertField("f5",  f5)
    assertField("f6",  f6)
    assertField("f7",  f7)
    assertField("f8",  f8)
    assertField("f9",  f9)
    assertIPField("f10", f10)
    assertField("f11", f11)
    assertField("f12", int64(f12))
    assertField("f13", f13)
    assertField("f14", uint16(f14))

    // pointers can't become pointers
    assertField("f1p",  uint16(f1))
    assertField("f2p",  int16(f2))
    assertField("f3p",  f3)
    assertField("f4p",  f4)
    assertField("f5p",  f5)
    assertField("f6p",  f6)
    assertField("f7p",  f7)
    assertField("f8p",  f8)
    assertField("f9p",  f9)
    assertIPField("f10p", f10)
    assertField("f11p", f11)
    assertField("f12p", int64(f12))
    assertField("f13p", f13)
    assertField("f14p", uint16(f14))
}

func TestEventSerializerNameLength(t *testing.T) {
    e := NewEvent()
    name := "aaaaaaaa"
    for ;len(name) <= MAX_SHORT_STRING_SIZE; { name += "aaaaaaaa" } // long string
    e.Name = name
    _, err := e.toBytes()

    if err == nil {
        t.Fatalf("expected name length (%d) to err", len(name))
    }
}

func TestEventSerializerKeyLength(t *testing.T) {
    e := NewEvent("Event")
    key := "aaaaaaaa"
    for ;len(key) <= MAX_SHORT_STRING_SIZE; { key += "aaaaaaaa" } // long string
    // should SetAttribute check length?
    e.SetAttribute(key, "too long")

    _, err := e.toBytes()

    if err == nil {
        t.Fatalf("expected key length (%d) to err", len(key))
    }
}
