package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/bhenderson/lwes-go"
)

var (
	addr    string
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

	go trapInt()
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

func trapInt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	for _ = range signalChan {
		fmt.Println("received interrupt")
		os.Exit(1)
	}
}
