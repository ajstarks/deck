/*
Package deck is a library for clients to make scalable presentations, using a standard markup language.
Clients read deck files into the Deck structure, and traverse the structure for display, publication, etc.
Clients may be interactive or produce standard formats such as SVG or PDF.

Elements

Here are the elements of a deck:

	deck: enclosing element
	canvas: describe the dimensions of the drawing canvas, one per deck
	slide: within a deck, any number of slides, specify the slide background and text colors.

within slides an number of:

	text: plain, textblock, or code
	list: plain, bullet, number
	image: JPEG or PNG images
	line: straight line
	rect: rectangle
	ellipse: ellipse
	curve: Quadratic Bezier curve
	arc: elliptical arc

Markup

Here is a sample deck in XML:

	<deck>
	     <canvas width="1024" height="768"/>
	     <slide bg="maroon" fg="white">
		<image xp="20" yp="30" width="256" height="256" name="picture.png"/>
		<text xp="20" yp="80" sp="3">Deck uses these elements</text>
		<line xp1="20" yp1="75" xp2="90" yp2="75" sp="0.3" color="rgb(127,0,0)"/>
		<list xp="20" yp="70" sp="1.5">
			<li>canvas<li>
			<li>slide</li>
			<li>text</li>
			<li>list</li>
			<li>line</li>
			<li>rect</li>
			<li>ellipse</li>
			<li>curve</li>
			<li>arc</li>
		</list>
		<line    xp1="20" yp1="10" xp2="30" yp2="10"/>
		<rect    xp="35"  yp="10" wp="4" hp="3" color="rgb(127,0,0)"/>
		<ellipse xp="45"  yp="10" wp="4" hr="75" color="rgb(0,127,0)"/>
		<arc     xp="55"  yp="10" wp="4" hr="75" a1="0" a2="180" color="rgb(0,0,127)"/>
	     </slide>
	</deck>

The list, text, rect, and ellipse elements have common attributes:

	xp: horizontal percentage
	yp: vertical percentage
	sp: font size percentage
	type: "bullet", "number" (list), "block", "code" (text)
	align: "left", "middle", "end"
	opacity: 0.0-1.0 (fully transparent - opaque)
	color: SVG names ("maroon"), or RGB "rgb(127,0,0)"
	font: "sans", "serif", "mono"

Layout

All layout in done in terms of percentages, using a coordinate system with the origin (0%, 0%) at the lower left.
The x (horizontal) direction increases to the right, and the y (vertical) direction increasing to upwards.
For example to place an element in the middle of the canvas, specify xp="50" yp="50". To place an element
one-third from the top, and one-third from the bottom: xp="66.6" yp="33.3".

The size of text is also scaled to the width of the canvas. For example sp="3" is a typical size for slide headings.
The sizes of graphical elements (width, height, stroke width) are also scaled to the canvas width.

The content of the slides are automatically scaled based on the specified canvas size
(sane defaults are should be set by clients, if dimensions are not specified).


Example


	package main

	import (
		"fmt"
		"log"

		"github.com/ajstarks/deck"
	)

	func dotext(x, y, size float64, text deck.Text) {
		fmt.Println("\tText:", x, y, size, text.Tdata)
	}

	func dolist(x, y, size float64, list deck.List) {
		fmt.Println("\tList:", x, y, size)
		for i, l := range list.Li {
			fmt.Println("\t\titem", i, l)
		}
	}
	func main() {
		presentation, err := deck.Read("deck.xml", 1024, 768) // open the deck
		if err != nil {
			log.Fatal(err)
		}
		for slidenumber, slide := range presentation.Slide { // for every slide...
			fmt.Println("Processing slide", slidenumber)
			for _, t := range slide.Text { // process the text elements
				x, y, size := deck.Dimen(presentation.Canvas, t.Xp, t.Yp, t.Sp)
				dotext(x, y, size, t)
			}
			for _, l := range slide.List { // process the list elements
				x, y, size := deck.Dimen(presentation.Canvas, l.Xp, l.Yp, l.Sp)
				dolist(x, y, size, l)
			}
		}
	}

*/
package deck
