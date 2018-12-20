// dchart - make charts in the deck format
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
	ts, left, right, top, bottom, ls, barw, umin, umax, psize, pwidth, volop, linewidth   float64
	xint, pmlen                                                                           int
	readcsv, showdot, datamin, showvolume, showscatter, showpct                           bool
	showbar, showval, showxlast, showline, showhbar, wbar, showaxis, shownote             bool
	showgrid, showtitle, fulldeck, showdonut, showpmap, showpgrid, showradial, showspokes bool
	bgcolor, datacolor, datafmt, chartitle, valpos, valuecolor, yaxr, csvcols, hline      string
)

var blue7 = []string{
	"rgb(8,69,148)",
	"rgb(33,113,181)",
	"rgb(66,146,198)",
	"rgb(107,174,214)",
	"rgb(158,202,225)",
	"rgb(198,219,239)",
	"rgb(239,243,255)",
}

var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

const (
	titlecolor   = "black"
	labelcolor   = "rgb(75,75,75)"
	dotlinecolor = "lightgray"
	defaultfmt   = "%.1f"
	wbop         = 30.0
	largest      = math.MaxFloat64
	smallest     = -math.MaxFloat64
	topclock     = math.Pi / 2
	fullcircle   = math.Pi * 2
	transparency = 50.0
)

func cmdflags() {
	// command line options
	flag.Float64Var(&ts, "textsize", 1.5, "text size")
	flag.Float64Var(&left, "left", -1.0, "left margin") // default set to out of bounds because different charts need individual defaults
	flag.Float64Var(&right, "right", 90.0, "right margin")
	flag.Float64Var(&top, "top", 80.0, "top of the plot")
	flag.Float64Var(&bottom, "bottom", 30.0, "bottom of the plot")
	flag.Float64Var(&ls, "ls", 2.4, "ls")
	flag.Float64Var(&barw, "barwidth", 0, "barwidth")
	flag.Float64Var(&umin, "min", -1, "minimum")
	flag.Float64Var(&umax, "max", -1, "maximum")
	flag.Float64Var(&psize, "psize", 40.0, "size of the donut")
	flag.Float64Var(&pwidth, "pwidth", ts*3, "width of the pmap/donut/radial")
	flag.Float64Var(&linewidth, "linewidth", 0.2, "width of line for line charts")
	flag.Float64Var(&volop, "volop", 50, "volume opacity")

	flag.BoolVar(&showbar, "bar", true, "show a bar chart")
	flag.BoolVar(&showdot, "dot", false, "show a dot chart")
	flag.BoolVar(&showvolume, "vol", false, "show a volume chart")
	flag.BoolVar(&showdonut, "donut", false, "show a donut chart")
	flag.BoolVar(&showpmap, "pmap", false, "show a proportional map")
	flag.BoolVar(&showline, "line", false, "show a line chart")
	flag.BoolVar(&showhbar, "hbar", false, "show a horizontal bar chart")
	flag.BoolVar(&showval, "val", true, "show data values")
	flag.BoolVar(&showaxis, "yaxis", false, "show y axis")
	flag.BoolVar(&showtitle, "title", true, "show title")
	flag.BoolVar(&showgrid, "grid", false, "show y axis grid")
	flag.BoolVar(&showscatter, "scatter", false, "show scatter chart")
	flag.BoolVar(&showradial, "radial", false, "show a radial chart")
	flag.BoolVar(&showspokes, "spokes", false, "show spokes on radial charts")
	flag.BoolVar(&showpgrid, "pgrid", false, "show proportional grid")
	flag.BoolVar(&shownote, "note", true, "show annotations")
	flag.BoolVar(&showxlast, "xlast", false, "show the last label")
	flag.BoolVar(&fulldeck, "fulldeck", true, "generate full markup")
	flag.BoolVar(&datamin, "dmin", false, "zero minimum")
	flag.BoolVar(&readcsv, "csv", false, "read CSV data")
	flag.BoolVar(&wbar, "wbar", false, "show word bar chart")
	flag.BoolVar(&showpct, "pct", false, "show computed percentages with values")
	flag.IntVar(&xint, "xlabel", 1, "x axis label interval (show every n labels, 0 to show no labels)")
	flag.IntVar(&pmlen, "pmlen", 20, "pmap label length")

	flag.StringVar(&chartitle, "chartitle", "", "specify the title (overiding title in the data)")
	flag.StringVar(&csvcols, "csvcol", "", "label,value from the CSV header")
	flag.StringVar(&valpos, "valpos", "t", "value position (t=top, b=bottom, m=middle)")
	flag.StringVar(&datacolor, "color", "lightsteelblue", "data color")
	flag.StringVar(&valuecolor, "vcolor", "rgb(127,0,0)", "value color")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&datafmt, "datafmt", defaultfmt, "data format")
	flag.StringVar(&yaxr, "yrange", "", "y-axis range (min,max,step)")
	flag.StringVar(&hline, "hline", "", "horizontal line value,label")

	flag.Parse()
}

