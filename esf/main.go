// Command esf is a generator for ESF (Event Specification File) to go files.
package main

import "os"

func main() {
	scanner := NewScanner(os.Stdin)
	esf, _ := scanner.Scan()
	os.Stdout.Write(esf)
}
