// lchart - make charts in the deck format
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck/generate"
)

// ChartData defines the name,value pairs
type ChartData struct {
	label string
	value float64
	note  string
}

var (
	ts, left, right, top, bottom, ls, barw, umin, umax                                                                           float64
	xint                                                                                                                         int
	readcsv, showdot, datamin, showvolume, showbar, showval, showxlast, connect, hbar, showaxis, showgrid, showtitle, fullmarkup bool
	bgcolor, datacolor, datafmt, chartitle, valpos, valuecolor, yaxr, csvcols                                                    string
)

const (
	titlecolor   = "black"
	labelcolor   = "rgb(75,75,75)"
	dotlinecolor = "lightgray"
	largest      = math.MaxFloat64
	smallest     = -math.MaxFloat64
)

// vmap maps one range into another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// getheader returns the indicies of the comma-separated list of fields
// by default or on error, return 0, 1
// For example given this header:
// First,Second,Third,Sum
// First,Sum returns 0,3 and First,Third returns 0,2
func getheader(s []string, lv string) (int, int) {
	li := 0
	vi := 1
	cv := strings.Split(lv, ",")
	if len(cv) != 2 {
		return li, vi
	}
	for i, p := range s {
		if p == cv[0] {
			li = i
		}
		if p == cv[1] {
			vi = i
		}
	}
	return li, vi
}

// getdata reads imput from a Reader, either tab-separated or CSV
func getdata(r io.ReadCloser) ([]ChartData, float64, float64, string) {
	var min, max float64
	var title string
	var data []ChartData
	if readcsv {
		data, min, max, title = csvdata(r)
	} else {
		data, min, max, title = tsvdata(r)
	}
	return data, min, max, title
}

// csvdata reads CSV structured name,value pairs, with optional comments,
// returning a slice with the data, allong with min, max and title
func csvdata(r io.ReadCloser) ([]ChartData, float64, float64, string) {
	var (
		data []ChartData
		d    ChartData
		err  error
	)
	input := csv.NewReader(r)
	maxval := smallest
	minval := largest
	title := ""
	n := 0
	li := 0
	vi := 1
	for {
		n++
		fields, csverr := input.Read()
		if csverr == io.EOF {
			break
		}
		if csverr != nil {
			fmt.Fprintf(os.Stderr, "%v %v\n", csverr, fields)
			continue
		}

		if len(fields) < 2 {
			continue
		}
		if fields[0] == "#" {
			title = fields[1]
			continue
		}
		if len(fields) == 3 {
			d.note = fields[2]
		} else {
			d.note = ""
		}
		if n == 1 && len(csvcols) > 0 { // column header is assumed to be the first row
			li, vi = getheader(fields, csvcols)
			title = fields[vi]
			continue
		}

		d.label = fields[li]
		d.value, err = strconv.ParseFloat(fields[vi], 64)
		if err != nil {
			d.value = 0
		}
		if d.value > maxval {
			maxval = d.value
		}
		if d.value < minval {
			minval = d.value
		}
		data = append(data, d)
	}
	r.Close()
	return data, minval, maxval, title
}

// tsvdata reads tab-delimited name,value pairs, with optional comments,
// returning a slice with the data, allong with min, max and title
func tsvdata(r io.ReadCloser) ([]ChartData, float64, float64, string) {
	var (
		data []ChartData
		d    ChartData
		err  error
	)

	maxval := smallest
	minval := largest
	title := ""
	scanner := bufio.NewScanner(r)
	// read a line, parse into name, value pairs
	// compute min and max values
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // process titles
			title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 2 {
			continue
		}
		if len(fields) == 3 {
			d.note = fields[2]
		} else {
			d.note = ""
		}
		d.label = fields[0]
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		if d.value > maxval {
			maxval = d.value
		}
		if d.value < minval {
			minval = d.value
		}
		data = append(data, d)
	}
	r.Close()
	return data, minval, maxval, title
}

// dottedvline makes dotted vertical line, using circles,
// with specified step
func dottedvline(deck *generate.Deck, x, y1, y2, dotsize, step float64, color string) {
	for y := y1; y <= y2; y += step {
		deck.Circle(x, y, dotsize, color)
	}
}

// dottedhline makes a dotted horizontal line, using circles,
// with specified step and separation
func dottedhline(d *generate.Deck, x, y, width, height, step, space float64, color string) {
	for xp := x; xp < x+width; xp += step {
		d.Circle(xp, y, height, color)
		xp += space
	}
}

