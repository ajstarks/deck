// pdfdeck: make PDF slide decks
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"code.google.com/p/gofpdf"
	"github.com/ajstarks/deck"
)

// fontmap maps generic font names to specific implementation names
var fontmap = map[string]string{}

const USLetterHeight = 612
const USLetterWidth = 792

// grid makes a grid using a percentage scale
func grid(doc *gofpdf.Fpdf, w, h float64, percent float64) {
	pw := w * (percent / 100)
	ph := h * (percent / 100)
	doc.SetLineWidth(1)
	for x := 0.0; x <= w; x += pw {
		doc.Line(x, 0, x, h)
	}
	for y := 0.0; y <= h; y += ph {
		doc.Line(0, y, w, y)
	}
}

// doline draws a line
func doline(doc *gofpdf.Fpdf, xp1, yp1, xp2, yp2, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Line(xp1, yp1, xp2, yp2)
}

// dorect draws a rectangle
func dorect(doc *gofpdf.Fpdf, x, y, w, h float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Rect(x, y, w, h, "F")
}

// bullet draws a rectangular bullet
func bullet(doc *gofpdf.Fpdf, x, y, size float64, color string) {
	rs := size / 2
	dorect(doc, x-size, y-rs, rs, rs, color)
}

// background places a colored rectangle
func background(doc *gofpdf.Fpdf, w, h float64, color string) {
	dorect(doc, 0, 0, w, h, color)
}

// dotext places text elements on the canvas according to type
func dotext(doc *gofpdf.Fpdf, cw, x, y, fs float64, wp float64, tdata, font, color, align, ttype string) {
	var tw float64

	td := strings.Split(tdata, "\n")
	red, green, blue := colorlookup(color)
	doc.SetTextColor(red, green, blue)
	if ttype == "code" {
		font = "mono"
		ch := float64(len(td)) * 1.8 * fs
		tw = deck.Pwidth(wp, cw, cw-x-20)
		dorect(doc, x-fs, y-fs, tw, ch, "rgb(240,240,240)")
	}
	if ttype == "block" {
		tw = deck.Pwidth(wp, cw, cw/2)
		textwrap(doc, x, y, tw, fs, fs*1.4, tdata, font)
	} else {
		ls := 1.8 * fs
		for _, t := range td {
			showtext(doc, x, y, t, fs, font, align)
			y += ls
		}
	}
}

// showtext places fully attributed text at the specified location
func showtext(doc *gofpdf.Fpdf, x, y float64, s string, fs float64, font, align string) {
	offset := 0.0
	doc.SetFont(fontlookup(font), "", fs)
	tw := doc.GetStringWidth(s)
	switch align {
	case "center":
		offset = (tw / 2)
	case "right":
		offset = tw
	}
	doc.Text(x-offset, y, s)
}

