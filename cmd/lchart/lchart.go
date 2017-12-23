// lchart - make line charts in the deck format
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck/generate"
)

// LineData defines the name,value pairs
type LineData struct {
	label string
	value float64
}

var (
	ts, left, right, top, bottom, ls, barw                                      float64
	xint                                                                        int
	showdot, datamin, showvolume, showbar, showval, connect, showaxis, showgrid bool
	datacolor, datafmt                                                          string
)

const (
	bgcolor      = "white"
	titlecolor   = "black"
	labelcolor   = "rgb(75,75,75)"
	dotlinecolor = "lightgray"
	valuecolor   = "rgb(127,0,0)"
)

// vmap maps one range into another
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// getdata name,value pairs, with optional comments,
// returning a slice with the data, allong with min, max and title
func getdata(r io.ReadCloser) ([]LineData, float64, float64, string) {
	var (
		data []LineData
		d    LineData
		err  error
	)

	maxval := -1.0
	minval := 1e308
	title := ""
	scanner := bufio.NewScanner(r)
	// read a line, parse into name, value pairs
	// compute min and max values
	for scanner.Scan() {
		t := scanner.Text()
		if t[0] == '#' && len(t) > 2 {
			title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) != 2 {
			continue
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

// dottedvline makes dotted vertical line
func dottedvline(deck *generate.Deck, x, y1, y2, dotsize, step float64, color string) {
	for y := y1; y <= y2; y += step {
		deck.Circle(x, y, dotsize, color)
	}
}

// axislabel constructs the axis label
func axislabel(deck *generate.Deck, x, y, min, max float64) {
	yf := vmap(y, min, max, bottom, top)
	deck.TextEnd(x, yf, fmt.Sprintf("%0.f", y), "sans", ts*0.75, "black")
	if showgrid {
		deck.Line(left, yf, right, yf, 0.1, "lightgray")
	}
}

// yaxis draws the yaxis (labels with optional grid)
func yaxis(deck *generate.Deck, x, min, max, steps float64) {
	l := math.Log10(max)
	p := math.Pow10(int(l))
	div := p / steps
	var yp float64
	for yp = min; yp <= max; yp += div {
		axislabel(deck, x, yp, min, max)
	}
	axislabel(deck, x, yp, min, max)
}

// makeplot makes the plot using input from the reader
func makeplot(deck *generate.Deck, r io.ReadCloser) {
	linedata, mindata, maxdata, title := getdata(r)
	if !datamin {
		mindata = 0
	}
	l := len(linedata)
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

	// Begin the slide with a centered title (if specified)
	deck.StartSlide(bgcolor)
	linespacing := ts * ls
	if len(title) > 0 {
		deck.TextMid(left+((right-left)/2), top+(linespacing*1.5), title, "sans", ts*1.5, titlecolor)
	}

	if showaxis {
		yaxis(deck, left-ts*1.5, mindata, maxdata, 2.0)
	}

	// for every name, value pair, make the draw the chart elements
	var px, py float64
	for i, data := range linedata {
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
			deck.TextMid(x, y+ts, fmt.Sprintf(datafmt, data.value), "sans", ts*0.75, valuecolor)
		}
		if xint > 0 && i%xint == 0 {
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
	deck.EndSlide()
}

func main() {
	// command line parameters
	flag.Float64Var(&ts, "textsize", 1.5, "text size")
	flag.Float64Var(&left, "left", 10.0, "left margin")
	flag.Float64Var(&right, "right", 100-left, "right margin")
	flag.Float64Var(&top, "top", 80.0, "top of the plot")
	flag.Float64Var(&bottom, "bottom", 30.0, "bottom of the plot")
	flag.Float64Var(&ls, "ls", 2.4, "ls")
	flag.Float64Var(&barw, "barwidth", 0, "barwidth")

	flag.BoolVar(&showbar, "bar", true, "show bar")
	flag.BoolVar(&showdot, "dot", false, "show dot")
	flag.BoolVar(&showvolume, "vol", false, "show volume")
	flag.BoolVar(&datamin, "dmin", false, "zero minimum")
	flag.BoolVar(&showval, "val", true, "show values")
	flag.BoolVar(&showaxis, "yaxis", true, "show y axis")
	flag.BoolVar(&showgrid, "grid", false, "show grid")
	flag.BoolVar(&connect, "connect", false, "connected line plot")

	flag.IntVar(&xint, "xlabel", 1, "x axis label interval")

	flag.StringVar(&datacolor, "color", "lightsteelblue", "data color")
	flag.StringVar(&datafmt, "datafmt", "%.1f", "data format")
	flag.Parse()

	// start the deck, for every file name make a slide.
	// if no files, read from standard input.
	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	if len(flag.Args()) > 0 {
		for _, file := range flag.Args() {
			r, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			makeplot(deck, r)
		}
	} else {
		makeplot(deck, os.Stdin)
	}
	deck.EndDeck()
}
