/*
svgdeck is a program for making SVG slides using the deck package.

Usage

	$ go get github.com/ajstarks/deck/svgdeck
	$ svgdeck deck.xml  # make deck-nnn.svg in current directory

One SVG file per slide is generated in the output directory.  For example, a deck input file named "deck.xml"
with 5 slides would generate deck-000.svg through deck-004.svg.

Clicking on any content will navigate to the next slide, cycling to the first slide when the last slide is reached.

the -grid percent option draws a grid scaled to the specifed percentage on each slide.

the -sans, -serif, and -mono options specify fonts.

the -outdir option specifies the directory where SVG files are written; defaults to the current directory.

the -title options adds title metadata.

the -pagesize option specifies the page dimensions (wxh or Letter, Legal, A3, A4, A5, ArchA, 4R, Index, Widescreen).

the -stdout option specified that output goes to the standard output file.

*/
package main