// dolists places lists on the canvas
func dolist(doc *gofpdf.Fpdf, x, y, fs float64, tdata []string, font, color, ltype string) {
	if font == "" {
		font = "sans"
	}
	doc.SetFont(fontlookup(font), "", fs)
	red, green, blue := colorlookup(color)
	doc.SetTextColor(red, green, blue)
	if ltype == "bullet" {
		x += fs
	}
	ls := 2.0 * fs
	for i, t := range tdata {
		if ltype == "number" {
			t = fmt.Sprintf("%d. ", i+1) + t
		}
		if ltype == "bullet" {
			bullet(doc, x, y, fs, color)
		}
		doc.Text(x, y, t)
		y += ls
	}
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// fontlookup maps font aliases to implementation font names
func fontlookup(s string) string {
	font, ok := fontmap[s]
	if ok {
		return font
	}
	return "sans"
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *gofpdf.Fpdf, x, y, w, fs, leading float64, s, font string) {
	const factor = 0.3
	doc.SetFont(fontlookup(font), "", fs)
	wordspacing := doc.GetStringWidth("m")
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w
	for _, s := range words {
		tw := doc.GetStringWidth(s)
		doc.Text(xp, yp, s)
		xp += tw + (wordspacing * factor)
		if xp > edge {
			xp = x
			yp += leading
		}
	}
}

// pct converts percentages to canvas measures
func pct(p, m float64) float64 {
	return (p / 100.0) * m
}

// doslides reads the deck file, making the PDF version
func doslides(doc *gofpdf.Fpdf, filename string, w, h int, gp float64) {
	var d deck.Deck
	var err error

	for _, v := range fontmap {
		doc.AddFont(v, "", v+".json")
	}
	d, err = deck.Read(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	d.Canvas.Width = USLetterWidth
	d.Canvas.Height = USLetterHeight

	for i := 0; i < len(d.Slide); i++ {
		pdfslide(doc, d, i, gp)
	}

	err = doc.Output(os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

// dimen returns canvas dimensions from percentages
func dimen(w, h, xp, yp, sp float64) (float64, float64, float64) {
	return pct(xp, w), pct(100-yp, h), pct(sp, w) * 1.2
}

// pdfslide makes a slide, one slide per PDF page
func pdfslide(doc *gofpdf.Fpdf, d deck.Deck, n int, gp float64) {
	if n < 0 || n > len(d.Slide)-1 {
		return
	}
	var x, y, fs float64

	doc.AddPage()
	cw := float64(d.Canvas.Width)
	ch := float64(d.Canvas.Height)

	slide := d.Slide[n]

	// set background, if specified
	if len(slide.Bg) > 0 {
		background(doc, cw, ch, slide.Bg)
	}
	// set the default foreground
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	if gp > 0 {
		grid(doc, cw, ch, gp)
	}
	// for every image on the slide...
	for _, im := range slide.Image {
		x, y, _ = dimen(cw, ch, im.Xp, im.Yp, 0)
		fw, fh := float64(im.Width), float64(im.Height)
		doc.Image(im.Name, x-(fw/2), y-(fh/2), fw, fh, false, "", 0, "")
		if len(im.Caption) > 0 {
			capsize := pct(2, cw)
			showtext(doc, x, y+(fh/2)+(capsize*2), im.Caption, capsize, "sans", "center")
		}
	}

	// every graphic on the slide
	const defaultColor = "rgb(127,127,127)"
	// line
	for _, line := range slide.Line {
		if line.Color == "" {
			line.Color = defaultColor
		}
		x1, y1, sw := dimen(cw, ch, line.Xp1, line.Yp1, line.Sp)
		x2, y2, _ := dimen(cw, ch, line.Xp2, line.Yp2, 0)
		if sw == 0 {
			sw = 2.0
		}
		doline(doc, x1, y1, x2, y2, sw, line.Color)
	}
	// rect
	for _, rect := range slide.Rect {
		x, y, _ := dimen(cw, ch, rect.Xp, rect.Yp, 0)
		w := pct(rect.Wp, cw)
		h := pct(rect.Hp, cw)
		if rect.Color == "" {
			rect.Color = defaultColor
		}
		dorect(doc, x-(w/2), y-(h/2), w, h, rect.Color)
	}
	// for every text element...
	for _, t := range slide.Text {
		if t.Color == "" {
			t.Color = slide.Fg
		}
		if t.Font == "" {
			t.Font = "sans"
		}
		x, y, fs = dimen(cw, ch, t.Xp, t.Yp, t.Sp)
		dotext(doc, cw, x, y, fs, t.Wp, t.Tdata, t.Font, t.Color, t.Align, t.Type)
	}
	// for every list element...
	for _, l := range slide.List {
		if l.Color == "" {
			l.Color = slide.Fg
		}
		x, y, fs = dimen(cw, ch, l.Xp, l.Yp, l.Sp)
		dolist(doc, x, y, fs, l.Li, l.Font, l.Color, l.Type)
	}
}

// dodeck kicks things off
func dodeck(filename, fontdir string, gp float64) {
	doc := gofpdf.New("L", "pt", "Letter", fontdir)
	doslides(doc, filename, USLetterWidth, USLetterHeight, gp)
}

// for every file, make a deck
func main() {
	var gridpct = flag.Float64("g", 0, "place percentage grid on each slide")
	var fontdir = flag.String("f", path.Join(os.Getenv("GOPATH"), "src/code.google.com/p/gofpdf/font"), "font directory")
	var sansfont = flag.String("sans", "helvetica", "sans font")
	var serifont = flag.String("serif", "times", "serif font")
	var monofont = flag.String("mono", "courier", "mono font")
	flag.Parse()
	fontmap["sans"] = *sansfont
	fontmap["serif"] = *serifont
	fontmap["mono"] = *monofont

	for _, f := range flag.Args() {
		dodeck(f, *fontdir, *gridpct)
	}
}
