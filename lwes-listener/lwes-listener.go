package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bhenderson/lwes-go"
)

var (
	addr   string
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

	lc := listen(listener)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	for {
		select {
		case <-sigc:
			listener.Close()
			return
		case event := <-lc:
			callback(event)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [opts]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func listen(l *lwes.Listener) chan *lwes.Event {
	ch := make(chan *lwes.Event)

	go func() {
		for {
			e, err := l.Recv()
			if err == nil {
				ch <- e
			}
		}
	}()

	return ch
}

func callback(event *lwes.Event) error {
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
