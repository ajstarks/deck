// pdfdeck: make PDF slide decks from deck markup
package main

import (
	"flag"
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"codeberg.org/go-pdf/fpdf" //"github.com/go-pdf/fpdf"
	"github.com/ajstarks/deck"
	"github.com/mandolyte/mdtopdf"
)

// command line options
type options struct {
	sansfont   string
	serifont   string
	monofont   string
	symbolfont string
	layers     string
	pages      string
	pagesize   string
	fontdir    string
	author     string
	title      string
	outdir     string
	gridpct    float64
	width      int
	height     int
	stdout     bool
	strictwrap bool
}

const (
	mm2pt       = 2.83464 // mm to pt conversion
	linespacing = 1.4
	listspacing = 2.0
	fontfactor  = 1.0
	listwrap    = 95.0
)

type TypedString struct {
	source   string
	datatype string
	data     string
}

// PageDimen describes page dimensions
// the unit field is used to convert to pt.
type PageDimen struct {
	width, height, unit float64
}

// opts are command line options
var opts options

// fontmap maps generic font names to specific implementation names
var fontmap = map[string]string{}

// transmap maps generic font names to the translation function
var transmap = map[string]func(string) string{}

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

// convert tabs to spaces
var codemap = strings.NewReplacer("\t", "    ")

