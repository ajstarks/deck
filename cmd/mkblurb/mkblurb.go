// mkblurb: make a text blurb slide for the deck package
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const slidefmt = "<deck>\n<slide bg=\"%s\" fg=\"%s\">\n<text xp=\"%g\" yp=\"%g\" sp=\"%g\" wp=\"%g\" font=\"%s\" type=\"block\">%s</text>\n</slide>\n</deck>\n"

var xmlrepl *strings.Replacer = strings.NewReplacer(
	"<", "&lt;",
	">", "&gt;",
	"&", "&amp;",
	"'", "&apos;",
	"‘", "&apos;",
	"’", "&apos;",
	"\"", "&quot;",
	"“", "&quot;",
	"”", "&quot;",
	"—", "--")

func xmlquote(s string) string {
	return xmlrepl.Replace(s)
}

func main() {
	var (
		bg    = flag.String("bg", "white", "background color")
		fg    = flag.String("fg", "black", "forground color")
		font  = flag.String("font", "sans", "font")
		size  = flag.Float64("size", 3, "font size %")
		width = flag.Float64("w", 60, "width percent")
		x     = flag.Float64("x", 20, "horizontal percent")
		y     = flag.Float64("y", 70, "vertical percent")
	)
	flag.Parse()

	var output []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		output = append(output, xmlquote(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return
	}
	fmt.Printf(slidefmt, *bg, *fg, *x, *y, *size, *width, *font, output[0])
}
