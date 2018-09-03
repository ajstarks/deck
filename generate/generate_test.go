package generate

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
)

var canvas *Deck

const (
	lorem    = `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
	hellogo  = "package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"hello, world\")\n}"
	hellorun = "$ go run hello.go\nhello, world"
	rgbfmt   = `rgb(%d,%d,%d)`
)

// randp returns a random float
func randp(n float64) float64 {
	x := math.Ceil(rand.Float64() * n)
	if x == 0 {
		return 1
	}
	return x
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

func TestMain(m *testing.M) {
	f, err := os.Create("decktest.xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	canvas = NewSlides(f, 0, 0)
	m.Run()
}

func BenchmarkEmptyDoc(b *testing.B) {
	canvas.StartDeck()
	canvas.EndDeck()
}

func BenchmarkEmptySlide(b *testing.B) {
	canvas.StartDeck()
	canvas.StartSlide()
	canvas.EndSlide()
	canvas.EndDeck()
}

func BenchmarkHello(b *testing.B) {
	canvas.StartDeck()
	canvas.StartSlide("black", "white")
	canvas.Circle(50, 0, 100, "blue")
	canvas.Text(50, 25, "hello, world", "sans", 10, "")
	canvas.EndSlide()
	canvas.EndDeck()
}

func BenchmarkList(b *testing.B) {
	items := []string{"One", "Two", "Three", "Four", "Five", "Six", "Seven"}
	canvas.StartDeck()
	canvas.StartSlide()
	canvas.Text(10, 90, "Important Items", "sans", 5, "")
	canvas.List(10, 70, 4, 0, 0, items, "bullet", "sans", "red")
	canvas.EndSlide()
	canvas.EndDeck()

}

func BenchmarkCircle(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		canvas.Circle(100, 100, 10, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkEllipse(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		canvas.Ellipse(100, 100, 10, 10, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkSquare(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		canvas.Square(100, 100, 10, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkRect(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		canvas.Rect(100, 100, 10, 10, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkArc(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		canvas.Arc(100, 100, 10, 10, 0.5, 90, 300, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkCurve(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		canvas.Curve(100, 100, 100, 100, 100, 100, 0.5, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkLine(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		x := 100.0
		y := 100.0
		if i%2 == 0 {
			canvas.Line(x, y, x+10, y, 0.5, "red", 100)
		} else {
			canvas.Line(x, y, x, y+10, 0.5, "red", 100)
		}
	}
	canvas.EndSlide()
}

func BenchmarkPolygon(b *testing.B) {
	canvas.StartSlide()
	for i := 0; i < b.N; i++ {
		px, py := rp(100, 100, 10, 5)
		canvas.Polygon(px, py, "red", 100)
	}
	canvas.EndSlide()
}

func BenchmarkImage(b *testing.B) {
	canvas.StartSlide("gray")
	y := 50.0
	for x := 20.0; x < 90.0; x += 20.0 {
		canvas.Image(x, y, 100, 100, "sm.png", "")
	}
	canvas.EndSlide()
}

func BenchmarkText(b *testing.B) {
	// Text Block
	canvas.StartSlide("rgb(180,180,180)")
	canvas.TextBlock(10, 90, lorem, "sans", 2.5, 30, "black")
	canvas.TextBlock(10, 50, lorem, "sans", 2.5, 50, "gray")
	canvas.TextBlock(10, 20, lorem, "sans", 2.5, 75, "white")
	canvas.EndSlide()

	// Text alignment
	canvas.StartSlide("rgb(180,180,180)")
	canvas.Text(50, 80, "left", "sans", 10, "black")
	canvas.TextMid(50, 50, "center", "serif", 10, "gray")
	canvas.TextEnd(50, 20, "right", "mono", 10, "white")
	canvas.Line(50, 100, 50, 0, 0.2, "black", 20)
	canvas.EndSlide()

	// Code
	canvas.StartSlide("rgb(180,180,180)")
	canvas.Code(35, 80, "$ edit hello.go", 2, 40, "rgb(0,0,127)")
	canvas.Code(35, 70, hellogo, 2, 40, "black")
	canvas.Code(35, 39, hellorun, 2, 40, "rgb(127,0,0)")
	canvas.EndSlide()
}
