/*

pdfdeck is a program for making PDF slides using the deck package.
The PDF is generated to stdout.

Usage

	$ go get github.com/ajstarks/deck/pdfdeck
	$ pdfdeck deck.xml > deck.pdf

the -g percent option draws a grid scaled to the specifed percentage on each slide
the -f option specifies the location of the font directory
the -sans, -serif, and -mono options specify fonts

*/
package main
