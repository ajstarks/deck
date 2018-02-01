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

var (
	ts, left, right, top, dx, dy, dsize, dwidth float64
	readcsv, showpmap, fullmarkup               bool
	datacolor, datafmt, chartitle, csvcols      string
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

const (
	largest  = math.MaxFloat64
	smallest = -math.MaxFloat64
)

// ChartData defines the name,value pairs
type ChartData struct {
	label string
	value float64
	note  string
}


// doflags processes command line flags
func cmdflags() {
	flag.Float64Var(&ts, "textsize", 1.2, "text size")
	flag.Float64Var(&left, "left", 10.0, "left margin")
	flag.Float64Var(&right, "right", 100-left, "right margin")
	flag.Float64Var(&top, "top", 80.0, "top of the plot")
	flag.Float64Var(&dx, "x", 50.0, "x location")
	flag.Float64Var(&dy, "y", 50.0, "y location")
	flag.Float64Var(&dsize, "size", 20.0, "size of the donut")
	flag.Float64Var(&dwidth, "width", 3.0, "width of the donut")
	flag.BoolVar(&showpmap, "pmap", false, "show a pmap")
	flag.BoolVar(&readcsv, "csv", false, "read CSV data")
	flag.BoolVar(&fullmarkup, "standalone", true, "generate full markup")
	flag.StringVar(&csvcols, "csvcol", "", "label,value from the CSV header")
	flag.StringVar(&datafmt, "datafmt", "%.1f", "data format string")
	flag.StringVar(&datacolor, "color", "darkblue", "data color")
	flag.Parse()
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

// polar converts polar to Cartesian coordinates
func polar(x, y, r, t float64) (float64, float64) {
	px := x + r*math.Cos(t)
	py := y + r*math.Sin(t)
	return px, py
}

// pmap draws a porpotional map
func pmap(deck *generate.Deck, data []ChartData, begin, end, y float64, color, title string) {
	x := begin
	bl := (end - begin) / 100.0
	tsize := 1.2 // dsize / 4
	hspace := 0.10
	var ty float64
	if len(title) > 0 {
		deck.TextMid(x+(end-begin)/2, y+(dsize*1.2), title, "sans", ts*2, "black")
	}
	for i, p := range pct(data) {
		bx := (p * bl)
		if p < 4 || len(data[i].label) > 10 {
			ty = y + dsize*1.5
		} else {
			ty = y
		}
		deck.TextMid(x+(bx/2), ty, data[i].label, "sans", tsize, "black")
		deck.TextMid(x+(bx/2), ty-(tsize*1.5), fmt.Sprintf(datafmt+"%%", p), "mono", tsize, "black")
		deck.Line(x, y, bx+x, y, dsize, color, p)
		x += bx - hspace
	}
}

// donut makes a donut chart
func donut(deck *generate.Deck, data []ChartData, x, y, size, width float64, color, title string) {
	a1 := 0.0
	var bcolor string
	var op float64
	if len(title) > 0 {
		deck.TextMid(x, y+(size*1.2), title, "sans", ts*2, "black")
	}
	for i, p := range pct(data) {
		angle := (p / 100) * 360.0
		a2 := a1 + angle
		mid := (a1 + a2) / 2
		// use either the standard color (cycling through a list) or define color based in value
		if color == "std" {
			bcolor = blue7[i%len(blue7)]
			op = 100
		} else {
			bcolor = color
			op = p
		}
		deck.Arc(x, y, size, size, width, a1, a2, bcolor, op)
		tx, ty := polar(x, y, size*.85, mid*(math.Pi/180))
		deck.TextMid(tx, ty, fmt.Sprintf("%s "+datafmt+"%%", data[i].label, p), "sans", ts, "black")
		a1 = a2
	}
}

// chart makes either a pmap or donut
func chart(deck *generate.Deck, r io.ReadCloser) {
	data, _, _, title := getdata(r)
	if fullmarkup {
		deck.StartSlide()
	}
	if showpmap {
		pmap(deck, data, left, right, top, datacolor, title)
	} else {
		donut(deck, data, dx, dy, dsize, dwidth, datacolor, title)
	}
	if fullmarkup {
		deck.EndSlide()
	}
}

func main() {
	cmdflags()
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
