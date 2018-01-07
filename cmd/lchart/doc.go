/*
lchart generates deck markup for bar, line, dot, and volume charts, reading data from the standard input or specified files.
Unless specified otherwise, each input source generates a slide in the deck.

The input data format a tab-separated list of label,data pairs where label is an arbitrary string,
and data is intepreted as a floating point value. A line beginning with "#" is parsed as a title,
with the title text beginning after the "#". A third column specifies an annotation.

Here is an example input data file:

 	# GOOG Stock Volume (Millions of Shares)
 	2017-01-01	33.1916
 	2017-02-01	25.6825
 	2017-03-01	33.8351	Peak value
 	2017-04-01	25.1619
 	2017-05-01	32.1801

The command line options are:

	-dmim        data minimum (default false, min=0)
	-min         set the minimum value
	-max         set the maximum value

	-bar         show bars (default true)
	-hbar        horizontal chart layout (default false)
	-line        show line chart (default false)
	-dot         show dot plot (default false)
	-grid        show gridlines on the y axis (default false)
	-val         show values (default true)
	-valpos      value position (t=top, b=bottom, m=middle) (default "t")
	-vol         show volume plot (default false)
	-yaxis       show a y axis (default false)
	-standalone  only generate internal markup (default false)
	-title       show title (default true)
	-chartitle   specify the title (overiding title in the data)

	-top         top of the plot (default 80)
	-bottom      bottom of the plot (default 30)
	-left        left margin (default 20)
	-right       right margin (default 80)

	-barwidth    barwidth (default computed from the number of data points)
	-ls          linespacing (default 2.4)
	-textsize    text size (default 1.5)
	-xlabel      x axis label interval (default 1, 0 to supress all labels)
	-color       data color (default "lightsteelblue")
	-vcolor      value color (default "rgb(127,0,0)")
	-datafmt     data format for values (default "%.1f")
*/
package main
