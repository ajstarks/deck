/*
 lchart reads data from the standard input file, or specified files, expecting a tab-separated list of text,data pairs
 where text is an arbitrary string, and data is intepreted as a floating point value.
 A line beginning with "#" is parsed as a title, with the title text beginning after the "#".

 For example:

	# PDF File Sizes
	casino.pdf	410907
	countdown.pdf	157784
	deck-12x8.pdf	837831
	deck-dejavu.pdf	1601595
	deck-fira-4x3.pdf	1196167
	deck-fira.pdf	1195517
	deck-gg.pdf	978688
	deck-gofont.pdf	1044627


 The command line options are:

   -bar
     	show bar (default true)
   -barwidth float
     	barwidth
   -bottom float
     	bottom of the plot (default 30)
   -color string
     	data color (default "lightsteelblue")
   -connect
     	connected line plot
   -datafmt string
     	data format (default "%.1f")
   -dmin
     	zero minimum
   -dot
     	show dot
   -grid
     	show grid
   -layout string
     	chart orientation (h=horizontal, v=vertical) (default "v")
   -left float
     	left margin (default 10)
   -ls float
     	ls (default 2.4)
   -max float
     	maximum (default -1)
   -min float
     	minimum (default -1)
   -right float
     	right margin (default 90)
   -textsize float
     	text size (default 1.5)
   -top float
     	top of the plot (default 80)
   -val
     	show values (default true)
   -vol
     	show volume
   -xlabel int
     	x axis label interval (default 1)
   -yaxis
     	show y axis (default true)
*/
package main
