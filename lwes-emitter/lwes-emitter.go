package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bhenderson/lwes-go"
)

var (
	addr    = flag.String("address", "224.2.2.22:12345", "Listen Channel")
	emitter *lwes.Emitter
)

func init() {
	flag.Usage = usage

}

func main() {
	flag.Parse()

	var err error
	emitter, err = lwes.NewEmitter(*addr)

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
	fmt.Fprintf(os.Stderr, "Usage: json_input | %s [opts]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Print("json_input optionally can contain Name.\n")
	os.Exit(1)
}
