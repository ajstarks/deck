package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/deck/generate"
)

type dimension struct {
	x, y, w, h float64
}

const (
	lorem    = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
	hellogo  = "package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"hello, world\")\n}"
	hellorun = "$ go run hello.go\nhello, world"
	rgbfmt   = `rgb(%d,%d,%d)`
)

// randcolor returns a random color in RGB format
func randcolor() string {
	return fmt.Sprintf(rgbfmt, rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

// randp returns a random float
func randp(n float64) float64 {
	x := math.Ceil(rand.Float64() * n)
	if x == 0 {
		return 1
	}
	return x
}

// randpoly returns the polygon coordinates centered from a point
func randpoly(cx, cy, size float64, np int) (string, string) {
	rx, ry := rp(cx, cy, size, np)
	return generate.Polycoord(rx, ry)
}

// rp returns n random coordinates radiating from (cx, cy)
func rp(cx, cy, size float64, np int) ([]float64, []float64) {
	adiv := 360.0 / float64(np)
	rx := make([]float64, np)
	ry := make([]float64, np)
	a := 0.0
	for i := 0; i < np; i++ {
		t := a * (math.Pi / 180)
		r := randp(size)
		rx[i] = (r * math.Cos(t)) + cx
		ry[i] = (r * math.Sin(t)) + cy
		a += adiv
	}
	return rx, ry
}

// section makes a slide "section"
func section(p *generate.Deck, name string, n int) {
	x := 10.0
	y := 70.0
	size := 5.0

	ry := y - (size / 2)
	p.Text(x, y, name, "sans", size, "")
	p.Line(x, ry, 100-x, ry, 0.1, "black")
	p.Text(80, 10, fmt.Sprintf("%d", n), "sans", size/3, "gray")
	list := []string{}
	for i := 0; i < 5; i++ {
		list = append(list, fmt.Sprintf("item number %d", (i+1)*(n*10)))
	}
	var ltype string
	switch {
	case n%2 == 0:
		ltype = "bullet"
	case n%3 == 0:
		ltype = "number"
	default:
		ltype = "plain"
	}
	p.List(x, 60, size/2, 0, 0, list, ltype, "sans", "")
}

// fitlist fits a list in a vertical range
func fitlist(p *generate.Deck, x, y, size float64, list []string) {
	yp := y
	colsize := 35.0
	interval := y/float64(len(list)) + size
	for _, s := range list {
		p.Text(x, yp, s, "sans", size, "black")
		yp -= interval
		if yp < size {
			x += colsize
			yp = y
		}
	}
}

// boxtext makes a colored box with centered text
func (d *dimension) boxtext(p *generate.Deck, s, font string, fontsize float64, bg, fg string) {
	p.Rect(d.x, d.y, d.w, d.h, bg)
	p.TextMid(d.x, d.y-(fontsize/2), s, font, fontsize, fg)
}

// rarrow draws an arrow pointing to the right
func (d *dimension) rarrow(p *generate.Deck, aw, ah float64, color string, opacity float64) {
	xw := d.x - d.w
	xa := d.x - aw
	h2 := d.h / 2
	ah2 := ah / 2
	px := []float64{xw, xa, xa, d.x, xa, xa, xw, xw}
	py := []float64{d.y + h2, d.y + h2, d.y + ah2, d.y, d.y - ah2, d.y - h2, d.y - h2, d.y + h2}
	p.Polygon(px, py, color, opacity)
}

// larrow draws an arrow pointing to the left
func (d *dimension) larrow(p *generate.Deck, aw, ah float64, color string, opacity float64) {
	xw := d.x + d.w
	xa := d.x + aw
	h2 := d.h / 2
	ah2 := ah / 2
	px := []float64{xw, xa, xa, d.x, xa, xa, xw, xw}
	py := []float64{d.y + h2, d.y + h2, d.y + ah2, d.y, d.y - ah2, d.y - h2, d.y - h2, d.y + h2}
	p.Polygon(px, py, color, opacity)
}

// darrow draws an arrow pointing down
func (d *dimension) darrow(p *generate.Deck, aw, ah float64, color string, opacity float64) {
	w2 := d.w / 2
	aw2 := aw / 2
	px := []float64{d.x, d.x + aw2, d.x + w2, d.x + w2, d.x - w2, d.x - w2, d.x - aw2, d.x}
	py := []float64{d.y, d.y + ah, d.y + ah, d.y + d.h, d.y + d.h, d.y + ah, d.y + ah, d.y}
	p.Polygon(px, py, color, opacity)

}

// uarrow draws an arrow pointing upwards
func (d *dimension) uarrow(p *generate.Deck, aw, ah float64, color string, opacity float64) {
	w2 := d.w / 2
	aw2 := aw / 2
	px := []float64{d.x, d.x - aw2, d.x - w2, d.x - w2, d.x + w2, d.x + w2, d.x + aw2, d.x}
	py := []float64{d.y, d.y - ah, d.y - ah, d.y - d.h, d.y - d.h, d.y - ah, d.y - ah, d.y}
	p.Polygon(px, py, color, opacity)
}

// cube makes a cube with its front-facing bottom edge at x,y, lit from the top
func cube(d *generate.Deck, x, y, w, h float64, color string) {
	xc := make([]float64, 4)
	yc := make([]float64, 4)

	w50 := w / 2
	h20 := h * .2
	h40 := h20 * 2
	yh := y + h

	xc[0], xc[1], xc[2], xc[3] = x-w50, x, x+w50, x
	yc[0], yc[1], yc[2], yc[3] = yh-h20, yh, yh-h20, yh-(h40)
	d.Polygon(xc, yc, color, 20) // left face

	xc[0], xc[1], xc[2], xc[3] = x-w50, x-w50, x, x
	yc[0], yc[1], yc[2], yc[3] = y+h20, yh-h20, yh-h40, y
	d.Polygon(xc, yc, color, 50) // right face

	xc[0], xc[1], xc[2], xc[3] = x, x, x+w50, x+w50
	yc[0], yc[1], yc[2], yc[3] = y, yh-h40, yh-h20, y+h20
	d.Polygon(xc, yc, color, 70) // top face
}

// Test slide generation
func main() {
	rand.Seed(time.Now().Unix() % 1e9)

	n := 200
	sections := []string{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten"}

	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()

	// Text
	fontnames := []string{"sans", "serif", "mono"}
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.TextMid(randp(100), randp(100), "hello", fontnames[rand.Intn(3)], randp(10), randcolor(), randp(100))
	}
	deck.EndSlide()

	// Text Block
	deck.StartSlide("rgb(180,180,180)")
	deck.TextBlock(10, 90, lorem, "sans", 2.5, 30, "black")
	deck.TextBlock(10, 50, lorem, "sans", 2.5, 50, "gray")
	deck.TextBlock(10, 20, lorem, "sans", 2.5, 75, "white")
	deck.EndSlide()

	// Text alignment
	deck.StartSlide("rgb(180,180,180)")
	deck.Text(50, 80, "left", "sans", 10, "black")
	deck.TextMid(50, 50, "center", "serif", 10, "gray")
	deck.TextEnd(50, 20, "right", "mono", 10, "white")
	deck.Line(50, 100, 50, 0, 0.2, "black", 20)
	deck.EndSlide()

	// Code
	deck.StartSlide("rgb(180,180,180)")
	deck.Code(35, 80, "$ edit hello.go", 2, 40, "rgb(0,0,127)")
	deck.Code(35, 70, hellogo, 2, 40, "black")
	deck.Code(35, 39, hellorun, 2, 40, "rgb(127,0,0)")
	deck.EndSlide()

	// Fit text
	deck.StartSlide()
	fitlist(deck, 10, 70, 5, sections)
	deck.EndSlide()

	// Lists
	for i := 0; i < 3; i++ {
		deck.StartSlide()
		section(deck, sections[i], i+1)
		deck.EndSlide()
	}

	// Image
	deck.StartSlide("gray")
	y := 50.0
	for x := 20.0; x < 90.0; x += 20.0 {
		deck.Image(x, y, 100, 100, "sm.png", "")
	}
	deck.EndSlide()

	// Circle
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.Circle(randp(100), randp(100), randp(10), randcolor(), randp(100))
	}
	deck.EndSlide()

	// Ellipse
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.Ellipse(randp(100), randp(100), randp(10), randp(10), randcolor(), randp(100))
	}
	deck.EndSlide()

	// Square
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.Square(randp(100), randp(100), randp(10), randcolor(), randp(100))
	}
	deck.EndSlide()

	// Rect
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.Rect(randp(100), randp(100), randp(10), randp(10), randcolor(), randp(100))
	}
	deck.EndSlide()

	// Arc
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.Arc(randp(100), randp(100), randp(10), randp(10), 0.5, randp(90), randp(300), randcolor(), randp(100))
	}
	deck.EndSlide()

	// Curve
	deck.StartSlide()
	for i := 0; i < n; i++ {
		deck.Curve(randp(100), randp(100), randp(100), randp(100), randp(100), randp(100), 0.5, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Line
	deck.StartSlide()
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		if i%2 == 0 {
			deck.Line(x, y, x+randp(10), y, 0.5, randcolor(), randp(100))
		} else {
			deck.Line(x, y, x, y+randp(10), 0.5, randcolor(), randp(100))
		}
	}
	deck.EndSlide()

	// Polygon
	deck.StartSlide()
	for i := 0; i < n; i++ {
		px, py := rp(randp(100), randp(100), randp(10), 5)
		deck.Polygon(px, py, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Left arrow
	deck.StartSlide()
	dim := &dimension{h: 3.0}
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.w = randp(50)
		aw := dim.h * 2
		ah := aw * 0.9
		dim.larrow(deck, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Right arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.w = randp(50)
		aw := dim.h * 2
		ah := aw * 0.9
		dim.rarrow(deck, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Up arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.w = 3.0
		dim.h = randp(50)
		aw := dim.w * 2
		ah := aw * 0.9
		dim.uarrow(deck, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Down arrow
	dim.w = 3.0
	deck.StartSlide()
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.h = randp(50)
		aw := dim.w * 2
		ah := aw * 0.9
		dim.darrow(deck, aw, ah, randcolor(), randp(100))

	}
	deck.EndSlide()

	// Left and right arrow
	dim.h = 3.0
	deck.StartSlide()
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.w = randp(50)
		aw := dim.h * 2
		ah := aw * 0.9
		if i%2 == 0 {
			dim.larrow(deck, aw, ah, randcolor(), randp(100))
		} else {
			dim.rarrow(deck, aw, ah, randcolor(), randp(100))
		}
	}
	deck.EndSlide()

	// Up and down arrow
	dim.w = 3.0
	deck.StartSlide()
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.h = randp(50)
		aw := dim.w * 2
		ah := aw * 0.9
		if i%2 == 0 {
			dim.uarrow(deck, aw, ah, randcolor(), randp(100))
		} else {
			dim.darrow(deck, aw, ah, randcolor(), randp(100))
		}
	}
	deck.EndSlide()

	// Left fat arrows
	dim.h = 3.0
	deck.StartSlide()
	for i := 0; i < n; i++ {
		dim.x = randp(100)
		dim.y = randp(100)
		dim.w = randp(10) // randp(50)
		aw := dim.w / 10  // dim.h * 2
		ah := dim.h / 3   // aw * 0.9
		dim.larrow(deck, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Colored text box
	btext := map[string]string{"eat": "green", "sleep": "gray", "pray": "blue", "love": "red"}
	boxcount := 0
	dim.x = 50
	dim.y = 80
	dim.w = 30
	dim.h = 15
	fontsize := dim.w / 10.0
	var font string
	deck.StartSlide()
	for s, color := range btext {
		boxcount++
		if boxcount%2 == 0 {
			font = "sans"
		} else {
			font = "serif"
		}
		dim.boxtext(deck, s, font, fontsize, color, "white")
		dim.y -= dim.h * 1.2
		if dim.y < dim.h/2 {
			dim.x += dim.w * 1.5
			dim.y = 80.0
		}
	}
	deck.EndSlide()

	// Diagram
	deck.StartSlide()
	dim.x = 25
	dim.y = 50
	dim.w = 20
	dim.h = 15
	dim.boxtext(deck, "Urgent", "sans", 3, "red", "white")
	dim.x = 75
	dim.boxtext(deck, "Important", "sans", 3, "green", "white")
	dim.x = 50
	dim.y = 50
	dim.w = 20
	dim.h = 15
	dim.rarrow(deck, 2, 5, "red", 40)
	dim.larrow(deck, 2, 5, "green", 40)
	deck.EndSlide()

	deck.StartSlide()
	cube(deck, 20, 50, 20, 25, "black")
	cube(deck, 50, 20, 20, 50, "gray")
	cube(deck, 80, 50, 20, 25, "rgb(127,0,0)")
	deck.EndSlide()

	deck.EndDeck()
}
