package main

import (
	"fmt"
	"github.com/ajstarks/deck/generate"
	"math"
	"math/rand"
	"os"
	"time"
)

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
	p.List(x, 60, size/2, list, ltype, "sans", "")
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
func boxtext(p *generate.Deck, x, y, w, h float64, s, font string, fontsize float64, bg, fg string) {
	p.Rect(x, y, w, h, bg)
	p.TextMid(x, y-(fontsize/2), s, font, fontsize, fg)
}

// rarrow draws an arrow pointing to the right
func rarrow(p *generate.Deck, x, y, w, h, aw, ah float64, color string, opacity float64) {
	xw := x - w
	xa := x - aw
	h2 := h / 2
	ah2 := ah / 2
	px := []float64{xw, xa, xa, x, xa, xa, xw, xw}
	py := []float64{y + h2, y + h2, y + ah2, y, y - ah2, y - h2, y - h2, y + h2}
	p.Polygon(px, py, color, opacity)
}

// larrow draws an arrow pointing to the left
func larrow(p *generate.Deck, x, y, w, h, aw, ah float64, color string, opacity float64) {
	xw := x + w
	xa := x + aw
	h2 := h / 2
	ah2 := ah / 2
	px := []float64{xw, xa, xa, x, xa, xa, xw, xw}
	py := []float64{y + h2, y + h2, y + ah2, y, y - ah2, y - h2, y - h2, y + h2}
	p.Polygon(px, py, color, opacity)
}

// darrow draws an arrow pointing down 
func darrow(p *generate.Deck, x, y, w, h, aw, ah float64, color string, opacity float64) {
	w2 := w / 2
	aw2 := aw / 2
	px := []float64{x, x + aw2, x + w2, x + w2, x - w2, x - w2, x - aw2, x}
	py := []float64{y, y + ah, y + ah, y + h, y + h, y + ah, y + ah, y}
	p.Polygon(px, py, color, opacity)

}

// uarrow draws an arrow pointing upwards
func uarrow(p *generate.Deck, x, y, w, h, aw, ah float64, color string, opacity float64) {
	w2 := w / 2
	aw2 := aw / 2
	px := []float64{x, x - aw2, x - w2, x - w2, x + w2, x + w2, x + aw2, x}
	py := []float64{y, y - ah, y - ah, y - h, y - h, y - ah, y - ah, y}
	p.Polygon(px, py, color, opacity)
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
		deck.TextMid(randp(100), randp(100), "hello", fontnames[rand.Intn(3)] , randp(10), randcolor(), randp(100))
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
	deck.Code(35, 35, hellorun, 2, 40, "rgb(127,0,0)")
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
		deck.Image(x, y, 100, 100, "sm.png")
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
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		w := randp(50)
		h := 3.0
		aw := h * 2
		ah := aw * 0.9
		larrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Right arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		w := randp(50)
		h := 3.0
		aw := h * 2
		ah := aw * 0.9
		rarrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Up arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		w := 3.0
		h := randp(50)
		aw := w * 2
		ah := aw * 0.9
		uarrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
	}
	deck.EndSlide()

	// Down arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		w := 3.0
		h := randp(50)
		aw := w * 2
		ah := aw * 0.9
		darrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))

	}
	deck.EndSlide()

	// Left and right arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		w := randp(50)
		h := 3.0
		aw := h * 2
		ah := aw * 0.9
		if i%2 == 0 {
			larrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
		} else {
			rarrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
		}
	}
	deck.EndSlide()

	// Up and down arrow
	deck.StartSlide()
	for i := 0; i < n; i++ {
		x := randp(100)
		y := randp(100)
		w := 3.0
		h := randp(50)
		aw := w * 2
		ah := aw * 0.9
		if i%2 == 0 {
			uarrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
		} else {
			darrow(deck, x, y, w, h, aw, ah, randcolor(), randp(100))
		}
	}
	deck.EndSlide()

	// Colored text box
	btext := map[string]string{"eat":"green", "sleep":"gray", "pray":"blue", "love":"red"}
	boxcount := 0
	bx := 50.0
	by := 80.0
	bw := 30.0
	bh := 15.0
	fontsize := bw/10.0
	var font string
	deck.StartSlide()
	for  s,color := range btext {
		boxcount++
		if boxcount%2 == 0 {
			font = "sans"
		} else {
			font = "serif"
		}
		boxtext(deck, bx, by, bw, bh, s, font, fontsize, color, "white")
		by -= bh  * 1.2
		if by < bh/2 {
			bx += bw * 1.5
			y = 80.0
		}
	}
	deck.EndSlide()

	// Diagram
	deck.StartSlide()
	boxtext(deck, 25, 50, 20, 15, "Urgent", "sans", 3, "red", "white")
	boxtext(deck, 75, 50, 20, 15, "Important", "sans", 3, "green", "white")
	rarrow(deck, 50, 50, 20, 2, 2, 5, "red", 40)
	larrow(deck, 50, 50, 20, 2, 2, 5, "green", 40)
	deck.EndSlide()

	deck.EndDeck()
}
