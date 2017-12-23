package main

// lchart reads data from the standard input file, or specified files, expecting a tab-separated list of text,data pairs
// where text is an arbitrary string, and data is intepreted as a floating point value.
// A line beginning with "#" is parsed as a title, with the title text beginning after the "#".
//
// For example:
//
//	# PDF File Sizes
//	casino.pdf	410907
//	countdown.pdf	157784
//	deck-12x8.pdf	837831
//	deck-dejavu.pdf	1601595
//	deck-fira-4x3.pdf	1196167
//	deck-fira.pdf	1195517
//	deck-gg.pdf	978688
//	deck-gofont.pdf	1044627
//
//
// The command line options are:
//
//	-dmim		zero minimum (default false)
//	-bar		show bars (default true)
//	-connect	connect data points (default false)
//	-dot		show dot plot (default false)
//	-grid		show gridlines on the y axis (default false)
//	-val		show values (default true)
//	-vol		show volume plot (default false)
//	-top		top of the plot (default 80)
//	-bottom 	bottom of the plot (default 30)
//	-left		left margin (default 10)
//	-right		right margin (default 90)
//	-barwidth	barwidth (default computed from the number of data points)
//	-ls		linespacing (default 2.4)
//	-textsize	text size (default 1.5)
//	-xlabel		x axis label interval (default 1)
//	-yaxis		show a y axis
//	-color 		data color (default "lightsteelblue")
//	-datafmt	data format (default "%.1f")