// xmlesc escapes XML
func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

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
			d.note = xmlesc(fields[2])
		} else {
			d.note = ""
		}
		if n == 1 && len(csvcols) > 0 { // column header is assumed to be the first row
			li, vi = getheader(fields, csvcols)
			title = fields[vi]
			continue
		}

		d.label = xmlesc(fields[li])
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
	return data, minval, maxval, xmlesc(title)
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
			d.note = xmlesc(fields[2])
		} else {
			d.note = ""
		}
		d.label = xmlesc(fields[0])
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
	return data, minval, maxval, xmlesc(title)
}

// dottedvline makes dotted vertical line, using circles,
// with specified step
func dottedvline(deck *generate.Deck, x, y1, y2, dotsize, step float64, color string) {

	if y1 < y2 { // positive
		for y := y1; y <= y2; y += step {
			deck.Circle(x, y, dotsize, color)
		}
	} else { // negative
		for y := y2; y <= y1; y += step {
			deck.Circle(x, y, dotsize, color)
		}
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

	if datafmt != defaultfmt {
		return fmt.Sprintf(datafmt, x)
	}

	frac := x - float64(int(x))
	if frac == 0 {
		return fmt.Sprintf("%0.f", x)
	}
	return fmt.Sprintf(datafmt, x)
}

// datasum computes the sum of the chart data
func datasum(data []ChartData) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	return sum
}

// pct computs the percentage of a range of values
func pct(data []ChartData) []float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}

	p := make([]float64, len(data))
	for i, d := range data {
		p[i] = (d.value / sum) * 100
	}
	return p
}

// pgrid makes a proportional grid with the specified rows and columns
func pgrid(deck *generate.Deck, data []ChartData, title string, rows, cols int) {
	// sanity checks

	if left < 0 {
		left = 30.0
	}

	if rows*cols != 100 {
		return
	}

	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	pct := make([]float64, len(data))
	for i, d := range data {
		pct[i] = math.Floor((d.value / sum) * 100)
	}

	// encode the data in a string vector
	chars := make([]string, 100)
	cb := 0
	for k := 0; k < len(data); k++ {
		for l := 0; l < int(pct[k]); l++ {
			chars[cb] = data[k].note
			cb++
		}
	}

	// make rows and cols
	n := 0
	y := top
	for i := 0; i < rows; i++ {
		x := left
		for j := 0; j < cols; j++ {
			if n >= 100 {
				break
			}
			deck.Circle(x, y, ts, chars[n])
			n++
			x += ls
		}
		y -= ls
	}

	// title and legend
	if len(title) > 0 && showtitle {
		deck.Text(left-ts/2, top+ts*2, title, "sans", ts*1.5, titlecolor)
	}
	cx := (float64(cols-1) * ls) + ls/2
	for i, d := range data {
		y -= ls * 1.2
		deck.Circle(left, y, ts, d.note)
		deck.Text(left+ts, y-(ts/2), d.label+" ("+dformat(pct[i])+"%)", "sans", ts, "black")
		if showval {
			deck.TextEnd(left+cx, y-(ts/2), dformat(d.value), "sans", ts, valuecolor)
		}
	}
}

// polar converts polar to Cartesian coordinates
func polar(x, y, r, t float64) (float64, float64) {
	px := x + r*math.Cos(t)
	py := y + r*math.Sin(t)
	return px, py
}

