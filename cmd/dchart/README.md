# dchart - charts for deck

```dchart``` generates deck markup for bar, line, scatter, dot, volume, donut and proportional charts, reading data from the standard input or specified files. 
Unless specified otherwise, each input source generates a slide in the deck.

The input data format a tab-separated or CSV formatted list of ```label,data``` pairs where label is an arbitrary string, 
and data is intepreted as a floating point value. 

A line beginning with "#" is parsed as a title, 
with the title text beginning after the "#".  If a third column is present, it serves as an annotation.


Here is an example input data file:

	# GOOG Stock Volume (Millions of Shares)
	2017-01-01	33.1916
	2017-02-01	25.6825
	2017-03-01	33.8351	Peak
	2017-04-01	25.1619
	2017-05-01	32.1801

Example CSV file:
	
	#,GOOG Stock Volume (Millions of Shares)
	2017-01-01,33.1916
	2017-02-01,25.6825
	2017-03-01,33.8351,Peak
	2017-04-01,25.1619
	2017-05-01,32.1801

Typically ```dchart``` generates input for deck clients like ```pdfdeck```, or ```pdi``` (shell script for pdfdeck which reads
deck markup on the standard input and produces PDF on the standard output).

	$ dchart foo.d bar.d baz.d > fbb.xml && pdfdeck fbb.xml && open fbb.pdf
	$ dchart -min=0 -max=700 -datafmt %0.2f -line -bar=f -vol -dot [A-Z]*.d | pdi > allvol.pdf
	$ ls -lS | awk 'BEGIN {print "# File Size"} NR > 1 {print $NF "\t" $5}' | dchart -hbar | pdi > fs.pdf

## Defaults

With no options, ```dchart``` makes a bar graph, showing data values and every data label.

## Placement

The plot is positioned and scaled on the deck canvas with the 
```-top```, ```-bottom```, ```-left```, and ```-right```, ```-x```, and ```-y``` flags. 
These flag values represent percentages on the deck canvas.

## Chart types and elements

The ```-bar```, ```-hbar```, ```-line```, ```-dot```, ```-scatter```, ```-vol```, 
```-pgrid```, ```-pmap```, and ```-donut``` 
flags specify the chart types.

The ```-grid```, ```-title```, ```-val```, and ```-yaxis``` 
flags control the visibility of plot components. 


## Command line options

	-dmim        data minimum (default false, min=0)
	-min         set the minimum value
	-max         set the maximum value
	-csv         read CSV files (default false)
	-csvcol      specify the columns to use for label,value

	-bar         show bar chart (default true)
	-wbar        show "word" bar chart (default false)
	-hbar        horizontal chart layout (default false)
	-scatter     show a scatter chart (default false)
	-dot         show dot plot (default false)
	-line        show line chart (default false)
	-vol         show volume plot (default false)
	-pgrid       show a proportional grid (default false)
	-pmap        show proportional map (default false)
	-donut       show a donut chart (default false)

	-grid        show gridlines on the y axis (default false)
	-val         show values (default true)
	-valpos      value position (t=top, b=bottom, m=middle) (default "t")
	-yaxis       show a y axis (default true)
	-yrange      specify the y axis labels (min,max,step)
	-fulldeck    generate full deck markup (default true)
	-title       show title (default true)
	-chartitle   specify the title (overiding title in the data)
	
	-top         top of the plot (default 80)
	-bottom      bottom of the plot (default 30)
	-left        left margin (default 20)
	-right       right margin (default 80)
	-x           x location of the donut chart (default 50)
	-y           y location of the donut chart (default 50)
	
	-psize       diameter of the donut (default 30)
	-pwidth      width of the donut or proportional map (default 3)
	-barwidth    barwidth (default computed from the number of data points)
	-ls          linespacing (default 2.4)
	-textsize    text size (default 1.5)
	-xlabel      x axis label interval (default 1, 0 to supress all labels)
	-xlast       show the last x label
	-color       data color (default "lightsteelblue")
	-vcolor      value color (default "rgb(127,0,0)")
	-datafmt     data format for values (default "%.1f")


## Examples

Using this data in ```AAPL.d```

	# AAPL Volume
	2017-01-01	563.122
	2017-02-01	574.969
	2017-03-01	561.628
	2017-04-01	373.304
	2017-05-01	653.755
	2017-06-01	684.178
	2017-07-01	421.992
	2017-08-01	661.069
	2017-09-01	679.879
	2017-10-01	504.291
	2017-11-01	600.663
	2017-12-01	417.354

here are some variations.

	$ dchart AAPL.d

![no-args](images/no-args.png)

	$ dchart -yrange=0,700,50 AAPL.d

![yrange](images/yrange.png)

	$ dchart -xlabel=2 -left 30 -right 70 -top 70 -bottom 40 -yaxis=f AAPL.d

![pos](images/pos.png)

	$ dchart -color gray AAPL.d

![bar-gray](images/bar-gray.png)

	$ dchart -grid AAPL.d # add a y axis grid

![bar-grid](images/bar-grid.png)

	$ dchart -grid -barwidth=1 AAPL.d # adjust the bar width

![barwidth](images/barwidth.png)

	$ dchart -bar=f -dot AAPL.d # no bars, dot plot

![dot](images/dot.png)

	$ dchart -bar=f -vol AAPL.d # no bars, volume plot

![vol](images/vol.png)

	$ dchart -datafmt %0.2f -bar=f -dot -line AAPL.d

![dot-line](images/dot-connect.png)

	$ dchart -bar=f -line AAPL.d # line chart

![connect](images/connect.png)

	$ dchart -bar=f -line -yaxis=f -val=f AAPL.d # only show line and x axis

![connect-no-axis-no-val](images/connect-no-axis-no-val.png)


	$ dchart -scatter -val=f -bar=f -yaxis=f AAPL.d

![scatter](images/scatter.png)


	$ dchart -bar=f -line -vol -dot AAPL.d # combine line, volume, and dot

![vol-dot](images/vol-dot.png)


	$ dchart -bar=f -line -vol -dot -yaxis=f AAPL.d # as above, removing the y-axis

![vol-dot-no-axis](images/vol-dot-no-axis.png)

	$ dchart -bar=f -line -vol -dot -grid AAPL.d

![connect-dot-vol-val-grid](images/connect-dot-vol-val-grid.png)

	$ dchart -hbar AAPL.d

![hlayout](images/hlayout.png)

	$ sort -k2 -nr pdf.d | dchart -left 20 -hbar

![sorted-bar](images/sorted-hbar.png)

Using this data in ``browser.d``

	# Browser Market Share Dec 2016-Dec 2017
	Chrome	53.72
	Safari	14.47
	Other	9.36
	UC	8.28
	Firefox	6.23
	IE	3.99
	Opera	3.9

here are views of proportional data:
	
	$ dchart -wbar browser.d
	
![wbar](images/wbar.png)
	
	$ dchart -donut -color=std -pwidth=5 browser.d 
	
![donut](images/donut.png)

	$ dchart -pmap -pwidth=5 -textsize=1 browser.d

![pmap](images/pmap.png)

Using this data in incar.d:

	# US Incarceration Rate
	White	39	white
	Hispanic	19	burlywood
	Black	40	black
	Other	2	gray

the note field may be used to specify the color

	$ dchart -x 10 -y 80 -ls 3 -bgcolor steelblue -pgrid  incar.d

![pgrid](images/pgrid.png)

