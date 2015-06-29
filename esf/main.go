package main

import "os"

func main() {
	scanner := NewScanner(os.Stdin)
	esf, _ := scanner.Scan()
	os.Stdout.Write(esf)
}