// spokes draws the points and lines like spokes on a wheel
func spokes(deck *generate.Deck, cx, cy, r, spokesize float64, n int, color string) {
	t := topclock
	step := fullcircle / float64(n)
	for i := 0; i < n; i++ {
		px, py := polar(cx, cy, r, t)
		deck.Line(cx, cy, px, py, spokesize, "lightgray")
		deck.Circle(px, py, 0.5, color)
		t -= step
	}
}

// radial draws a radial plot
func radial(deck *generate.Deck, data []ChartData, title string, maxd float64) {

	if left < 0 {
		left = 50.0
	}

	dx := left
	dy := top
	if len(title) > 0 && showtitle {
		deck.TextMid(dx, dy, title, "sans", ts*1.5, titlecolor)
	}
	if umax > 0 {
		maxd = umax
	}
	t := topclock
	deck.Circle(dx, dy, pwidth*2, "silver", 10)
	step := fullcircle / float64(len(data))
	var color string
	for _, d := range data {
		cv := vmap(d.value, 0, maxd, 2, psize)
		px, py := polar(dx, dy, pwidth, t)
		tx, ty := polar(dx, dy, pwidth+(psize/2)+(ts*2), t)

		if len(d.note) > 0 {
			color = d.note
		} else {
			color = datacolor
		}

		deck.TextMid(tx, ty, d.label, "sans", ts/2, "black")
		if showval {
			deck.TextMid(px, py-ts/3, dformat(d.value), "mono", ts, valuecolor)
		}
		if showspokes {
			spokes(deck, px, py, psize/2, 0.05, int(d.value), color)
		} else {
			deck.Circle(px, py, cv, color, transparency)
			deck.Line(tx, ty, px, py, 0.05, "gray", 50)
		}
		t -= step
	}
}

// pmap draws a porpotional map
func pmap(deck *generate.Deck, data []ChartData, title string) {
	if left < 0 {
		left = 20.0
	}
	x := left
	pl := (right - left)
	bl := pl / 100.0
	hspace := 0.10
	var ty float64
	var textcolor string
	if len(title) > 0 && showtitle {
		deck.TextMid(x+pl/2, top+(pwidth*2), title, "sans", ts*1.5, titlecolor)
	}
	for i, p := range pct(data) {
		bx := (p * bl)
		if p < 3 || len(data[i].label) > pmlen {
			ty = top - pwidth*1.5
			deck.Line(x+(bx/2), ty+(ts*1.5), x+(bx/2), top, 0.1, dotlinecolor)
		} else {
			ty = top
		}
		linecolor, lineop := stdcolor(i, data[i].note, datacolor, p)
		deck.Line(x, top, bx+x, top, pwidth, linecolor, lineop)
		if lineop == 100 {
			textcolor = "white"
		} else {
			textcolor = "black"
		}

		deck.TextMid(x+(bx/2), ty+(pwidth), data[i].label, "sans", ts*0.75, textcolor)
		if showval {
			deck.TextMid(x+(bx/2), ty-pwidth, dformat(data[i].value), "mono", ts/2, valuecolor)
		}
		deck.TextMid(x+(bx/2), ty-(ts/2), fmt.Sprintf(datafmt+"%%", p), "sans", ts, textcolor)

		x += bx - hspace
	}
}

// stdcolor uses either the standard color (cycling through a list) or specified color and opacity
func stdcolor(i int, dcolor, color string, op float64) (string, float64) {
	if color == "std" {
		return blue7[i%len(blue7)], 100
	}
	if len(dcolor) > 0 {
		return dcolor, 40
	}
	return color, op
}

