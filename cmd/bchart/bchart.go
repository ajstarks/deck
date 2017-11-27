//
// bchart - make barcharts in the deck format
//
// bchart reads data from the standard input file, expecting a tab-separated list of text,data pairs
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
//	  -color barcolor (default "rgb(175,175,175)")
//	  -datafmt data format (default "%.1f")
//	  -dmin zero minimum
//	  -dot draw a line and dot instead of a solid bar
//	  -left left margin (default 20)
//	  -textsize text size (default 1.2)
//	  -top top of the chart (default 90)
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/deck/generate"
)

type Bardata struct {
	label string
	value float64
}

func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func getdata(r io.Reader) ([]Bardata, float64, float64, string) {
	var (
		data []Bardata
		d    Bardata
		err  error
	)

	maxval := -1.0
	minval := 1e50
	title := ""
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if t[0] == '#' && len(t) > 2 {
			title = t[1:]
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
	return data, minval, maxval, title
}
func main() {
	var (
		ts, left, top      float64
		dot, datamin       bool
		datacolor, datafmt string
	)
	const (
		bgcolor      = "white"
		titlecolor   = "black"
		labelcolor   = "rgb(75,75,75)"
		dotlinecolor = "lightgray"
		valuecolor   = "rgb(127,0,0)"
	)
	flag.Float64Var(&ts, "textsize", 1.2, "text size")
	flag.Float64Var(&left, "left", 20.0, "left margin")
	flag.Float64Var(&top, "top", 90.0, "top")
	flag.BoolVar(&dot, "dot", false, "dot and line")
	flag.BoolVar(&datamin, "dmin", false, "zero minimum")
	flag.StringVar(&datacolor, "color", "rgb(175,175,175)", "bar color")
	flag.StringVar(&datafmt, "datafmt", "%.1f", "data format")
	flag.Parse()

	hts := ts / 2
	right := 100 - left
	linespacing := ts * 2.4

	bardata, mindata, maxdata, title := getdata(os.Stdin)
	if !datamin {
		mindata = 0
	}
	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	deck.StartSlide(bgcolor)
	if title != "" {
		deck.TextMid(50, top+(linespacing*1.5), title, "serif", ts*1.5, titlecolor)
	}
	y := top
	for _, data := range bardata {
		deck.TextEnd(left-hts, y, data.label, "sans", ts, labelcolor)
		bv := vmap(data.value, mindata, maxdata, left, right)
		if dot {
			deck.Line(left, y+hts, bv, y+hts, ts/4, dotlinecolor)
			deck.Circle(bv, y+hts, hts, datacolor)
		} else {
			deck.Line(left, y+hts, bv, y+hts, ts, datacolor)
		}
		deck.Text(bv+hts, y+(hts/2), fmt.Sprintf(datafmt, data.value), "mono", hts, valuecolor)
		y -= linespacing
	}
	deck.EndSlide()
	deck.EndDeck()
}
