/*
Package generate provides a high-level API for the creation of slide decks using
the structures of the deck package (github.com/ajstarks/deck).

Initialization of the package specifies the io.Writer destination for the generated markup,
and the width and height of the slides's canvas. (Speciying (0,0) allows the client to use default dimensions).

Each deck element (text, list, image, rect, ellipse, line, curve, arc, and polygon) are supported.
Slides use a percentage-based coordinate system (origin at the lower left corner,
x increasing left to right, 0-100%, y increasing to the upwards, 0-100%).

By default slides use black text on a white background.
Elements may have colors and opacities applied to them.

Example:

	package main

	import (
		"os"
		"github.com/ajstarks/deck/generate"
	)
	
	func main() {
		deck := generate.NewSlides(os.Stdout, 0, 0)
		deck.StartDeck() // start the deck
	
		// Text alignment
		deck.StartSlide("rgb(180,180,180)")         // New slide with a gray background
		deck.Text(50, 80, "left", 10, "black")      // left-aligned black text
		deck.TextMid(50, 50, "center", 10, "gray")  // centered gray text
		deck.TextEnd(50, 20, "right", 10, "white")  // right-aligned white text
		deck.Line(50, 100, 50, 0, 0.2, "black", 20) // vertical line
		deck.EndSlide() // end the slide

		// List
		items := []string{"First", "Second", "Third", "Fourth", "Fifth"}
		deck.StartSlide()                            // start a new slide
		deck.Text(10, 90, "Imporant Items", 5, "")   // title for the list
		deck.List(10, 80, 4, items, "bullet", "red") // make a bullet list
		deck.EndSlide()                              // end the slide
	
		deck.EndDeck() // end the deck
	}
*/
package generate