// donut makes a donut chart
func donut(deck *generate.Deck, data []ChartData, title string) {
	if left < 0 {
		left = 50.0
	}
	a1 := 0.0
	dx := left // + (psize / 2)
	dy := top - (psize / 2)
	if len(title) > 0 && showtitle {
		deck.TextMid(dx, dy+(psize*1.2), title, "sans", ts*1.5, titlecolor)
	}
	for i, p := range pct(data) {
		angle := (p / 100) * 360.0
		a2 := a1 + angle
		mid := (a1 + a2) / 2

		bcolor, op := stdcolor(i, data[i].note, datacolor, p)
		deck.Arc(dx, dy, psize, psize, pwidth, a1, a2, bcolor, op)
		tx, ty := polar(dx, dy, psize*.85, mid*(math.Pi/180))
		deck.TextMid(tx, ty, fmt.Sprintf("%s "+datafmt+"%%", data[i].label, p), "sans", ts, "black")
		if showval {
			deck.TextMid(tx, ty-ts*1.5, fmt.Sprintf(dformat(data[i].value)), "sans", ts, valuecolor)
		}
		a1 = a2
	}
}

// pchart draws proportional data, either a pmap, pgrid, radial or donut using input from a Reader
func pchart(deck *generate.Deck, r io.ReadCloser) {
	data, _, maxdata, title := getdata(r)
	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}
	if fulldeck {
		deck.StartSlide(bgcolor)
	}
	switch {
	case showdonut:
		donut(deck, data, title)
	case showpmap:
		pmap(deck, data, title)
	case showpgrid:
		pgrid(deck, data, title, 10, 10)
	case showradial:
		radial(deck, data, title, maxdata)
	}
	if fulldeck {
		deck.EndSlide()
	}
}

// wbchart makes a word bar chart
func wbchart(deck *generate.Deck, r io.ReadCloser) {
	if left < 0 {
		left = 20.0
	}
	hts := ts / 2
	mts := ts * 0.75
	linespacing := ts * ls

	bardata, mindata, maxdata, title := getdata(r)
	if !datamin {
		mindata = 0
	}
	if fulldeck {
		deck.StartSlide(bgcolor)
	}

	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	if len(title) > 0 && showtitle {
		deck.Text(left, top+(linespacing*1.5), title, "sans", ts*1.5, titlecolor)
	}

	var sum float64
	if showpct {
		sum = datasum(bardata)
	}

	// for every name, value pair, make the chart
	y := top
	for _, data := range bardata {
		deck.Text(left+hts, y, data.label, "sans", ts, labelcolor)
		bv := vmap(data.value, mindata, maxdata, left, right)
		deck.Line(left+hts, y+hts, bv, y+hts, ts*1.5, datacolor, wbop)
		if showval {
			if showpct {
				avgs := fmt.Sprintf(" ("+datafmt+"%%)", 100*(data.value/sum))
				deck.TextEnd(left, y+(hts/2), dformat(data.value)+avgs, "mono", mts, valuecolor)
			} else {
				deck.TextEnd(left, y+(hts/2), dformat(data.value), "mono", mts, valuecolor)
			}
		}
		y -= linespacing
	}
	if fulldeck {
		deck.EndSlide()
	}
}

// hchart makes horizontal bar charts using input from a Reader
func hchart(deck *generate.Deck, r io.ReadCloser) {
	hts := ts / 2
	mts := ts * 0.75
	linespacing := ts * ls

	bardata, mindata, maxdata, title := getdata(r)

	if left < 0 {
		left = 30.0
	}

	if !datamin {
		mindata = 0
	}
	if fulldeck {
		deck.StartSlide(bgcolor)
	}

	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	if len(title) > 0 && showtitle {
		deck.TextMid(50, top+(linespacing*1.5), title, "sans", ts*1.5, titlecolor)
	}

	var sum float64
	if showpct {
		sum = datasum(bardata)
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
		if showval {
			if showpct {
				avgs := fmt.Sprintf(" ("+datafmt+"%%)", 100*(data.value/sum))
				deck.Text(bv+hts, y+(hts/2), dformat(data.value)+avgs, "mono", mts, valuecolor)
			} else {
				deck.Text(bv+hts, y+(hts/2), dformat(data.value), "mono", mts, valuecolor)
			}
		}
		y -= linespacing
	}
	if fulldeck {
		deck.EndSlide()
	}
}