// yrange parses the min, max, step for axis labels
func yrange(s string) (float64, float64, float64) {
	var min, max, step float64
	n, err := fmt.Sscanf(s, "%f,%f,%f", &min, &max, &step)
	if n != 3 || err != nil {
		return 0, 0, 0
	}
	return min, max, step
}

// cyrange computes "optimal" min, max, step for axis labels
// rounding the max to the appropriate number, given the number of labels
func cyrange(min, max float64, n int) (float64, float64, float64) {
	l := math.Log10(max)
	p := math.Pow10(int(l))
	pl := math.Ceil(max / p)
	ymax := pl * p
	return min, ymax, ymax / float64(n)
}

// yaxis constructs y axis labels
func yaxis(deck *generate.Deck, x, dmin, dmax float64) {
	var axismin, axismax, step float64
	if yaxr == "" {
		axismin, axismax, step = cyrange(dmin, dmax, 5)
	} else {
		axismin, axismax, step = yrange(yaxr)
	}
	if step <= 0 {
		return
	}
	for y := axismin; y <= axismax; y += step {
		yp := vmap(y, dmin, dmax, bottom, top)
		deck.TextEnd(x, yp, fmt.Sprintf("%0.f", y), "sans", ts*0.75, "black")
		if showgrid {
			deck.Line(left, yp, right, yp, 0.1, "lightgray")
		}
	}
}

// dformat returns the string representation of a float64
// according to the datafmt flag value.
// if there is no fractional portion of the float64, override the flag and
// return the string with no decimals.
func dformat(x float64) string {
	frac := x - float64(int(x))
	if frac == 0 {
		return fmt.Sprintf("%0.f", x)
	}
	return fmt.Sprintf(datafmt, x)
}

// hbar makes horizontal bar charts using input from a Reader
func hchart(deck *generate.Deck, r io.ReadCloser) {
	hts := ts / 2
	mts := ts * 0.75
	linespacing := ts * ls

	bardata, mindata, maxdata, title := getdata(r)
	if !datamin {
		mindata = 0
	}
	deck.StartSlide(bgcolor)

	if len(chartitle) > 0 {
		title = chartitle
	}

	if len(title) > 0 && showtitle {
		deck.TextMid(50, top+(linespacing*1.5), title, "sans", ts*1.5, titlecolor)
	}

	// for every name, value pair, make the chart
	y := top
	for _, data := range bardata {
		deck.TextEnd(left-hts, y, data.label, "sans", ts, labelcolor)
		bv := vmap(data.value, mindata, maxdata, left, right)
		if showdot {
			dottedhline(deck, left, y+hts, bv-left, ts/5, 1, 0.25, dotlinecolor)
			deck.Circle(bv, y+hts, mts, datacolor)
		} else {
			deck.Line(left, y+hts, bv, y+hts, ts, datacolor)
		}
		deck.Text(bv+hts, y+(hts/2), dformat(data.value), "mono", mts, valuecolor)
		y -= linespacing
	}
	deck.EndSlide()
}

// vchart makes charts using input from a Reader
// the types of charts are bar (column), dot, line, and volume
func vchart(deck *generate.Deck, r io.ReadCloser) {
	chartdata, mindata, maxdata, title := getdata(r)
	if !datamin {
		mindata = 0
	}

	if umin >= 0 {
		mindata = umin
	}

	if umax >= 0 && umax > mindata {
		maxdata = umax
	}

	l := len(chartdata)
	dlen := float64(l - 1)

	// define the width of bars
	var dw = (right-left)/dlen - 1
	if barw > 0 && barw <= dw {
		dw = barw
	}

	// for volume plots, allocate, fill in the extrema
	var xvol, yvol []float64
	if showvolume {
		xvol = make([]float64, l+2)
		yvol = make([]float64, l+2)
		xvol[0] = left
		yvol[0] = bottom
		xvol[l+1] = left
		yvol[l+1] = bottom
	}

	linespacing := ts * ls
	spacing := ts * 1.5

	if fullmarkup {
		deck.StartSlide(bgcolor)
	}

	if len(chartitle) > 0 {
		title = chartitle
	}

	if len(title) > 0 && showtitle {
		deck.TextMid(left+((right-left)/2), top+(linespacing*1.5), title, "sans", spacing, titlecolor)
	}

	if showaxis {
		yaxis(deck, left-spacing-(dw*0.5), mindata, maxdata)
	}

	// for every name, value pair, make the chart elements
	var px, py float64
	for i, data := range chartdata {
		x := vmap(float64(i), 0, dlen, left, right)
		y := vmap(data.value, mindata, maxdata, bottom, top)

		if showvolume {
			xvol = append(xvol, x)
			yvol = append(yvol, y)
		}
		if connect && i > 0 {
			deck.Line(px, py, x, y, 0.2, datacolor)
		}
		if showdot {
			dottedvline(deck, x, bottom, y, ts/6, 1, dotlinecolor)
			deck.Circle(x, y, ts*.6, datacolor)
		}
		if showbar {
			deck.Line(x, bottom, x, y, dw, datacolor)
		}
		if showval {
			yv := y + ts
			switch valpos {
			case "t":
				yv = y + ts
			case "b":
				yv = bottom + ts
			case "m":
				yv = y - ((y - bottom) / 2)
			}
			deck.TextMid(x, yv, dformat(data.value), "sans", ts*0.75, valuecolor)
		}
		if len(data.note) > 0 {
			deck.TextMid(x, y, data.note, "serif", ts*0.6, labelcolor)
		}
		// show x label every xinit times, show the last, if specified
		if xint > 0 && (i%xint == 0 || (showxlast && i == l-1)) {
			deck.TextMid(x, bottom-(ts*2), data.label, "sans", ts*0.8, labelcolor)
		}
		px = x
		py = y
	}
	if showvolume {
		xvol = append(xvol, right)
		yvol = append(yvol, bottom)
		deck.Polygon(xvol, yvol, datacolor, 50)
	}
	if fullmarkup {
		deck.EndSlide()
	}
}

