/*
dchart generates deck markup for bar, line, dot, and volume charts, reading data from the standard input or specified files.
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
	-csv         read CSV files (default false)
	-csvcol      specify the columns to use for label,value

	-bar         show bars (default true)
	-hbar        horizontal chart layout (default false)
	-scatter     show scatter chart (default false)
	-wbar        show "word" bar chart (default false)
	-line        show line chart (default false)
	-dot         show dot plot (default false)
	-grid        show gridlines on the y axis (default false)
	-val         show values (default true)
	-rline       show a regression line (default false)
	-frame       show a frame outlining the chart (default false)
	-pct         show percentages with values (default false)
	-valpos      value position (t=top, b=bottom, m=middle) (default "t")
	-vol         show volume plot (default false)
	-pgrid       show a proportional grid (default false)
	-pmap        show proportional map (default false)
	-donut       show a donut chart (default false)
	-radial      show a radial chart (default false)
	-spokes      show spokes on the radial chart (default false)
	-yaxis       show a y axis (default true)
	-yrange      define the y axis range (min,max,step)
	-fulldeck    generate full markup (default true)
	-title       show title (default true)
	-chartitle   specify the title (overiding title in the data)
	-hline       horizontal line with optional label (value,label)
	-noteloc     note location (c-center, r-right, l-left, default c)


	-top         top of the plot (default 80)
	-bottom      bottom of the plot (default 30)
	-left        left margin (default 20)
	-right       right margin (default 80)
	-x           x location of the donut chart (default 50)
	-y           y location of the donut chart (default 50)
	-psize       diameter of the donut (default 30)

	-pwidth      width of the donut or proportional map (default 3)
	-barwidth    barwidth (default computed from the number of data points)
	-linewidth   linewidth for line charts (default 0.2)
	-ls          linespacing (default 2.4)
	-textsize    text size (default 1.5)
	-xlabel      x axis label interval (default 1, 0 to supress all labels)
	-xlast       show the last x label
	-color       data color (default "lightsteelblue")
	-vcolor      value color (default "rgb(127,0,0)")
	-rlcolor     regression line color (default "rgb(127,0,0)")
	-framecolor  frame color (default "rgb(127,0,0)")
	-volop       volume opacity (default 50)
	-datafmt     data format for values (default "%.1f")
	-note        show annotation (default true)
*/
package main
