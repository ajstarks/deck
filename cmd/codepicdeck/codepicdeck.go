// codepicdeck: make code+pic slide decks
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ajstarks/deck/generate"
)

// expand tabs to spaces, and escape XML
var codemap = strings.NewReplacer(
	"\t", "    ",
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;")

// includefile returns the content of a file as a tab-expanded, XML-escaped string
func includefile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}

// imagesize returns the dimensions (w,h) of an image file
func imagesize(filename string) (int, int) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 0, 0
	}
	defer f.Close()
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 0, 0
	}
	return img.Width, img.Height
}

// index makes an image index slide
func index(deck *generate.Deck, filenames []string, title string) {
	deck.StartSlide()
	x, y := 5.0, 90.0
	for _, f := range filenames {
		imagefile := swapext(f, ".go", ".png")
		deck.Image(x, y, 72, 72, imagefile)
		deck.TextMid(x, y-7, f, "sans", 1, "black")
		x += 10.0
		if x > 95 {
			x = 5.0
			y -= 20.0
		}
	}
	deck.TextMid(50, 96, title, "sans", 2, "black")
	deck.EndSlide()
}

// swapext swaps the specified file extensions
func swapext(s, fromext, toext string) string {
	i := strings.LastIndex(s, fromext)
	if i < 0 {
		return ""
	}
	return s[:i] + toext
}

// codepic makes a code and picture slides
func codepic(deck *generate.Deck, filenames []string) {
	slide := 0
	for _, codefile := range filenames {
		imagefile := swapext(codefile, ".go", ".png")
		if len(imagefile) == 0 {
			fmt.Fprintf(os.Stderr, "cannot get the basename for %s\n", codefile)
			continue
		}
		imw, imh := imagesize(imagefile)
		slide++
		deck.StartSlide()
		deck.Image(75, 68, imw, imh, imagefile)
		deck.Text(2.5, 96, includefile(codefile), "mono", 1.2, "black")
		deck.TextEnd(90, 2.5, codefile, "sans", 2, "black")
		deck.TextEnd(95, 2.5, fmt.Sprintf("[%d]", slide), "sans", 2, "gray")
		deck.EndSlide()
	}
}

func main() {
	doindex := flag.Bool("index", true, "generates an index slide")
	title := flag.String("title", "", "index title")
	flag.Parse()
	files := flag.Args()
	deck := generate.NewSlides(os.Stdout, 0, 0)
	deck.StartDeck()
	if *doindex {
		index(deck, files, *title)
	}
	codepic(deck, files)
	deck.EndDeck()
}
