# lchart - linecharts in the deck format

lchart reads data from the standard input file, or specified files, expecting a tab-separated list 
of text,data pairs where text is an arbitrary string, and data is intepreted as a floating point value.
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


Typical use is the generating data from a deck client like pdfdeck,

    $ lchart foo.d bar.d baz.d > fbb.xml && pdfdeck fbb.xml && open fbb.pdf
    $ ls -lS | awk 'BEGIN { print "# File Size"} NR > 1 {print $NF "\t" $5}' | lchart | pdi > fs.pdf

The plot it positioned and scaled on the deck canvas with the ```-top```, ```-bottom```, ```-left```, and ```-right``` flags.

The  ```-bar```, ```-connect```, ```-dot```, ```-grid```, ```-val```, ```-vol```, and ```-yaxis``` 
flags toggle the visibility of plot components.  With no options, the plot is a bar graph with yaxis labels,
showing data values, and every data label.


The command line options are:

	-dmim     data minimum (default false, min=0)
	-bar      show bars (default true)
	-connect  connect data points (default false)
	-dot      show dot plot (default false)
	-grid     show gridlines on the y axis (default false)
	-val      show values (default true)
	-vol      show volume plot (default false)
	-yaxis    show a y axis (default false)
	
	
	-top      top of the plot (default 80)
	-bottom   bottom of the plot (default 30)
	-left     left margin (default 20)
	-right    right margin (default 80)
	
	
	-barwidth barwidth (default computed from the number of data points)
	-ls       linespacing (default 2.4)
	-textsize text size (default 1.5)
	-xlabel   x axis label interval (default 1)
	-color    data color (default "lightsteelblue")
	-datafmt  data format for values (default "%.1f")
