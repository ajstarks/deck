package main

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"math"
	"os"
	"time"
)

var (
	canvas = svg.New(os.Stdout)
	width  = 500
	height = 500
)

func vmap(value float64, l1 float64, h1 float64,
	l2 float64, h2 float64) float64 {
	return l2 + (h2-l2)*(value-l1)/(h1-l1)
}

// See: Processing (Reas and Fry), pg. 247
func main() {
	w2, h2 := width/2, height/2
	h, m, s := time.Now().Clock()
	sec := vmap(float64(s), 0, 60, 0, math.Pi*2) - math.Pi/2
	min := vmap(float64(m), 0, 60, 0, math.Pi*2) - math.Pi/2
	hour := vmap(float64(h%12), 0, 12, 0, math.Pi*2) - math.Pi/2
	secpct := float64(width) * 0.38
	minpct := float64(width) * 0.30
	hourpct := float64(width) * 0.25
	facepct := (width * 40) / 100
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	canvas.Text(w2, 25, fmt.Sprintf("%02d:%02d:%02d", h, m, s), "text-anchor:middle;font-size:12pt;fill:white")
	canvas.Circle(w2, h2, facepct, canvas.RGB(100, 100, 100))
	canvas.Gstyle("stroke:white;stroke-width:20;stroke-opacity:0.6;stroke-linecap:round")
	canvas.Line(w2, h2, int(math.Cos(sec)*secpct)+w2, int(math.Sin(sec)*secpct)+h2, "stroke:red;stroke-width:5")
	canvas.Line(w2, h2, int(math.Cos(min)*minpct)+w2, int(math.Sin(min)*minpct)+h2)
	canvas.Line(w2, h2, int(math.Cos(hour)*hourpct)+w2, int(math.Sin(hour)*hourpct)+h2)
	canvas.Gend()
	canvas.End()
}