// vchart makes charts using input from a Reader
// the types of charts are bar (column), dot, line, and volume
func vchart(deck *generate.Deck, r io.ReadCloser) {
	chartdata, mindata, maxdata, title := getdata(r)

	if left < 0 {
		left = 10.0
	}

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
		xvol[l+1] = right
		yvol[l+1] = bottom
	}

	linespacing := ts * ls
	spacing := ts * 1.5

	if fulldeck {
		deck.StartSlide(bgcolor)
	}

	if len(chartitle) > 0 {
		title = xmlesc(chartitle)
	}

	if len(title) > 0 && showtitle {
		deck.TextMid(left+((right-left)/2), top+(linespacing*1.5), title, "sans", spacing, titlecolor)

	}

	if showaxis {
		yaxis(deck, left-spacing-(dw*0.5), mindata, maxdata)
	}

	if len(hline) > 0 {
		var hl float64
		var hs string
		fmt.Sscanf(hline, "%f,%s", &hl, &hs)
		hy := vmap(hl, mindata, maxdata, bottom, top)
		deck.Line(left, hy, right, hy, 0.1, valuecolor, 50)
		if len(hs) > 0 {
			deck.Text(right+ts/2, hy-ts/4, hs, "serif", ts*0.75, labelcolor)
		}
	}

	var sum float64
	if showpct {
		sum = datasum(chartdata)
	}

	// for every name, value pair, make the chart elements
	var px, py float64
	for i, data := range chartdata {
		x := vmap(float64(i), 0, dlen, left, right)
		y := vmap(data.value, mindata, maxdata, bottom, top)

		if showvolume {
			xvol[i+1] = x
			yvol[i+1] = y
		}
		if showline && i > 0 {
			deck.Line(px, py, x, y, linewidth, datacolor)
		}
		if showdot {
			dottedvline(deck, x, bottom, y, ts/6, 1, dotlinecolor)
			deck.Circle(x, y, ts*.6, datacolor)
		}
		if showscatter {
			deck.Circle(x, y, ts*.6, datacolor)
		}
		if showbar {
			deck.Line(x, bottom, x, y, dw, datacolor)
		}
		if showval {
			yv := y + ts
			switch valpos {
			case "t":
				if data.value < 0 {
					yv = y - ts
				} else {
					yv = y + ts
				}
			case "b":
				yv = bottom + ts
			case "m":
				yv = y - ((y - bottom) / 2)
			}
			if showpct {
				avgs := fmt.Sprintf(" ("+datafmt+"%%)", 100*(data.value/sum))
				deck.TextMid(x, yv, dformat(data.value)+avgs, "sans", ts*0.75, valuecolor)
			} else {
				deck.TextMid(x, yv, dformat(data.value), "sans", ts*0.75, valuecolor)
			}
		}
		if len(data.note) > 0 && shownote {
			deck.TextMid(x, y+0.1, data.note, "serif", ts*0.6, labelcolor)
		}
		// show x label every xinit times, show the last, if specified
		if xint > 0 && (i%xint == 0 || (showxlast && i == l-1)) {
			xlabels := strings.Split(data.label, `\n`)
			xly := bottom - (ts * 2)
			for _, xl := range xlabels {
				deck.TextMid(x, xly, xl, "sans", ts*0.8, labelcolor)
				xly -= ts * 1.2
			}
		}
		px = x
		py = y
	}
	if showvolume {
		deck.Polygon(xvol, yvol, datacolor, volop)
	}
	if fulldeck {
		deck.EndSlide()
	}
}

// chart makes charts according to the orientation:
// horizontal bar or line, bar, dot, or donut volume charts
func chart(deck *generate.Deck, r io.ReadCloser) {
	switch {
	case showhbar:
		hchart(deck, r)
	case wbar:
		wbchart(deck, r)
	case showdonut, showpmap, showpgrid, showradial:
		pchart(deck, r)
	default:
		vchart(deck, r)
	}
}

func main() {
	// process command line options,
	// start the deck, for every file name make a slide.
	// Read from standard input, if no files are specified.
	cmdflags()
	deck := generate.NewSlides(os.Stdout, 0, 0)
	if fulldeck {
		deck.StartDeck()
	}
	if len(flag.Args()) > 0 {
		for _, file := range flag.Args() {
			r, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			chart(deck, r)
		}
	} else {
		chart(deck, os.Stdin)
	}
	if fulldeck {
		deck.EndDeck()
	}
}
