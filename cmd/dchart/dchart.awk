# dchart -- make various charts using deck markup
BEGIN {
	# set operating parameters and defaults
	if (fulldeck == "") {
		fulldeck="t"
	}
	if (datafmt == "") {
		datafmt="%.2f"
	}
	if (xlabel == 0) {
		xlabel=1
	}
	if (line == "") {
		line = "t"
	}
	if (bar == "") {
		bar = "f"
	}
	if (scatter == "") {
		scatter = "f"
	}
	if (dotop == 0) {
		dotop=50
	}
	if (volop == 0) {
		volop=50
	}
	if (top == 0) {
		top=80
	}
	if (bottom == 0) {
		bottom=20
	}
	if (right == 0) {
		right=90
	}
	if (left == 0) {
		left=10
	}
	if (linesize == 0) {
		linesize = 0.1
	}
	if (barwidth == 0) {
		barwidth=0.1
	}
	if (dotsize == 0) {
		dotsize=linesize*10
	}
	if (color == "") {
		color = "lightsteelblue"
	}
	if (val == "") {
		val = "t"
	}
	if (ygrid == "") {
		ygrid = "f"
	}
	if (valcolor == "") {
		valcolor="maroon"
	}
	if (labelsize == 0) {
		labelsize = 1.5
	}
	if (valsize == 0) {
		valsize = labelsize*0.75
	}
	if (labelcolor == "") {
		labelcolor = "gray"
	}
	lx=left-2
	width=right-left
}

NR == 1 { 
	max = $2
	# begin the deck
	if (fulldeck == "t") {
		print "deck"
	}
}
{ 
	# save label,data pairs, determine maxima
	label[NR]=$1
	data[NR]=$2
	if ($2 > max) {
		max = $2
	}
}

# Map one range to another
function vmap(value, low1, high1, low2, high2) {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

END {
	# variables to be used in decksh markup
	min=0
	interval=width/(NR-1)
	printf "\tcolor=\"%s\"\n",      color
	printf "\tvalcolor=\"%s\"\n",   valcolor
	printf "\tvalsize=%g\n",        valsize
	printf "\tlabelsize=%g\n",      labelsize
	printf "\tlabelcolor=\"%s\"\n", labelcolor
	printf "\tlinesize=%g\n",       linesize
	printf "\tdotsize=%g\n",        dotsize
	printf "\tbarwidth=%g\n",       barwidth
	printf "\tleft=%g\n",           left
	printf "\tright=%g\n",          right
	printf "\ttop=%g\n",            top
	printf "\tbottom=%g\n",         bottom
	printf "\tdotop=%g\n",          dotop
	printf "\tvolop=%g\n",          volop

	# begin the slide
	if (fulldeck == "t") {
		printf "\tslide\n"
	}

	# y Axis
	if (yaxis == "t") {
		if (yrange == "") {
			ymin=min
			ymax=max
			yint=(ymax-ymin)/5
		} else {
			split(yrange, yr, ",")
			ymin=yr[1]
			yint=yr[2]
			ymax=yr[3]
		}
		for (yl=ymin; yl <= ymax; yl+=yint) {
			y = vmap(yl, min, max, bottom, top)
			printf "\t\tetext \"%s\" %g %g labelsize \"sans\" labelcolor\n", sprintf(datafmt, yl), left-2, y-(labelsize/3)
			if (ygrid == "t") {
				printf "\t\tline left %g right %g 0.05 \"gray\"\n", y, y 
			}
		}
	}

	if (chartitle != "") {
		printf "text \"%s\" left %g %g\n", chartitle, top+5, labelsize*1.5
	} 

	# for every label,data line, draw the elements of line, bar, scatter, and volume charts (each may be individually specified)
	x=left
	for (i=1; i <= NR; i++) {
		y=vmap(data[i], min, max, bottom, top)
		if (line == "t" && i < NR) {
			printf "\t\tline %g %g %g %g linesize color\n", x, y, x+interval, vmap(data[i+1], min, max, bottom, top)
		}
		if (bar == "t") {
			printf "\t\tline %g %g %g %g barwidth color\n", x, bottom, x, y
		}
		if (scatter == "t") {
			printf "\t\tcircle %g %g dotsize color dotop\n", x, y
		}
		if (val == "t") {
			printf "\t\tctext \"%s\" %g %g valsize \"serif\" valcolor\n", sprintf(datafmt, data[i]), x, y+1
		}
		# x axis labels
		if (i%xlabel == 0 || i == 1) {
			printf "\t\tctext \"%s\" %g %g labelsize \"sans\" labelcolor\n", label[i], x, bottom-(labelsize*2)
		}
		x+=interval
	}
	# prepare data for volume (area) chart
	if (vol == "t") {
		xa[1] = left
		ya[1] = bottom
		xv=left
		n=2
		for (i=1; i <= NR; i++) {
			xa[n] = xv
			ya[n] = vmap(data[i], min, max, bottom, top)
			xv += interval
			n++
		}
		xa[n] = right
		ya[n] = bottom

		n++
		xa[n] = xa[1]
		ya[n] = ya[1]

		lp=length(xa)
		printf "\t\tpolygon \""
		for (i=1; i <= lp; i++) {
			printf " %g", xa[i]
		}
		printf "\" \""
		for (i=1; i <= lp; i++) {
			printf (" %g", ya[i])
		}
		printf "\" color volop\n"
	}

	# close out slide and deck
	if (fulldeck == "t") {
		print "\teslide"
		print "edeck"
	}
}