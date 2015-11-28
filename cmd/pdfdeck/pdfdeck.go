// pdfdeck: make PDF slide decks
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ajstarks/deck"
	"github.com/jung-kurt/gofpdf"
)

const (
	mm2pt       = 2.83464 // mm to pt conversion
	linespacing = 1.4
	listspacing = 2.0
	fontfactor  = 1.0
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
	"Letter":     {792, 612, 1},
	"Legal":      {1008, 612, 1},
	"Tabloid":    {1224, 792, 1},
	"ArchA":      {864, 648, 1},
	"Widescreen": {1152, 648, 1},
	"4R":         {432, 288, 1},
	"Index":      {360, 216, 1},
	"A2":         {420, 594, mm2pt},
	"A3":         {420, 297, mm2pt},
	"A4":         {297, 210, mm2pt},
	"A5":         {210, 148, mm2pt},
}

var codemap = strings.NewReplacer("\t", "    ")

// translate is the function that does unicode character replacement
var translate func(string) string

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
			showtext(doc, x, h-fs, fmt.Sprintf("%.0f", pl), fs, "sans", "center", "")
		}
		pl += percent
	}
	for y, pl := 0.0, 0.0; y <= h; y += ph {
		doc.Line(0, y, w, y)
		if pl < 100 {
			showtext(doc, fs, y+(fs/3), fmt.Sprintf("%.0f", 100-pl), fs, "sans", "center", "")
		}
		pl += percent
	}
}

// bullet draws a bullet
func bullet(doc *gofpdf.Fpdf, x, y, size float64, color string) {
	rs := size / 2
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Circle(x-size*2, y-rs, rs, "F")
	//dorect(doc, x-size, y-rs, rs, rs, color)
}

// background places a colored rectangle
func background(doc *gofpdf.Fpdf, w, h float64, color string) {
	dorect(doc, 0, 0, w, h, color)
}

