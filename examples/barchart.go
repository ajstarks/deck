package main

import (
   "fmt"
   "github.com/ajstarks/deck/generate"
   "os"
)
type Bardata struct {
   label string
   value float64
}
func vmap(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
   return low2 + (high2-low2)*(value-low1)/(high1-low1)
}
func main() {
   benchmarks := []Bardata{
      {"Macbook Air", 154.701}, {"MacBook Pro (2008)", 289.603}, {"BeagleBone Black", 2896.037}, {"Raspberry Pi", 5765.568},
   }
   maxdata := 5800.0
   ts := 2.5
   hts := ts / 2
   x, y := 10.0, 60.0
   bx1 := x + (ts * 12)
   bx2 := bx1 + 50.0
   linespacing := ts * 2.0
   deck := generate.NewSlides(os.Stdout, 0, 0)
   deck.StartDeck()
   deck.StartSlide("rgb(255,255,255)")
   deck.Text(x, y+20, "Go 1.1.2 Build and Test Times", "sans", ts*2, "black")
   for _, data := range benchmarks {
      deck.Text(x, y, data.label, "sans", ts, "rgb(100,100,100)")
      bv := vmap(data.value, 0, maxdata, bx1, bx2)
      deck.Line(bx1, y+hts, bv, y+hts, ts, "lightgray")
      deck.Text(bv+0.5, y+(hts/2), fmt.Sprintf("%.1f", data.value), "sans", hts, "rgb(127,0,0)")
      y -= linespacing
   }
   deck.EndSlide()
   deck.EndDeck()
}
