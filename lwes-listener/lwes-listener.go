package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/bhenderson/lwes-go"
    "os"
)

var (
    addr string
    pretty bool
    printj bool
)

func init() {
    flag.Usage = usage

    flag.StringVar(&addr, "address", "224.2.2.22:12345", "Listen Channel")
    flag.BoolVar(&pretty, "pretty", false, "Pretty print event")
    flag.BoolVar(&printj, "json", false, "Print event as json")
}

func main() {
    flag.Parse()

    listener, err := lwes.NewListener(addr)

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

func callback(event *lwes.Event, err error) error {
    if err != nil {
        fmt.Println(err)
        return nil
    }

    switch {
    default:
        fmt.Println(event)
    case printj:
        var b []byte
        if pretty {
            b, _ = json.MarshalIndent(event, "", "  ")
        } else {
            b, _ = json.Marshal(event)
        }
        fmt.Println(string(b))
    case pretty:
        fmt.Println(event.Name)
        for k, _ := range event.Attributes {
            fmt.Printf("%s: %v\n", k, event.Attributes[k])
        }
        fmt.Println("")
    }
    os.Stdout.Sync()

    return nil
}
