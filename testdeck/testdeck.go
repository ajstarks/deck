// testdeck: dump a slide deck 
package main

import (
	"os"
	"github.com/ajstarks/deck"
)

// for every file, dump a deck
func main() {
	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			d, err := deck.Read(f, 1024, 768)
			if err != nil {
				continue
			}
			deck.Dump(d)
		}
	}
}
