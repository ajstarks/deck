// pdfdeck: make PDF slide decks
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/gofpdf"
	"github.com/ajstarks/deck"
)

const (
	mm2pt       = 2.83464 // mm to pt conversion
	linespacing = 1.4
	listspacing = 1.8
	fontfactor  = 1.2
)

// PageDimen describes page dimensions
// the unit field is used to convert to pt.
type PageDimen struct {
	width, height, unit float64
}

// fontmap maps generic font names to specific implementation names
var fontmap = map[string]string{}

// pagemap defines page dimensions
var pagemap = map[string]PageDimen{
	"Letter": {612, 792, 1},
	"Legal":  {612, 1008, 1},
	"A3":     {297, 420, mm2pt},
	"A4":     {210, 297, mm2pt},
	"A5":     {148, 210, mm2pt},
}

// pct converts percentages to canvas measures
func pct(p, m float64) float64 {
	return (p / 100.0) * m
}

// dimen returns canvas dimensions from percentages
func dimen(w, h, xp, yp, sp float64) (float64, float64, float64) {
	return pct(xp, w), pct(100-yp, h), pct(sp, w) * fontfactor
}

// setopacity sets the alpha value:
// 0 == default value (opaque)
// -1 == fully transparent
// > 0 set opacity percent
func setopacity(doc *gofpdf.Fpdf, v float64) {
	switch {
	case v < 0:
		doc.SetAlpha(0, "Normal")
	case v > 0:
		doc.SetAlpha(v/100, "Normal")
	case v == 0:
		doc.SetAlpha(1, "Normal")
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

// grid makes a percentage scale
func grid(doc *gofpdf.Fpdf, w, h float64, color string, percent float64) {
	pw := w * (percent / 100)
	ph := h * (percent / 100)
	doc.SetLineWidth(0.5)
	r, g, b := colorlookup(color)
	doc.SetDrawColor(r, g, b)
	doc.SetTextColor(r, g, b)
	fs := pct(1, w)
	for x, pl := 0.0, 0.0; x <= w; x += pw {
		doc.Line(x, 0, x, h)
		if pl > 0 {
			showtext(doc, x, h-fs, fmt.Sprintf("%.0f", pl), fs, "sans", "center")
		}
		pl += percent
	}
	for y, pl := 0.0, 0.0; y <= h; y += ph {
		doc.Line(0, y, w, y)
		if pl < 100 {
			showtext(doc, fs, y+(fs/3), fmt.Sprintf("%.0f", 100-pl), fs, "sans", "center")
		}
		pl += percent
	}
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

// doline draws a line
func doline(doc *gofpdf.Fpdf, xp1, yp1, xp2, yp2, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Line(xp1, yp1, xp2, yp2)
}

// doarc draws a line
func doarc(doc *gofpdf.Fpdf, x, y, w, h, a1, a2, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Arc(x, y, w, h, 0, a1, a2, "D")
}

// docurve draws a bezier curve
func docurve(doc *gofpdf.Fpdf, xp1, yp1, xp2, yp2, xp3, yp3, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Curve(xp1, yp1, xp2, yp2, xp3, yp3, "D")
}

// dorect draws a rectangle
func dorect(doc *gofpdf.Fpdf, x, y, w, h float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Rect(x, y, w, h, "F")
}

// doellipse draws a rectangle
func doellipse(doc *gofpdf.Fpdf, x, y, w, h float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Ellipse(x, y, w, h, 0, "F")
}

// dotext places text elements on the canvas according to type
func dotext(doc *gofpdf.Fpdf, cw, x, y, fs float64, wp float64, tdata, font, color, align, ttype string) {
	var tw float64

	td := strings.Split(tdata, "\n")
	red, green, blue := colorlookup(color)
	doc.SetTextColor(red, green, blue)
	if ttype == "code" {
		font = "mono"
		ch := float64(len(td)) * listspacing * fs
		tw = deck.Pwidth(wp, cw, cw-x-20)
		dorect(doc, x-fs, y-fs, tw, ch, "rgb(240,240,240)")
	}
	if ttype == "block" {
		tw = deck.Pwidth(wp, cw, cw/2)
		textwrap(doc, x, y, tw, fs, fs*linespacing, tdata, font)
	} else {
		ls := listspacing * fs
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

// doslides reads the deck file, making the PDF version
func doslides(doc *gofpdf.Fpdf, filename, author, title string, w, h int, gp float64) {
	var d deck.Deck
	var err error

	for _, v := range fontmap {
		doc.AddFont(v, "", v+".json")
	}
	d, err = deck.Read(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pdfdeck: %v\n", err)
		return
	}

	// Lanscape mode switch w,h
	d.Canvas.Width = h
	d.Canvas.Height = w
	doc.SetCreator("PDFDeck", true)
	if len(title) > 0 {
		doc.SetTitle(title, true)
	}
	if len(author) > 0 {
		doc.SetAuthor(author, true)
	}
	for i := 0; i < len(d.Slide); i++ {
		pdfslide(doc, d, i, gp)
	}
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
		grid(doc, cw, ch, slide.Fg, gp)
	}
	// for every image on the slide...
	for _, im := range slide.Image {
		x, y, _ = dimen(cw, ch, im.Xp, im.Yp, 0)
		fw, fh := float64(im.Width), float64(im.Height)
		midx := fw / 2
		midy := fh / 2
		doc.Image(im.Name, x-midx, y-midy, fw, fh, false, "", 0, "")
		if len(im.Caption) > 0 {
			capsize := deck.Pwidth(im.Sp, cw, pct(2, cw))
			if im.Font == "" {
				im.Font = "sans"
			}
			if im.Color == "" {
				im.Color = slide.Fg
			}
			if im.Align == "" {
				im.Align = "center"
			}
			switch im.Align {
			case "left", "start":
				x -= midx
			case "right", "end":
				x += midx
			}
			capr, capg, capb := colorlookup(im.Color)
			doc.SetTextColor(capr, capg, capb)
			showtext(doc, x, y+(midy)+(capsize*2), im.Caption, capsize, im.Font, im.Align)
		}
	}
	// every graphic on the slide
	const defaultColor = "rgb(127,127,127)"
	// rect
	for _, rect := range slide.Rect {
		x, y, _ := dimen(cw, ch, rect.Xp, rect.Yp, 0)
		w := pct(rect.Wp, cw)
		h := pct(rect.Hp, cw)
		if rect.Color == "" {
			rect.Color = defaultColor
		}
		setopacity(doc, rect.Opacity)
		dorect(doc, x-(w/2), y-(h/2), w, h, rect.Color)
	}
	// ellipse
	for _, ellipse := range slide.Ellipse {
		x, y, _ := dimen(cw, ch, ellipse.Xp, ellipse.Yp, 0)
		w := pct(ellipse.Wp, cw)
		h := pct(ellipse.Hp, cw)
		if ellipse.Color == "" {
			ellipse.Color = defaultColor
		}
		setopacity(doc, ellipse.Opacity)
		doellipse(doc, x, y, w/2, h/2, ellipse.Color)
	}
	// curve
	for _, curve := range slide.Curve {
		if curve.Color == "" {
			curve.Color = defaultColor
		}
		setopacity(doc, curve.Opacity)
		x1, y1, sw := dimen(cw, ch, curve.Xp1, curve.Yp1, curve.Sp)
		x2, y2, _ := dimen(cw, ch, curve.Xp2, curve.Yp2, 0)
		x3, y3, _ := dimen(cw, ch, curve.Xp3, curve.Yp3, 0)
		if sw == 0 {
			sw = 2.0
		}
		docurve(doc, x1, y1, x2, y2, x3, y3, sw, curve.Color)
	}
	// arc
	for _, arc := range slide.Arc {
		if arc.Color == "" {
			arc.Color = defaultColor
		}
		setopacity(doc, arc.Opacity)
		x, y, sw := dimen(cw, ch, arc.Xp, arc.Yp, arc.Sp)
		w := pct(arc.Wp, cw)
		h := pct(arc.Hp, cw)
		if sw == 0 {
			sw = 2.0
		}
		doarc(doc, x, y, w/2, h/2, arc.A1, arc.A2, sw, arc.Color)
	}
	// line
	for _, line := range slide.Line {
		if line.Color == "" {
			line.Color = defaultColor
		}
		setopacity(doc, line.Opacity)
		x1, y1, sw := dimen(cw, ch, line.Xp1, line.Yp1, line.Sp)
		x2, y2, _ := dimen(cw, ch, line.Xp2, line.Yp2, 0)
		if sw == 0 {
			sw = 2.0
		}
		doline(doc, x1, y1, x2, y2, sw, line.Color)
	}

	// for every text element...
	for _, t := range slide.Text {
		if t.Color == "" {
			t.Color = slide.Fg
		}
		if t.Font == "" {
			t.Font = "sans"
		}
		setopacity(doc, t.Opacity)
		x, y, fs = dimen(cw, ch, t.Xp, t.Yp, t.Sp)
		dotext(doc, cw, x, y, fs, t.Wp, t.Tdata, t.Font, t.Color, t.Align, t.Type)
	}
	// for every list element...
	for _, l := range slide.List {
		if l.Color == "" {
			l.Color = slide.Fg
		}
		setopacity(doc, l.Opacity)
		x, y, fs = dimen(cw, ch, l.Xp, l.Yp, l.Sp)
		dolist(doc, x, y, fs, l.Li, l.Font, l.Color, l.Type)
	}
}

// dodeck turns deck input files into PDFs
// if the sflag is set, all output goes to the standard output file,
// otherwise, PDFs are written the destination directory, to filenames based on the input name.
func dodeck(files []string, sflag bool, pagesize, outdir, fontdir, author, title string, gp float64) {
	p, ok := pagemap[pagesize]
	if !ok {
		p = pagemap["Letter"]
	}
	cw := int(p.width * p.unit)
	ch := int(p.height * p.unit)
	if sflag { // combined output to standard output
		doc := gofpdf.New("L", "pt", pagesize, fontdir)
		for _, filename := range files {
			doslides(doc, filename, author, title, cw, ch, gp)
		}
		err := doc.Output(os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	} else { // output to individual files
		for _, filename := range files {
			base := strings.Split(filepath.Base(filename), ".xml")
			out, err := os.Create(filepath.Join(outdir, base[0]+".pdf"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "pdfdeck: %v\n", err)
				continue
			}
			doc := gofpdf.New("L", "pt", pagesize, fontdir)
			doslides(doc, filename, author, title, cw, ch, gp)
			err = doc.Output(out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pdfdeck: %v\n", err)
				continue
			}
			out.Close()
		}
	}
}

// for every file, make a deck
func main() {
	var (
		sansfont = flag.String("sans", "helvetica", "sans font")
		serifont = flag.String("serif", "times", "serif font")
		monofont = flag.String("mono", "courier", "mono font")
		pagesize = flag.String("pagesize", "Letter", "pagesize (Letter, Legal, A3, A4, A5)")
		fontdir  = flag.String("fontdir", filepath.Join(os.Getenv("GOPATH"), "src/code.google.com/p/gofpdf/font"), "directory for fonts")
		outdir   = flag.String("outdir", ".", "output directory")
		stdout   = flag.Bool("stdout", false, "output to standard output")
		title    = flag.String("title", "", "document title")
		author   = flag.String("author", "", "document author")
		gridpct  = flag.Float64("grid", 0, "place percentage grid on each slide")
	)
	flag.Parse()
	fontmap["sans"] = *sansfont
	fontmap["serif"] = *serifont
	fontmap["mono"] = *monofont
	dodeck(flag.Args(), *stdout, *pagesize, *outdir, *fontdir, *author, *title, *gridpct)
}
