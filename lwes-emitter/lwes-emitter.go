package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/bhenderson/lwes-go"
    "io"
    "log"
    "os"
)

var (
    addr string
    emitter *lwes.Emitter
    eventName string
)

func init() {
    flag.Usage = usage

    flag.StringVar(&addr, "address", "224.2.2.22:12345", "Listen Channel")
    flag.StringVar(&eventName,   "event_name", "LWES_GO::TestEvent", "The name of the event")
}

func main() {
    flag.Parse()

    var err error
    emitter, err = lwes.NewEmitter(addr)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    emitter.Heartbeat = 0

    fromJson()
}

func emit(e *lwes.Event) {
    err := emitter.Emit(e)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func fromJson() {
    e := lwes.NewEvent()
    attrs := make(map[string]interface{})
    e.Name = eventName

    dec := json.NewDecoder(os.Stdin)

    for {

        if err := dec.Decode(&attrs); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }

        e.Attributes = attrs
        emit(e)
    }
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage: %s [opts]\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}
