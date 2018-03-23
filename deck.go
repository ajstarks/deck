// Package deck makes slide decks
package deck

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// Deck defines the structure of a presentation deck
// The size of the canvas, and series of slides
type Deck struct {
	Title       string  `xml:"title"`
	Creator     string  `xml:"creator"`
	Subject     string  `xml:"subject"`
	Publisher   string  `xml:"publisher"`
	Description string  `xml:"description"`
	Date        string  `xml:"date"`
	Canvas      canvas  `xml:"canvas"`
	Slide       []Slide `xml:"slide"`
}

type canvas struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

// Slide is the structure of an individual slide within a deck
// <slide bg="black" fg="rgb(255,255,255)" duration="2s" note="hello, world">
// <slide gradcolor1="black" gradcolor2="white" gp="20" duration="2s" note="wassup">
type Slide struct {
	Bg          string    `xml:"bg,attr"`
	Fg          string    `xml:"fg,attr"`
	Gradcolor1  string    `xml:"gradcolor1,attr"`
	Gradcolor2  string    `xml:"gradcolor2,attr"`
	GradPercent float64   `xml:"gp,attr"`
	Duration    string    `xml:"duration,attr"`
	Note        string    `xml:"note"`
	List        []List    `xml:"list"`
	Text        []Text    `xml:"text"`
	Image       []Image   `xml:"image"`
	Ellipse     []Ellipse `xml:"ellipse"`
	Line        []Line    `xml:"line"`
	Rect        []Rect    `xml:"rect"`
	Curve       []Curve   `xml:"curve"`
	Arc         []Arc     `xml:"arc"`
	Polygon     []Polygon `xml:"polygon"`
}

// CommonAttr are the common attributes for text and list
type CommonAttr struct {
	Xp      float64 `xml:"xp,attr"`      // X coordinate
	Yp      float64 `xml:"yp,attr"`      // Y coordinate
	Sp      float64 `xml:"sp,attr"`      // size
	Lp      float64 `xml:"lp,attr"`      // linespacing (leading) percentage
	Type    string  `xml:"type,attr"`    // type: block, plain, code, number, bullet
	Align   string  `xml:"align,attr"`   // alignment: center, end, begin
	Color   string  `xml:"color,attr"`   // item color
	Opacity float64 `xml:"opacity,attr"` // opacity percentage
	Font    string  `xml:"font,attr"`    // font type: i.e. sans, serif, mono
	Link    string  `xml:"link,attr"`    // reference to other content (i.e. http:// or mailto:)
}

// Dimension describes a graphics object with width and height
type Dimension struct {
	CommonAttr
	Wp float64 `xml:"wp,attr"` // width percentage
	Hp float64 `xml:"hp,attr"` // height percentage
	Hr float64 `xml:"hr,attr"` // height relative percentage
	Hw float64 `xml:"hw,attr"` // height by width
}

// ListItem describes a list item
// <list xp="20" yp="70" sp="1.5">
//	<li>canvas<li>
//	<li>slide</li>
// </list>
type ListItem struct {
	Color    string  `xml:"color,attr"`
	Opacity  float64 `xml:"opacity,attr"`
	Font     string  `xml:"font,attr"`
	ListText string  `xml:",chardata"`
}

// List describes the list element
type List struct {
	CommonAttr
	Wp float64    `xml:"wp,attr"`
	Li []ListItem `xml:"li"`
}

// Text describes the text element
type Text struct {
	CommonAttr
	Wp    float64 `xml:"wp,attr"`
	File  string  `xml:"file,attr"`
	Tdata string  `xml:",chardata"`
}

// Image describes an image
// <image xp="20" yp="30" width="256" height="256" scale="50" name="picture.png" caption="Pretty picture"/>
type Image struct {
	CommonAttr
	Width     int     `xml:"width,attr"`     // image width
	Height    int     `xml:"height,attr"`    // image height
	Scale     float64 `xml:"scale,attr"`     // image scale percentage
	Autoscale string  `xml:"autoscale,attr"` // scale the image to the canvas
	Name      string  `xml:"name,attr"`      // image file name
	Caption   string  `xml:"caption,attr"`   // image caption
}

// Ellipse describes a rectangle with x,y,w,h
// <ellipse xp="45"  yp="10" wp="4" hr="75" color="rgb(0,127,0)"/>
type Ellipse struct {
	Dimension
}

// Rect describes a rectangle with x,y,w,h
// <rect xp="35"  yp="10" wp="4" hp="3"/>
type Rect struct {
	Dimension
}

// Line defines a straight line
// <line xp1="20" yp1="10" xp2="30" yp2="10"/>
type Line struct {
	Xp1     float64 `xml:"xp1,attr"`     // begin x coordinate
	Yp1     float64 `xml:"yp1,attr"`     // begin y coordinate
	Xp2     float64 `xml:"xp2,attr"`     // end x coordinate
	Yp2     float64 `xml:"yp2,attr"`     // end y coordinate
	Sp      float64 `xml:"sp,attr"`      // line thickness
	Color   string  `xml:"color,attr"`   // line color
	Opacity float64 `xml:"opacity,attr"` // line opacity (1-100)
}

