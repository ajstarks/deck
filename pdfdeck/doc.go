/*
pdfdeck is a program for making PDF slides using the deck package.

Usage

	$ go get github.com/ajstarks/deck/pdfdeck
	$ pdfdeck deck.xml  # make deck.pdf in current directory

the -grid percent option draws a grid scaled to the specifed percentage on each slide.

the -fontdir option specifies the location of the font directory.

the -sans, -serif, and -mono options specify fonts.

the -outdir option specifies the directory where PDF files are written; defaults to the current directory.

the -author option adds author metadata.

the -title options adds title metadata.

the -pagesize option specifies the page dimensions (Letter, Legal, Tabloid, A2, A3, A4, A5, ArchA, Index, 4R, Widescreen).

the -pagewidth and pageheight options specify custom page dimensions.

the -stdout option specified that output goes to the standard output file.
*/
package main
