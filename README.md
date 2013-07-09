#deck: A Go package for slide decks

Package deck provides an interface, via a standard markup language for making scalable, percentage-based layout slide decks.
Clients read deck files into the Deck structure, and traverse the structure for display, publication, etc.
From a single markup language, clients may be interactive or produce standard formats such as HTML or PDF.

## Elements ##

* deck: enclosing element 
* canvas: describe the dimensions of the drawing canvas, one per deck
* slide: within a deck, any number of slides, specify the slide background and text colors.

within slides an number of:
* text: plain, textblock, or code
* list: plain, bullet, number
* image: JPEG or PNG images

## Markup ##

Here is a sample deck in XML:

	<deck>
	   <canvas width="1024" height="768"/>
	      <slide bg="maroon" fg="white">
		<image xp="20" yp="30" width="256" height="256" name="picture.png"/>
	         <text xp="20" yp="80" sp="3">Deck uses these elements</text>
	         <list xp="20" yp="70" sp="1.5">
	            <li>canvas<li>
	            <li>slide</li>
	            <li>text</li>
	            <li>list</li>
	            <li>image</li>
	         </list>
	      </slide>
	</deck>


The list and text elements have common attributes:

	xp: horizontal percentage
	yp: vertical percentage
	sp: font size percentage
	type: "bullet", "number" (list), "block", "code" (text)
	align: "left", "middle", "end"
	color: SVG names ("maroon"), or RGB "rgb(127,0,0)"
	font: "sans", "serif", "mono"
	
## Layout ##

All layout in done in terms of percentages, using a coordinate system with the origin (0%, 0%) at the lower left.
The x (horizontal) direction increases to the right, with the y (vertical) direction increasing to upwards.
For example, to place an element in the middle of the canvas, specify xp="50" yp="50". To place an element
one-third from the top, and one-third from the bottom: xp="66.6" yp="33.3".

The size of text is also scaled to the width of the canvas. For example sp="3" is a typical size for slide headings.

The content of the slides are automatically scaled based on the specified canvas size 
(sane defaults are should be set the clients, if dimensions not specified)