// Curve defines a quadratic Bezier curve
// The begining, ending, and control points are required:
// <curve xp1="60" yp1="10" xp2="75" yp2="20" xp3="70" yp3="10" />
type Curve struct {
	Xp1     float64 `xml:"xp1,attr"`
	Yp1     float64 `xml:"yp1,attr"`
	Xp2     float64 `xml:"xp2,attr"`
	Yp2     float64 `xml:"yp2,attr"`
	Xp3     float64 `xml:"xp3,attr"`
	Yp3     float64 `xml:"yp3,attr"`
	Sp      float64 `xml:"sp,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Arc defines an elliptical arc
// the arc is defined by a beginning and ending angle in percentages
// <arc xp="55"  yp="10" wp="4" hr="75" a1="0" a2="180"/>
type Arc struct {
	Dimension
	A1      float64 `xml:"a1,attr"`
	A2      float64 `xml:"a2,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// Polygon defines a polygon, x and y coordinates are specified by
// strings of space-separated percentages:
// <polygon xc="10 20 30" yc="30 40 50"/>
type Polygon struct {
	XC      string  `xml:"xc,attr"`
	YC      string  `xml:"yc,attr"`
	Color   string  `xml:"color,attr"`
	Opacity float64 `xml:"opacity,attr"`
}

// ReadDeck reads the deck description file from a io.Reader
func ReadDeck(r io.ReadCloser, w, h int) (Deck, error) {
	var d Deck
	err := xml.NewDecoder(r).Decode(&d)
	if d.Canvas.Width == 0 {
		d.Canvas.Width = w
	}
	if d.Canvas.Height == 0 {
		d.Canvas.Height = h
	}
	r.Close()
	return d, err
}

// Read reads the deck description file
func Read(filename string, w, h int) (Deck, error) {
	var d Deck
	if filename == "-" {
		return ReadDeck(os.Stdin, w, h)
	}
	r, err := os.Open(filename)
	if err != nil {
		return d, err
	}
	return ReadDeck(r, w, h)
}

// Dimen computes the coordinates and size of an object
func Dimen(c canvas, xp, yp, sp float64) (x, y, s float64) {
	x = (xp / 100) * float64(c.Width)
	y = (yp / 100) * float64(c.Height)
	s = (sp / 100) * float64(c.Width)
	return
}

// Pwidth computes the percent width based on canvas size
func Pwidth(wp, cw, defval float64) float64 {
	if wp == 0 {
		return defval
	}
	return (wp / 100) * cw
}

// Search searches the deck for the specified text, returning the slide number if found
func Search(d Deck, s string) int {
	// for every slide...
	for i, slide := range d.Slide {
		// search lists
		for _, l := range slide.List {
			for _, ll := range l.Li {
				if strings.Contains(ll.ListText, s) {
					return i
				}
			}
		}
		// search text
		for _, t := range slide.Text {
			if strings.Contains(t.Tdata, s) {
				return i
			}
		}
	}
	return -1
}

// Dump shows the decoded description
func Dump(d Deck) {
	fmt.Printf("Title: %#v\nCreator: %#v\nDescription: %#v\nDate: %#v\nPublisher: %#v\nSubject: %#v\n",
		d.Title, d.Creator, d.Description, d.Date, d.Publisher, d.Subject)
	fmt.Printf("Canvas = %v\n", d.Canvas)
	for i, s := range d.Slide {
		fmt.Printf("Slide [%d] = %+v %+v %+v %+v %+v %+v\n", i, s.Bg, s.Fg, s.Duration, s.Gradcolor1, s.Gradcolor2, s.GradPercent)
		for j, l := range s.List {
			fmt.Printf("\tList [%d] = %+v\n", j, l)
		}
		for k, t := range s.Text {
			fmt.Printf("\tText [%d] = %+v\n", k, t)
		}
		for m, im := range s.Image {
			fmt.Printf("\tImage [%d] = %+v\n", m, im)
		}
		for l, line := range s.Line {
			fmt.Printf("\tLine [%d] = %+v\n", l, line)
		}
		for r, rect := range s.Rect {
			fmt.Printf("\tRect [%d] = %+v\n", r, rect)
		}
		for a, arc := range s.Arc {
			fmt.Printf("\tArc [%d] = %+v\n", a, arc)
		}
		for c, curve := range s.Curve {
			fmt.Printf("\tCurve [%d] = %+v\n", c, curve)
		}
		for e, ellipse := range s.Ellipse {
			fmt.Printf("\tEllipse [%d] = %+v\n", e, ellipse)
		}
		for p, polygon := range s.Polygon {
			fmt.Printf("\tPolygon [%d] = %+v\n", p, polygon)
		}
	}
}
