package main

import (
    "net"
    "fmt"
    "os"
    "github.com/bhenderson/lwes"
)

func main() {
    var h, t int8
    var ifc *net.Interface
    emitter, err := lwes.NewEmitter("224.2.2.22", 12345, h, t, ifc)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    event := lwes.NewEvent()
    event.Name = "Event4"
    event.SetAttribute("field1", int16(3))

    err = emitter.Emit(event)

    if err != nil {
        fmt.Println(err)
    }
}
