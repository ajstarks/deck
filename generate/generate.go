// package generate performs slide deck generation
package generate

import (
	"fmt"
	"github.com/ajstarks/deck"
	"io"
)

const (
	circlefmt  = `<ellipse xp="%.2f" yp="%.2f" wp="%.2f" hr="%.2f" opacity="%.2f" color="%s"/>`
	squarefmt  = `<rect xp="%.2f" yp="%.2f" wp="%.2f" hr="%.2f" opacity="%.2f" color="%s"/>`
	ellipsefmt = `<ellipse xp="%.2f" yp="%.2f" wp="%.2f" hp="%.2f" opacity="%.2f" color="%s"/>`
	rectfmt    = `<rect xp="%.2f" yp="%.2f" wp="%.2f" hp="%.2f" opacity="%.2f" color="%s"/>`
	arcfmt     = `<arc xp="%.2f" yp="%.2f" wp="%.2f" hp="%.2f" sp="%.2f" a1="%.2f" a2="%.2f" opacity="%.2f" color="%s"/>`
	linefmt    = `<line xp1="%.2f" yp1="%.2f" xp2="%.2f" yp2="%.2f" sp="%.2f" opacity="%.2f" color="%s"/>`
	curvefmt   = `<curve xp1="%.2f" yp1="%.2f" xp2="%.2f" yp2="%.2f" xp3="%.2f" yp3="%.2f" sp="%.2f" opacity="%.2f" color="%s"/>`
	polygonfmt = `<polygon xc="%s" yc="%s" opacity="%.2f" color="%s"/>`
	textfmt    = `<text xp="%.2f" yp="%.2f" sp="%.2f" align="%s" wp="%.2f" font="%s" opacity="%.2f" color="%s" type="%s">%s</text>`
	imagefmt   = `<image xp="%.2f" yp="%.2f" width="%d" height="%d" name="%s"/>`
	listfmt    = `<list type="%s" xp="%.2f" yp="%.2f" sp="%.2f" font="%s" color="%s">`
	lifmt      = `<li>%s</li>`
	closelist  = `</list>`
	slidefmt   = `<slide>`
	slidebg    = `<slide bg="%s">`
	slidebgfg  = `<slide bg="%s" fg="%s">`
	closeslide = `</slide>`
	deckfmt    = `<deck><canvas width="%d" height="%d"/>`
	closedeck  = `</deck>`
)

// Deck is the generated deck structure.
type Deck struct {
	width, height int
	dest          io.Writer
}

// NewSlides initializes he generated deck structure.
func NewSlides(where io.Writer, w, h int) *Deck {
	return &Deck{dest: where, width: w, height: h}
}

// StartDeck begins a slide deck.
func (p *Deck) StartDeck() {
	fmt.Fprintf(p.dest, deckfmt, p.width, p.height)
}

// EndDeck ends a slide.
func (p *Deck) EndDeck() {
	fmt.Fprintln(p.dest, closedeck)
}

// StartSlide begins a slide.
func (p *Deck) StartSlide(colors ...string) {
	switch len(colors) {
	case 1:
		fmt.Fprintf(p.dest, slidebg, colors[0])
	case 2:
		fmt.Fprintf(p.dest, slidebgfg, colors[0], colors[1])
	default:
		fmt.Fprintln(p.dest, slidefmt)
	}
}

// EndSlide ends a slide.
func (p *Deck) EndSlide() {
	fmt.Fprintln(p.dest, closeslide)
}

// square makes square markup from the rect structure.
func (p *Deck) square(r deck.Rect) {
	fmt.Fprintf(p.dest, squarefmt, r.Xp, r.Yp, r.Wp, r.Hr, r.Opacity, r.Color)
}

// circle makes square markup from the ellipse structure.
func (p *Deck) circle(e deck.Ellipse) {
	fmt.Fprintf(p.dest, circlefmt, e.Xp, e.Yp, e.Wp, e.Hr, e.Opacity, e.Color)
}

// ellipse makes ellipse markup from the ellipse structure.
func (p *Deck) ellipse(e deck.Ellipse) {
	fmt.Fprintf(p.dest, ellipsefmt, e.Xp, e.Yp, e.Wp, e.Hp, e.Opacity, e.Color)
}

// rect makes rect markup rom the rect structure.
func (p *Deck) rect(r deck.Rect) {
	fmt.Fprintf(p.dest, rectfmt, r.Xp, r.Yp, r.Wp, r.Hp, r.Opacity, r.Color)
}