// chart makes charts according to the orientation:
// horizontal bar or line, bar, dot or volume charts
func chart(deck *generate.Deck, r io.ReadCloser) {
	if hbar {
		hchart(deck, r)
	} else {
		vchart(deck, r)
	}
}

func main() {
	// command line options
	flag.Float64Var(&ts, "textsize", 1.5, "text size")
	flag.Float64Var(&left, "left", 10.0, "left margin")
	flag.Float64Var(&right, "right", 100-left, "right margin")
	flag.Float64Var(&top, "top", 80.0, "top of the plot")
	flag.Float64Var(&bottom, "bottom", 30.0, "bottom of the plot")
	flag.Float64Var(&ls, "ls", 2.4, "ls")
	flag.Float64Var(&barw, "barwidth", 0, "barwidth")
	flag.Float64Var(&umin, "min", -1, "minimum")
	flag.Float64Var(&umax, "max", -1, "maximum")

	flag.BoolVar(&showbar, "bar", true, "show a bar chart")
	flag.BoolVar(&showdot, "dot", false, "show a dot chart")
	flag.BoolVar(&showvolume, "vol", false, "show a volume chart")
	flag.BoolVar(&connect, "line", false, "show a line chart")
	flag.BoolVar(&datamin, "dmin", false, "zero minimum")
	flag.BoolVar(&hbar, "hbar", false, "horizontal bar")
	flag.BoolVar(&showval, "val", true, "show values")
	flag.BoolVar(&showaxis, "yaxis", true, "show y axis")
	flag.BoolVar(&showtitle, "title", true, "show title")
	flag.BoolVar(&showgrid, "grid", false, "show grid")
	flag.BoolVar(&fullmarkup, "standalone", true, "generate full markup")
	flag.BoolVar(&showxlast, "xlast", false, "show the last label")
	flag.BoolVar(&readcsv, "csv", false, "read CSV data")
	flag.IntVar(&xint, "xlabel", 1, "x axis label interval (show every n labels, 0 to show no labels)")

	flag.StringVar(&chartitle, "chartitle", "", "specify the title (overiding title in the data)")
	flag.StringVar(&csvcols, "csvcol", "", "label,value from the CSV header")
	flag.StringVar(&valpos, "valpos", "t", "value position (t=top, b=bottom, m=middle)")
	flag.StringVar(&datacolor, "color", "lightsteelblue", "data color")
	flag.StringVar(&valuecolor, "vcolor", "rgb(127,0,0)", "value color")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&datafmt, "datafmt", "%.1f", "data format")
	flag.StringVar(&yaxr, "yrange", "", "y-axis range (min,max,step)")
	flag.Parse()

	// start the deck, for every file name make a slide.
	// if no files, read from standard input.
	deck := generate.NewSlides(os.Stdout, 0, 0)
	if fullmarkup {
		deck.StartDeck()
	}
	if len(flag.Args()) > 0 {
		for _, file := range flag.Args() {
			r, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			chart(deck, r)
		}
	} else {
		chart(deck, os.Stdin)
	}
	if fullmarkup {
		deck.EndDeck()
	}
}
