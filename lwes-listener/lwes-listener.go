package main

import (
    "flag"
    "fmt"
    "encoding/json"
    "os"
    "github.com/bhenderson/lwes"
)

var addr string
var port int
var pretty bool
var printj bool

func init() {
    flag.Usage = usage

    flag.StringVar(&addr,   "address", "224.2.2.22", "Listen Address")
    flag.IntVar(   &port,   "port",    12345,       "Listen Port")
    flag.BoolVar(  &pretty, "pretty",  false,       "Pretty print event")
    flag.BoolVar(  &printj, "json",    false,       "Print event as json")
}

func main() {
    flag.Parse()

    listener, err := lwes.NewListener(addr, port)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    listener.Each(callback)
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage: %s [opts]\n", os.Args[0])
    flag.PrintDefaults()
    os.Exit(1)
}

func callback(event *lwes.Event, err error) {
    if err != nil {
        fmt.Println(err)
        return
    }

    switch {
    default:
        fmt.Println(event)
    case pretty:
        fmt.Println(event.PrettyString())
    case printj:
        b, _ := json.Marshal(event)
        fmt.Println(string(b))
    }
    os.Stdout.Sync()
}