// line makes line markup from the deck line structure.
func (p *Deck) line(l deck.Line) {
	fmt.Fprintf(p.dest, linefmt, l.Xp1, l.Yp1, l.Xp2, l.Yp2, l.Sp, l.Opacity, l.Color)
}

// curve makes curve markup from the curve structure.
func (p *Deck) curve(c deck.Curve) {
	fmt.Fprintf(p.dest, curvefmt, c.Xp1, c.Yp1, c.Xp2, c.Yp2, c.Xp3, c.Yp3, c.Sp, c.Opacity, c.Color)
}

// arc makes arc markup from the arc structure.
func (p *Deck) arc(a deck.Arc) {
	fmt.Fprintf(p.dest, arcfmt, a.Xp, a.Yp, a.Wp, a.Hp, a.Sp, a.A1, a.A2, a.Opacity, a.Color)
}

// polygon makes polygon markup from the polygon structure.
func (p *Deck) polygon(poly deck.Polygon) {
	fmt.Fprintf(p.dest, polygonfmt, poly.XC, poly.YC, poly.Opacity, poly.Color)
}

// text makes text markup from the deck text structure.
func (p *Deck) text(t deck.Text) {
	fmt.Fprintf(p.dest, textfmt, t.Xp, t.Yp, t.Sp, t.Align, t.Wp, t.Font, t.Opacity, t.Color, t.Type, t.Tdata)
}

// image makes image markup from the deck image structure.
func (p *Deck) image(pic deck.Image) {
	fmt.Fprintf(p.dest, imagefmt, pic.Xp, pic.Yp, pic.Width, pic.Height, pic.Name)
}

// list makes markup from the list deck structure.
func (p *Deck) list(l deck.List, items []string, ltype, font, color string) {
	fmt.Fprintf(p.dest, listfmt, ltype, l.Xp, l.Yp, l.Sp, l.Font, l.Color)
	for _, s := range items {
		fmt.Fprintf(p.dest, lifmt, s)
	}
	fmt.Fprintln(p.dest, closelist)
}

