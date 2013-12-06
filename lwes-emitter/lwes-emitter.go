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
)

func init() {
    flag.Usage = usage

    flag.StringVar(&addr, "address", "224.2.2.22:12345", "Listen Channel")
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

    dec := json.NewDecoder(os.Stdin)

    for {

        if err := dec.Decode(e); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal("json error: ", err)
        }

        emit(e)
    }
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage: %s [opts]\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}
