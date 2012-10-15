package main

import (
    "fmt"
    "os"
    "github.com/bhenderson/lwes"
    "time"
)

func main() {
    emitter, err := lwes.NewEmitter("224.2.2.22", 12345)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    event := lwes.NewEvent("Event4")
    event.SetAttribute("field1", 15)

    for _ = range time.Tick(2 * time.Second) {
        err = emitter.Emit(event)

        if err != nil {
            fmt.Println(err)
        }
    }
}
