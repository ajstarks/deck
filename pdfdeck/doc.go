/*

pdfdeck is a program for making PDF slides using the deck package.

Usage

	$ go get github.com/ajstarks/deck/pdfdeck
	$ pdfdeck deck.xml  # make deck.pdf in current directory

the -g percent option draws a grid scaled to the specifed percentage on each slide
the -f option specifies the location of the font directory
the -sans, -serif, and -mono options specify fonts
the -outdir option specifies the directory where PDF files are written; defaults to the current directory

*/
package main