// gradient sets the background color gradient
func gradient(doc *gofpdf.Fpdf, w, h float64, gc1, gc2 string, gp float64) {
	r1, g1, b1 := colorlookup(gc1)
	r2, g2, b2 := colorlookup(gc2)
	gp /= 100.0
	doc.LinearGradient(0, 0, w, h, r1, g1, b1, r2, g2, b2, 0, gp, 0, 0)
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

// dopoly draws a polygon
func dopoly(doc *gofpdf.Fpdf, xc, yc, color string, cw, ch float64) {
	xs := strings.Split(xc, " ")
	ys := strings.Split(yc, " ")
	if len(xs) != len(ys) {
		return
	}
	if len(xs) < 3 || len(ys) < 3 {
		return
	}
	poly := make([]gofpdf.PointType, len(xs))
	for i := 0; i < len(xs); i++ {
		x, err := strconv.ParseFloat(xs[i], 64)
		if err != nil {
			poly[i].X = 0
		} else {
			poly[i].X = pct(x, cw)
		}
		y, err := strconv.ParseFloat(ys[i], 64)
		if err != nil {
			poly[i].Y = 0
		} else {
			poly[i].Y = pct(100-y, ch)
		}
	}
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Polygon(poly, "F")
}

// dotext places text elements on the canvas according to type
func dotext(doc *gofpdf.Fpdf, cw, x, y, fs, wp, spacing float64, tdata, font, color, align, ttype, tlink string) {
	var tw float64

	td := strings.Split(tdata, "\n")
	red, green, blue := colorlookup(color)
	doc.SetTextColor(red, green, blue)
	if ttype == "code" {
		font = "mono"
		ch := float64(len(td)) * spacing * fs
		tw = deck.Pwidth(wp, cw, cw-x-20)
		dorect(doc, x-fs, y-fs, tw, ch, "rgb(240,240,240)")
	}
	if ttype == "block" {
		tw = deck.Pwidth(wp, cw, cw/2)
		textwrap(doc, x, y, tw, fs, fs*spacing, translate(tdata), font, tlink)
	} else {
		ls := spacing * fs
		for _, t := range td {
			showtext(doc, x, y, t, fs, font, align, tlink)
			y += ls
		}
	}
}

// showtext places fully attributed text at the specified location
func showtext(doc *gofpdf.Fpdf, x, y float64, s string, fs float64, font, align, link string) {
	offset := 0.0
	doc.SetFont(fontlookup(font), "", fs)
	t := translate(s)
	tw := doc.GetStringWidth(t)
	switch align {
	case "center", "middle", "mid", "c":
		offset = (tw / 2)
	case "right", "end", "e":
		offset = tw
	}
	doc.Text(x-offset, y, t)
	if len(link) > 0 {
		doc.LinkString(x-offset, y-fs, tw, fs, link)
	}
}

// dolists places lists on the canvas
func dolist(doc *gofpdf.Fpdf, x, y, fs, spacing float64, list []deck.ListItem, font, color, ltype string) {
	if font == "" {
		font = "sans"
	}
	red, green, blue := colorlookup(color)

	if ltype == "bullet" {
		x += fs * 1.2
	}
	ls := spacing * fs
	var t string
	for i, tl := range list {
		doc.SetFont(fontlookup(font), "", fs)
		doc.SetTextColor(red, green, blue)
		if ltype == "number" {
			t = fmt.Sprintf("%d. ", i+1) + tl.ListText
		} else {
			t = tl.ListText
		}
		if ltype == "bullet" {
			bullet(doc, x, y, fs/2, color)
		}
		if len(tl.Color) > 0 {
			tlred, tlgreen, tlblue := colorlookup(tl.Color)
			doc.SetTextColor(tlred, tlgreen, tlblue)
		}
		if len(tl.Font) > 0 {
			doc.SetFont(fontlookup(tl.Font), "", fs)
		}
		doc.Text(x, y, translate(t))
		y += ls
	}
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *gofpdf.Fpdf, x, y, w, fs, leading float64, s, font, link string) {
	var factor = 0.3
	if font == "mono" {
		factor = 1.0
	}
	doc.SetFont(fontlookup(font), "", fs)
	wordspacing := doc.GetStringWidth("M")
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
	if len(link) > 0 {
		doc.LinkString(x, y-fs, edge, (yp-y)+fs, link)
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
	// set default background
	if slide.Bg == "" {
		slide.Bg = "white"
	}
	background(doc, cw, ch, slide.Bg)

	if slide.GradPercent <= 0 || slide.GradPercent > 100 {
		slide.GradPercent = 100
	}
	// set gradient background, if specified. You need both colors
	if len(slide.Gradcolor1) > 0 && len(slide.Gradcolor2) > 0 {
		gradient(doc, cw, ch, slide.Gradcolor1, slide.Gradcolor2, slide.GradPercent)
	}
	// set the default foreground
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	// for every image on the slide...
	for _, im := range slide.Image {
		x, y, _ = dimen(cw, ch, im.Xp, im.Yp, 0)
		fw, fh := float64(im.Width), float64(im.Height)
		if im.Scale > 0 {
			fw *= (im.Scale / 100)
			fh *= (im.Scale / 100)
		}
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
			showtext(doc, x, y+(midy)+(capsize*1.5), im.Caption, capsize, im.Font, im.Align, "")
		}
	}
	// every graphic on the slide
	const defaultColor = "rgb(127,127,127)"
	// rect
	for _, rect := range slide.Rect {
		x, y, _ := dimen(cw, ch, rect.Xp, rect.Yp, 0)
		var w, h float64
		w = pct(rect.Wp, cw)
		if rect.Hr == 0 {
			h = pct(rect.Hp, ch)
		} else {
			h = pct(rect.Hr, w)
		}
		if rect.Color == "" {
			rect.Color = defaultColor
		}
		setopacity(doc, rect.Opacity)
		dorect(doc, x-(w/2), y-(h/2), w, h, rect.Color)
	}
	// ellipse
	for _, ellipse := range slide.Ellipse {
		x, y, _ := dimen(cw, ch, ellipse.Xp, ellipse.Yp, 0)
		var w, h float64
		w = pct(ellipse.Wp, cw)
		if ellipse.Hr == 0 {
			h = pct(ellipse.Hp, ch)
		} else {
			h = pct(ellipse.Hr, w)
		}
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
	// polygon
	for _, poly := range slide.Polygon {
		if poly.Color == "" {
			poly.Color = defaultColor
		}
		setopacity(doc, poly.Opacity)
		dopoly(doc, poly.XC, poly.YC, poly.Color, cw, ch)
	}

	// for every text element...
	var tdata string
	for _, t := range slide.Text {
		if t.Color == "" {
			t.Color = slide.Fg
		}
		if t.Font == "" {
			t.Font = "sans"
		}
		setopacity(doc, t.Opacity)
		x, y, fs = dimen(cw, ch, t.Xp, t.Yp, t.Sp)
		if t.File != "" {
			tdata = includefile(t.File)
		} else {
			tdata = t.Tdata
		}
		if t.Lp == 0 {
			t.Lp = linespacing
		}
		dotext(doc, cw, x, y, fs, t.Wp, t.Lp, tdata, t.Font, t.Color, t.Align, t.Type, t.Link)
	}
	// for every list element...
	for _, l := range slide.List {
		if l.Color == "" {
			l.Color = slide.Fg
		}
		if l.Lp == 0 {
			l.Lp = listspacing
		}
		setopacity(doc, l.Opacity)
		x, y, fs = dimen(cw, ch, l.Xp, l.Yp, l.Sp)
		dolist(doc, x, y, fs, l.Lp, l.Li, l.Font, l.Color, l.Type)
	}
	// add a grid, if specified
	if gp > 0 {
		grid(doc, cw, ch, slide.Fg, gp)
	}
}

// doslides reads the deck file, making the PDF version
func doslides(doc *gofpdf.Fpdf, pc gofpdf.InitType, filename, author, title string, gp float64) {
	var d deck.Deck
	var err error

	translate = doc.UnicodeTranslatorFromDescriptor("")
	w := int(pc.Size.Wd)
	h := int(pc.Size.Ht)
	for _, v := range fontmap {
		doc.AddFont(v, "", v+".json")
	}
	d, err = deck.Read(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pdfdeck: %v\n", err)
		return
	}
	if pc.OrientationStr == "L" {
		w, h = h, w
	}
	d.Canvas.Width = w
	d.Canvas.Height = h
	doc.SetDisplayMode("fullpage", "single") // optimal set for presentations
	doc.SetCreator("pdfdeck", true)

	// Document-supplied overrides command-line specified metadata
	if len(d.Creator) > 0 {
		author = d.Creator
	}

	if len(d.Title) > 0 {
		title = d.Title
	}

	if len(title) > 0 {
		doc.SetTitle(title, true)
	}

	if len(author) > 0 {
		doc.SetAuthor(author, true)
	}

	if len(d.Subject) > 0 {
		doc.SetSubject(d.Subject, true)
	}
	for i := 0; i < len(d.Slide); i++ {
		pdfslide(doc, d, i, gp)
	}
}

// dodeck turns deck input files into PDFs
// if the sflag is set, all output goes to the standard output file,
// otherwise, PDFs are written the destination directory, to filenames based on the input name.
func dodeck(files []string, pageconfig gofpdf.InitType, w, h float64, sflag bool, outdir, author, title string, gp float64) {
	pc := &pageconfig
	if sflag { // combined output to standard output
		doc := gofpdf.NewCustom(pc)
		for _, filename := range files {
			doslides(doc, pageconfig, filename, author, title, gp)
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
			doc := gofpdf.NewCustom(pc)
			doslides(doc, pageconfig, filename, author, title, gp)
			err = doc.Output(out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pdfdeck: %v\n", err)
				continue
			}
			out.Close()
		}
	}
}

// includefile returns the contents of a file as string
func includefile(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}

// for every file, make a deck
func main() {
	var (
		sansfont   = flag.String("sans", "helvetica", "sans font")
		serifont   = flag.String("serif", "times", "serif font")
		monofont   = flag.String("mono", "courier", "mono font")
		symbolfont = flag.String("symbol", "zapfdingbats", "symbol font")
		pagesize   = flag.String("pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
		fontdir    = flag.String("fontdir", os.Getenv("DECKFONTS"), "directory for fonts (defaults to DECKFONTS environment variable)")
		outdir     = flag.String("outdir", ".", "output directory")
		title      = flag.String("title", "", "document title")
		author     = flag.String("author", "", "document author")
		gridpct    = flag.Float64("grid", 0, "draw a percentage grid on each slide")
		stdout     = flag.Bool("stdout", false, "output to standard output")
	)
	flag.Parse()

	var pw, ph float64
	nd, err := fmt.Sscanf(*pagesize, "%g,%g", &pw, &ph)
	if nd != 2 || err != nil {
		pw, ph = 0.0, 0.0
	}
	if pw == 0 && ph == 0 {
		p, ok := pagemap[*pagesize]
		if !ok {
			p = pagemap["Letter"]
		}
		pw = p.width * p.unit
		ph = p.height * p.unit
	}

	pageconfig := gofpdf.InitType{
		UnitStr:    "pt",
		SizeStr:    *pagesize,
		Size:       gofpdf.SizeType{Wd: pw, Ht: ph},
		FontDirStr: *fontdir,
	}
	fontmap["sans"] = *sansfont
	fontmap["serif"] = *serifont
	fontmap["mono"] = *monofont
	fontmap["symbol"] = *symbolfont
	dodeck(flag.Args(), pageconfig, pw, ph, *stdout, *outdir, *author, *title, *gridpct)
}