// Text places plain text aligned at (x,y), with specified font, size and color. Opacity is optional
func (p *Deck) Text(x, y float64, s, font string, size float64, color string, opacity ...float64) {
	t := deck.Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Color = color
	t.Tdata = s
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextMid places centered text aligned at (x,y), with specified font, size and color. Opacity is optional.
func (p *Deck) TextMid(x, y float64, s, font string, size float64, color string, opacity ...float64) {
	t := deck.Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Tdata = s
	t.Color = color
	t.Align = "center"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextEnd places right-justified text aligned at (x,y), with specified font, size and color. Opacity is optional.
func (p *Deck) TextEnd(x, y float64, s, font string, size float64, color string, opacity ...float64) {
	t := deck.Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Tdata = s
	t.Color = color
	t.Align = "right"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// TextBlock makes a block of text aligned at (x,y), wrapped at margin; with specified font, size and color. Opacity is optional.
func (p *Deck) TextBlock(x, y float64, s, font string, size, margin float64, color string, opacity ...float64) {
	t := deck.Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Font = font
	t.Wp = margin
	t.Tdata = s
	t.Color = color
	t.Type = "block"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// Code makes a code block at (x,y), with specified size and color (opacity is optional),
// on a light gray background with the specified margin width.
func (p *Deck) Code(x, y float64, s string, size, margin float64, color string, opacity ...float64) {
	t := deck.Text{}
	t.Xp = x
	t.Yp = y
	t.Sp = size
	t.Wp = margin
	t.Tdata = s
	t.Color = color
	t.Type = "code"
	if len(opacity) > 0 {
		t.Opacity = opacity[0]
	} else {
		t.Opacity = 100
	}
	p.text(t)
}

// List makes a plain, bullet, or plain list with the specified font, size and color.
func (p *Deck) List(x, y, size float64, items []string, ltype, font, color string) {
	l := deck.List{}
	l.Xp = x
	l.Yp = y
	l.Sp = size
	l.Font = font
	l.Color = color
	p.list(l, items, ltype, font, color)
}

// Square makes a square, centered at (x,y), with width w, at the specified color and optional opacity.
func (p *Deck) Square(x, y, w float64, color string, opacity ...float64) {
	r := deck.Rect{}
	r.Xp = x
	r.Yp = y
	r.Wp = w
	r.Hr = 100
	r.Color = color
	if len(opacity) > 0 {
		r.Opacity = opacity[0]
	} else {
		r.Opacity = 100
	}
	p.square(r)
}

// Circle makes a circle, centered at (x,y) with width w, at the specified color and optional opacity.
func (p *Deck) Circle(x, y, w float64, color string, opacity ...float64) {
	e := deck.Ellipse{}
	e.Xp = x
	e.Yp = y
	e.Wp = w
	e.Hr = 100
	e.Color = color
	if len(opacity) > 0 {
		e.Opacity = opacity[0]
	} else {
		e.Opacity = 100
	}
	p.circle(e)
}

// Rect makes a rectangle, centered at (x,y), with (w,h) dimensions, at the specified color and optional opacity.
func (p *Deck) Rect(x, y, w, h float64, color string, opacity ...float64) {
	r := deck.Rect{}
	r.Xp = x
	r.Yp = y
	r.Wp = w
	r.Hp = h
	r.Color = color
	if len(opacity) > 0 {
		r.Opacity = opacity[0]
	} else {
		r.Opacity = 100
	}
	p.rect(r)
}

// Ellipse makes a ellipse graphic, centered at (x,y), with (w,h) dimensions, at the specified color and optional opacity.
func (p *Deck) Ellipse(x, y, w, h float64, color string, opacity ...float64) {
	e := deck.Ellipse{}
	e.Xp = x
	e.Yp = y
	e.Wp = w
	e.Hp = h
	e.Color = color
	if len(opacity) > 0 {
		e.Opacity = opacity[0]
	} else {
		e.Opacity = 100
	}
	p.ellipse(e)
}

// Line makes a line from (x1,y1) to (x2, y2), with the specified color with optional opacity; thickness is size.
func (p *Deck) Line(x1, y1, x2, y2, size float64, color string, opacity ...float64) {
	l := deck.Line{Xp1: x1, Xp2: x2, Yp1: y1, Yp2: y2, Sp: size, Color: color}
	if len(opacity) > 0 {
		l.Opacity = opacity[0]
	} else {
		l.Opacity = 100
	}
	p.line(l)
}

// Arc makes an arc centered at (x,y), with specified color (with optional opacity),
// with dimensions (w,h), between angle a1 and a2 (specified in degrees).
func (p *Deck) Arc(x, y, w, h, size, a1, a2 float64, color string, opacity ...float64) {
	a := deck.Arc{A1: a1, A2: a2}
	a.Xp = x
	a.Yp = y
	a.Wp = w
	a.Hp = h
	a.Sp = size
	a.Color = color
	if len(opacity) > 0 {
		a.Opacity = opacity[0]
	} else {
		a.Opacity = 100
	}
	p.arc(a)
}

// Curve makes a Bezier curve between (x1, y2) and (x3, y3), with control points at (x2, y2), thickness is specified by size.
func (p *Deck) Curve(x1, y1, x2, y2, x3, y3, size float64, color string, opacity ...float64) {
	c := deck.Curve{Xp1: x1, Xp2: x2, Xp3: x3, Yp1: y1, Yp2: y2, Yp3: y3, Sp: size, Color: color}
	if len(opacity) > 0 {
		c.Opacity = opacity[0]
	} else {
		c.Opacity = 100
	}
	p.curve(c)
}

// Polygon makes a polygon with the specified color (with optional opacity), with coordinates in x and y slices.
func (p *Deck) Polygon(x, y []float64, color string, opacity ...float64) {
	xc, yc := Polycoord(x, y)
	poly := deck.Polygon{XC: xc, YC: yc, Color: color}
	if len(opacity) > 0 {
		poly.Opacity = opacity[0]
	}
	p.polygon(poly)
}

// Polycoord converts slices of coordinates to strings.
func Polycoord(px, py []float64) (string, string) {
	var xc, yc string
	np := len(px)
	if np < 3 || len(py) != np {
		return xc, yc
	}
	for i := 0; i < np-1; i++ {
		xc += fmt.Sprintf("%.2f ", px[i])
		yc += fmt.Sprintf("%.2f ", py[i])
	}
	xc += fmt.Sprintf("%.2f", px[np-1])
	yc += fmt.Sprintf("%.2f", py[np-1])
	return xc, yc
}

// Image places the named image centered at (x, y), with dimensions of (w, h).
func (p *Deck) Image(x, y float64, w, h int, name string) {
	i := deck.Image{Width: w, Height: h, Name: name}
	i.Xp = x
	i.Yp = y
	p.image(i)
}
