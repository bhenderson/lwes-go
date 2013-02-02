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

var addr string
var stdin bool
var emitter *lwes.Emitter

func init() {
    flag.Usage = usage

    flag.StringVar(&addr,   "udpaddr", "224.2.2.22:12345", "Listen Channel")
    flag.BoolVar(  &stdin,  "stdin",   false,        "Emit an event for each line of json on stdin ({name:...attributes:...})")
}

func main() {
    flag.Parse()

    var err error
    emitter, err = lwes.NewEmitter(addr)
    emitter.Heartbeat = 0

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    event := lwes.NewEvent()

    if stdin {
        fromJson(event)
    } else {
        // Just emit a default event.
        event.Name = "LWES::TestEmitter"
        event.SetAttribute("field1", 15)
        emit(event)
    }
}

func emit(e *lwes.Event) {
    err := emitter.Emit(e)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func fromJson(e *lwes.Event) {
    dec := json.NewDecoder(os.Stdin)

    for {

        if err := dec.Decode(e); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }

        emit(e)
    }
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage: %s [opts]\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}