// pagerange returns the begin and end using a "-" string
func pagerange(s string) (int, int) {
	p := strings.Split(s, "-")
	if len(p) != 2 {
		return 0, 0
	}
	b, berr := strconv.Atoi(p[0])
	e, err := strconv.Atoi(p[1])
	if berr != nil || err != nil {
		return 0, 0
	}
	if b > e {
		return 0, 0
	}
	return b, e
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
func setopacity(doc *fpdf.Fpdf, v float64) {
	switch {
	case v < 0:
		doc.SetAlpha(0, "Normal")
	case v > 0:
		doc.SetAlpha(v/100, "Normal")
	case v == 0:
		doc.SetAlpha(1, "Normal")
	}
}

// linesettings set the line style
func linesettings(doc *fpdf.Fpdf) {
	doc.SetLineCapStyle("butt")
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// fontlookup maps font aliases to implementation font names
func fontlookup(s string) string {
	if font, ok := fontmap[s]; ok {
		return font
	}
	return "sans"
}

// grid makes a percentage scale
func grid(doc *fpdf.Fpdf, w, h float64, color string) {
	percent := opts.gridpct
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
func bullet(doc *fpdf.Fpdf, x, y, size float64, color string) {
	rs := size / 2
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Circle(x-size*2, y-rs, rs, "F")
}

// background places a colored rectangle
func background(doc *fpdf.Fpdf, w, h float64, color string) {
	rectangle(doc, 0, 0, w, h, color)
}

// gradientbg  sets the background color gradient
func gradientbg(doc *fpdf.Fpdf, w, h float64, gc1, gc2 string, gp float64) {
	gradient(doc, 0, 0, w, h, gc1, gc2, gp)
}

// gradient  sets the background color gradient
func gradient(doc *fpdf.Fpdf, x, y, w, h float64, gc1, gc2 string, gp float64) {
	r1, g1, b1 := colorlookup(gc1)
	r2, g2, b2 := colorlookup(gc2)
	gp /= 100.0
	doc.LinearGradient(x, y, w, h, r1, g1, b1, r2, g2, b2, 0, gp, 1, 1)
}

// line draws a line
func line(doc *fpdf.Fpdf, xp1, yp1, xp2, yp2, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Line(xp1, yp1, xp2, yp2)
}

// arc draws a line
func arc(doc *fpdf.Fpdf, x, y, w, h, a1, a2, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Arc(x, y, w, h, 0, a1, a2, "D")
}

// quadcurve draws a quadradic bezier curve
func quadcurve(doc *fpdf.Fpdf, xp1, yp1, xp2, yp2, xp3, yp3, sw float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetLineWidth(sw)
	doc.SetDrawColor(r, g, b)
	doc.Curve(xp1, yp1, xp2, yp2, xp3, yp3, "D")
}

// rectangle draws a rectangle
func rectangle(doc *fpdf.Fpdf, x, y, w, h float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Rect(x, y, w, h, "F")
}

// ellipse draws a rectangle
func ellipse(doc *fpdf.Fpdf, x, y, w, h float64, color string) {
	r, g, b := colorlookup(color)
	doc.SetFillColor(r, g, b)
	doc.Ellipse(x, y, w, h, 0, "F")
}

// polygon draws a polygon
func polygon(doc *fpdf.Fpdf, xc, yc, color string, cw, ch float64) {
	xs := strings.Split(xc, " ")
	ys := strings.Split(yc, " ")
	if len(xs) != len(ys) {
		return
	}
	if len(xs) < 3 || len(ys) < 3 {
		return
	}
	poly := make([]fpdf.PointType, len(xs))
	for i := range xs {
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

// content places text elements on the canvas according to type
func textcontent(doc *fpdf.Fpdf, cw, x, y, fs float64, tdata TypedString, t deck.Text) {
	wp := t.Wp
	rotation := t.Rotation
	spacing := t.Lp
	font := t.Font
	align := t.Align
	tlink := t.Link

	if rotation > 0 {
		doc.TransformBegin()
		doc.TransformRotate(rotation, x, y)
	}
	red, green, blue := colorlookup(t.Color)
	doc.SetTextColor(red, green, blue)

	var tw float64
	switch t.Type {
	case "code":
		font = "mono"
		codemap.Replace(tdata.data)
		td := strings.Split(tdata.data, "\n")
		ch := float64(len(td)) * spacing * fs
		tw = deck.Pwidth(wp, cw, cw-x-20)
		rectangle(doc, x-fs, y-fs, tw, ch, "rgb(240,240,240)")
		plaintext(doc, td, x, y, spacing, fs, font, align, tlink)
	case "block":
		tw = deck.Pwidth(wp, cw, cw/2)
		textwrap(doc, x, y, tw, fs, fs*spacing, transmap[font](tdata.data), font, tlink)
	case "markdown":
		markdown(tdata)
	default:
		codemap.Replace(tdata.data)
		td := strings.Split(tdata.data, "\n")
		plaintext(doc, td, x, y, spacing, fs, font, align, tlink)
	}
	if rotation > 0 {
		doc.TransformEnd()
	}
}

// plaintext places lines of text
func plaintext(doc *fpdf.Fpdf, td []string, x, y, spacing, fs float64, font, align, tlink string) {
	ls := spacing * fs
	for _, t := range td {
		showtext(doc, x, y, t, fs, font, align, tlink)
		y += ls
	}
}

// markdown creates a separate PDF from markdown
func markdown(tdata TypedString) {
	pf := mdtopdf.NewPdfRenderer("", "", tdata.source+".pdf", "")
	pf.Process([]byte(tdata.data))
}

// showtext places fully attributed text at the specified location
func showtext(doc *fpdf.Fpdf, x, y float64, s string, fs float64, font, align, link string) {
	offset := 0.0
	doc.SetFont(fontlookup(font), "", fs)
	var t string
	tf, ok := transmap[font]
	if ok {
		t = tf(s)
	} else {
		return
	}
	//t := transmap[font](s)
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

// list places lists on the canvas
func list(doc *fpdf.Fpdf, cw, x, y, fs float64, l deck.List) {
	rotation := l.Rotation
	font := l.Font
	color := l.Color
	align := l.Align
	ltype := l.Type

	if font == "" {
		font = "sans"
	}
	red, green, blue := colorlookup(color)

	if ltype == "bullet" {
		x += fs * 1.2
	}
	ls := l.Lp * fs
	tw := deck.Pwidth(l.Wp, cw, cw/2)

	var t string
	var yw int

	if rotation > 0 {
		doc.TransformBegin()
		doc.TransformRotate(rotation, x, y)
	}
	defont := font
	for i, tl := range l.Li {
		doc.SetFont(fontlookup(font), "", fs)
		doc.SetTextColor(red, green, blue)
		setopacity(doc, tl.Opacity)
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
			font = tl.Font
		} else {
			font = defont
		}
		if align == "center" || align == "c" {
			showtext(doc, x, y, transmap[font](t), fs, font, align, "")
			y += ls
		} else {

			yw = textwrap(doc, x, y, tw, fs, ls, transmap[font](t), font, "")

			y += ls
			if yw >= 1 {
				y += ls * float64(yw)
			}
		}
	}
	if rotation > 0 {
		doc.TransformEnd()
	}
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(doc *fpdf.Fpdf, x, y, w, fs, leading float64, s, font, link string) int {
	var factor = 0.3
	if font == "mono" {
		factor = 1.0
	}
	nbreak := 0
	doc.SetFont(fontlookup(font), "", fs)
	wordspacing := doc.GetStringWidth("M")
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w

	for _, s := range words {
		if s == "\\n" { // magic new line
			xp = x
			yp += (leading * 1.5)
			nbreak++
			continue
		}
		if opts.strictwrap {
			tw := doc.GetStringWidth(s)
			if xp+tw > edge {
				xp = x
				yp += leading
			}
			doc.Text(xp, yp, s)
			xp += tw + (wordspacing * factor)

		} else {
			tw := doc.GetStringWidth(s)
			doc.Text(xp, yp, s)
			xp += tw + (wordspacing * factor)
			if xp > edge {
				xp = x
				yp += leading
				nbreak++
			}
		}
	}
	if len(link) > 0 {
		doc.LinkString(x, y-fs, edge, (yp-y)+fs, link)
	}
	return nbreak
}

// content reads data from a file, returning a tab-expanded string
func filecontent(scheme, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}

// includefile returns the contents of a file as string
func includefile(filetype, filename string) TypedString {
	var ts TypedString
	ts.source = filename
	ts.datatype = filetype
	ts.data = filecontent(filetype, filename)
	return ts
}

// pdfslide makes a slide, one slide per PDF page
func pdfslide(doc *fpdf.Fpdf, d deck.Deck, n int, showslide bool) {
	if n < 0 || n > len(d.Slide)-1 || !showslide {
		return
	}

	var x, y, fs float64
	var imgopt fpdf.ImageOptions
	imgopt.AllowNegativePosition = true

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
		gradientbg(doc, cw, ch, slide.Gradcolor1, slide.Gradcolor2, slide.GradPercent)
	}
	// set the default foreground
	if slide.Fg == "" {
		slide.Fg = "black"
	}

	const defaultColor = "rgb(127,127,127)"
	layerlist := strings.Split(opts.layers, ":")
	// draw elements in the order of the layer list
	for il := range layerlist {
		switch layerlist[il] {
		case "image":
			// for every image on the slide...
			for _, im := range slide.Image {
				x, y, _ = dimen(cw, ch, im.Xp, im.Yp, 0)
				fw, fh := float64(im.Width), float64(im.Height)
				// scale the image by the specified percentage
				if im.Scale > 0 {
					fw *= (im.Scale / 100)
					fh *= (im.Scale / 100)
				}
				// scale the image to fit the canvas width
				if im.Autoscale == "on" && fw > cw {
					fh *= (cw / fw)
					fw = cw
				}
				// scale the image to a percentage of the canvas width
				if im.Height == 0 && im.Width > 0 {
					nw, nh := imageInfo(im.Name)
					if nh > 0 {
						imscale := (fw / 100) * cw
						fw = imscale
						fh = imscale / (float64(nw) / float64(nh))
					}
				}
				midx := fw / 2
				midy := fh / 2
				doc.ImageOptions(im.Name, x-midx, y-midy, fw, fh, false, imgopt, 0, im.Link)
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
		case "rect":
			// rect
			for _, r := range slide.Rect {
				x, y, _ := dimen(cw, ch, r.Xp, r.Yp, 0)
				var w, h float64
				w = pct(r.Wp, cw)
				if r.Hr == 0 {
					h = pct(r.Hp, ch)
				} else {
					h = pct(r.Hr, w)
				}
				if r.Color == "" {
					r.Color = defaultColor
				}
				if len(r.Gradcolor1) > 0 && len(r.Gradcolor2) > 0 {
					gradient(doc, x-(w/2), y-(h/2), w, h, r.Gradcolor1, r.Gradcolor2, r.GradPercent)
				} else {
					setopacity(doc, r.Opacity)
					rectangle(doc, x-(w/2), y-(h/2), w, h, r.Color)
				}
			}
		case "ellipse":
			// ellipse
			for _, e := range slide.Ellipse {
				x, y, _ := dimen(cw, ch, e.Xp, e.Yp, 0)
				var w, h float64
				w = pct(e.Wp, cw)
				if e.Hr == 0 {
					h = pct(e.Hp, ch)
				} else {
					h = pct(e.Hr, w)
				}
				if e.Color == "" {
					e.Color = defaultColor
				}
				setopacity(doc, e.Opacity)
				ellipse(doc, x, y, w/2, h/2, e.Color)
			}
		case "curve":
			// curve
			for _, c := range slide.Curve {
				if c.Color == "" {
					c.Color = defaultColor
				}
				setopacity(doc, c.Opacity)
				x1, y1, sw := dimen(cw, ch, c.Xp1, c.Yp1, c.Sp)
				x2, y2, _ := dimen(cw, ch, c.Xp2, c.Yp2, 0)
				x3, y3, _ := dimen(cw, ch, c.Xp3, c.Yp3, 0)
				if sw == 0 {
					sw = 2.0
				}
				quadcurve(doc, x1, y1, x2, y2, x3, y3, sw, c.Color)
			}
		case "arc":
			// arc
			for _, a := range slide.Arc {
				if a.Color == "" {
					a.Color = defaultColor
				}
				setopacity(doc, a.Opacity)
				x, y, sw := dimen(cw, ch, a.Xp, a.Yp, a.Sp)
				w := pct(a.Wp, cw)
				h := pct(a.Hp, cw)
				if sw == 0 {
					sw = 2.0
				}
				arc(doc, x, y, w/2, h/2, a.A1, a.A2, sw, a.Color)
			}
		case "line":
			// line
			for _, l := range slide.Line {
				if l.Color == "" {
					l.Color = defaultColor
				}
				setopacity(doc, l.Opacity)
				x1, y1, sw := dimen(cw, ch, l.Xp1, l.Yp1, l.Sp)
				x2, y2, _ := dimen(cw, ch, l.Xp2, l.Yp2, 0)
				if sw == 0 {
					sw = 2.0
				}
				line(doc, x1, y1, x2, y2, sw, l.Color)
			}
		case "poly":
			// polygon
			for _, p := range slide.Polygon {
				if p.Color == "" {
					p.Color = defaultColor
				}
				setopacity(doc, p.Opacity)
				polygon(doc, p.XC, p.YC, p.Color, cw, ch)
			}
		case "text":
			// for every text element...
			var tdata TypedString
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
					tdata = includefile(t.Type, t.File)
				} else {
					tdata.data = t.Tdata
				}
				if t.Lp == 0 {
					t.Lp = linespacing
				}
				textcontent(doc, cw, x, y, fs, tdata, t)
			}
		case "list":
			// for every list element...
			for _, l := range slide.List {
				if l.Color == "" {
					l.Color = slide.Fg
				}
				if l.Lp == 0 {
					l.Lp = listspacing
				}
				if l.Wp == 0 {
					l.Wp = listwrap
				}
				setopacity(doc, l.Opacity)
				x, y, fs = dimen(cw, ch, l.Xp, l.Yp, l.Sp)
				list(doc, cw, x, y, fs, l)
			}
		}
	}
	// add a grid, if specified
	if opts.gridpct > 0 {
		grid(doc, cw, ch, slide.Fg)
	}
}

// nulltrans is the null translation function
func nulltrans(s string) string {
	return s
}

// slides reads the deck file, making the PDF version
func slides(doc *fpdf.Fpdf, pc fpdf.InitType, filename string, begin, end int) {
	var d deck.Deck
	var err error

	w := int(pc.Size.Wd)
	h := int(pc.Size.Ht)
	for k, v := range fontmap {
		fontfile := filepath.Join(pc.FontDirStr, v)
		_, err := os.Stat(fontfile + ".json")
		if err != nil {
			doc.AddUTF8Font(v, "", v+".ttf")
			transmap[k] = nulltrans
		} else {
			doc.AddFont(v, "", v+".json")
			transmap[k] = doc.UnicodeTranslatorFromDescriptor("")
		}
	}
	d, err = deck.Read(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pdfdeck: file %q - %v\n", filename, err)
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
	author := opts.author
	title := opts.title
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
	for i := range d.Slide {
		pdfslide(doc, d, i, (i+1 >= begin && i+1 <= end))
	}
}

// pdfdeck turns deck input files into PDFs
// if the sflag is set, all output goes to the standard output file,
// otherwise, PDFs are written the destination directory, to filenames based on the input name.
func pdfdeck(files []string, pageconfig fpdf.InitType, begin, end int) {
	pc := &pageconfig
	if opts.stdout { // combined output to standard output
		doc := fpdf.NewCustom(pc)
		linesettings(doc)
		for _, filename := range files {
			slides(doc, pageconfig, filename, begin, end)
		}
		err := doc.Output(os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pdfdeck: %v\n", err)
		}
		return
	}
	// output to individual files
	for _, filename := range files {
		base := strings.Split(filepath.Base(filename), ".xml")
		out, err := os.Create(filepath.Join(opts.outdir, base[0]+".pdf"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "pdfdeck: file %q - %v\n", filename, err)
			continue
		}
		doc := fpdf.NewCustom(pc)
		linesettings(doc)
		slides(doc, pageconfig, filename, begin, end)
		err = doc.Output(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pdfdeck: file %q - %v\n", filename, err)
			continue
		}
		out.Close()
	}
}

// setpagesize parses the page size string (wxh)
func setpagesize(s string) (float64, float64) {
	var width, height float64
	var err error
	d := strings.FieldsFunc(s, func(c rune) bool { return !unicode.IsNumber(c) })
	if len(d) != 2 {
		return 0, 0
	}
	width, err = strconv.ParseFloat(d[0], 64)
	if err != nil {
		return 0, 0
	}
	height, err = strconv.ParseFloat(d[1], 64)
	if err != nil {
		return 0, 0
	}
	return width, height
}

// setfontdir determines the font directory:
// if the string argument is non-empty, use that, otherwise
// use the contents of the DECKFONT environment variable,
// if that is not set, or empty, use $HOME/deckfonts
func setfontdir(s string) string {
	if len(s) > 0 {
		return s
	}
	envdef := os.Getenv("DECKFONTS")
	if len(envdef) > 0 {
		return envdef
	}
	return path.Join(os.Getenv("HOME"), "deckfonts")
}

// imageInfo returns the dimensions of an image
func imageInfo(s string) (int, int) {
	f, err := os.Open(s)
	defer f.Close()
	if err != nil {
		return 0, 0
	}
	im, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0
	}
	return im.Width, im.Height
}

var usage = `
pdfdeck [options] file...

Options     Default                                            Description
..................................................................................................
-sans       helvetica                                          Sans Serif font
-serif      times                                              Serif font
-mono       courier                                            Monospace font
-symbol     zapfdingbats                                       Symbol font

-layers     image:rect:ellipse:curve:arc:line:poly:text:list   Drawing order
-grid       0                                                  Draw a grid at specified %
-pages      1-1000000                                          Pages to output (first-last)
-pagesize   Letter                                             Page size (w,h) or Letter, Legal,
                                                               Tabloid, A[3-5], ArchA, 4R, Index)

-fontdir    $HOME/deckfonts                                    Font directory
-outdir     Current directory                                  Output directory
-stdout     false                                              Output to standard output
-sw         false                                              Use strict text wrapping
-author     ""                                                 Document author
-title      ""                                                 Document title
....................................................................................................`

func cmdUsage() {
	fmt.Fprintln(flag.CommandLine.Output(), usage)
}

// for every file, make a deck
func main() {
	// process command line
	flag.StringVar(&opts.sansfont, "sans", "helvetica", "sans font")
	flag.StringVar(&opts.serifont, "serif", "times", "serif font")
	flag.StringVar(&opts.monofont, "mono", "courier", "mono font")
	flag.StringVar(&opts.symbolfont, "symbol", "zapfdingbats", "symbol font")
	flag.StringVar(&opts.pagesize, "pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
	flag.StringVar(&opts.fontdir, "fontdir", setfontdir(""), "directory for fonts")
	flag.StringVar(&opts.outdir, "outdir", ".", "output directory")
	flag.StringVar(&opts.title, "title", "", "document title")
	flag.StringVar(&opts.author, "author", "", "document author")
	flag.StringVar(&opts.layers, "layers", "image:rect:ellipse:curve:arc:line:poly:text:list", "Layer order")
	flag.Float64Var(&opts.gridpct, "grid", 0, "draw a percentage grid on each slide")
	flag.StringVar(&opts.pages, "pages", "1-1000000", "page range (first-last)")
	flag.BoolVar(&opts.stdout, "stdout", false, "output to standard output")
	flag.BoolVar(&opts.strictwrap, "sw", false, "strict text wrap")
	flag.Usage = cmdUsage
	flag.Parse()

	// set page dimensions
	pw, ph := setpagesize(opts.pagesize)
	if pw == 0 && ph == 0 {
		p, ok := pagemap[opts.pagesize]
		if !ok {
			p = pagemap["Letter"]
		}
		pw = p.width * p.unit
		ph = p.height * p.unit
	}
	pageconfig := fpdf.InitType{
		UnitStr:    "pt",
		SizeStr:    opts.pagesize,
		Size:       fpdf.SizeType{Wd: pw, Ht: ph},
		FontDirStr: setfontdir(opts.fontdir),
	}

	// set default fonts
	fontmap["sans"] = opts.sansfont
	fontmap["serif"] = opts.serifont
	fontmap["mono"] = opts.monofont
	fontmap["symbol"] = opts.symbolfont

	// make slides
	begin, end := pagerange(opts.pages)
	pdfdeck(flag.Args(), pageconfig, begin, end)
}
