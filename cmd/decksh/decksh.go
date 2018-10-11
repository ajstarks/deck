// decksh: a little language that generates deck markup
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/scanner"
)

// emap is the id=expression map
var emap = map[string]string{}

// xmlmap defines the XML substitutions
var xmlmap = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;")

var shell = "/bin/sh"
var shellflag = "-c"

// on init, set the shell info, if on windows.
func init() {
	if runtime.GOOS == "windows" {
		shell = "cmd"
		shellflag = "/c"
	}
}

// xmlesc escapes XML
func xmlesc(s string) string {
	return xmlmap.Replace(s)
}

// assign creates an assignment by filling in the global id map
func assign(s []string, linenumber int) error {
	if len(s) < 3 {
		return fmt.Errorf("line %d: assignment needs id=<expression>", linenumber)
	}
	emap[s[0]] = s[2]
	return nil
}

// eval evaluates an id=<expressipn> string
func eval(s string) string {
	v, ok := emap[s]
	if ok {
		return v
	}
	return s
}

// parse takes a line of input and returns a string slice containing the parsed tokens
func parse(src string) []string {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	tokens := []string{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		tokens = append(tokens, s.TokenText())
	}
	for i := 1; i < len(tokens); i++ {
		tokens[i] = eval(tokens[i])
	}
	return tokens
}

// dumptokens show the parsed tokens
func dumptokens(w io.Writer, s []string, linenumber int) {
	fmt.Fprintf(w, "line %d: args=%d [ ", linenumber, len(s))
	for i, t := range s {
		fmt.Fprintf(w, "%d:%s ", i, t)
	}
	fmt.Fprintln(w, "]")
}

// deck produces the "deck" element
func deck(w io.Writer, s []string, linenumber int) error {
	_, err := fmt.Fprintln(w, "<deck>")
	return err
}

// canvas produces the "canvas" element
func canvas(w io.Writer, s []string, linenumber int) error {
	e := fmt.Errorf("line %d: canvas width height", linenumber)
	if len(s) != 3 {
		return e
	}
	for i := 1; i < 3; i++ {
		s[i] = eval(s[i])
	}
	fmt.Fprintf(w, "<canvas width=%q height=%q/>\n", s[1], s[2])
	return nil
}

// slide produces the "slide" element
func slide(w io.Writer, s []string, linenumber int) error {
	switch len(s) {
	case 1:
		fmt.Fprintln(w, "<slide>")
	case 2:
		fmt.Fprintf(w, "<slide bg=%s>\n", s[1])
	case 3:
		fmt.Fprintf(w, "<slide bg=%s fg=%s>\n", s[1], s[2])
	default:
		return fmt.Errorf("line %d: slide [bgcolor] [fgcolor]", linenumber)
	}
	return nil
}

// elist ends a deck, slide, or list
func endtag(w io.Writer, s []string, linenumber int) error {
	tag := s[0]
	if len(tag) < 2 || tag[0:1] != "e" {
		return fmt.Errorf("line %d: edeck, eslide, or elist", linenumber)
	}
	fmt.Fprintf(w, "</%s>\n", tag[1:])
	return nil
}

// fontColorOp generates markup for font, color, and opacity
func fontColorOp(s []string) string {
	switch len(s) {
	case 1:
		return fmt.Sprintf("font=%s", s[0])
	case 2:
		return fmt.Sprintf("font=%s color=%s", s[0], s[1])
	case 3:
		return fmt.Sprintf("font=%s color=%s opacity=%q", s[0], s[1], s[2])
	default:
		return ""
	}
}

// remove quotes from a string, and XML escape it
func qesc(s string) string {
	if len(s) < 3 {
		return ""
	}
	return (xmlesc(s[1 : len(s)-1]))
}

// text generates markup for text
func text(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	if n < 5 {
		return fmt.Errorf("line %d: %s \"text\" x y size [font] [color] [opacity]", linenumber, s[0])
	}
	fco := fontColorOp(s[5:])
	switch s[0] {
	case "text":
		fmt.Fprintf(w, "<text xp=%q yp=%q sp=%q %s>%s</text>\n", s[2], s[3], s[4], fco, qesc(s[1]))
	case "ctext":
		fmt.Fprintf(w, "<text align=\"c\" xp=%q yp=%q sp=%q %s>%s</text>\n", s[2], s[3], s[4], fco, qesc(s[1]))
	case "etext":
		fmt.Fprintf(w, "<text align=\"e\" xp=%q yp=%q sp=%q %s>%s</text>\n", s[2], s[3], s[4], fco, qesc(s[1]))
	case "textfile":
		fmt.Fprintf(w, "<text file=%s xp=%q yp=%q sp=%q %s/>\n", s[1], s[2], s[3], s[4], fco)
	}
	return nil
}

// text generates markup for a block of text
func textblock(w io.Writer, s []string, linenumber int) error {
	if len(s) < 6 {
		return fmt.Errorf("line %d: %s \"text\" x y width size [font] [color] [opacity]", linenumber, s[0])
	}
	fmt.Fprintf(w, "<text type=\"block\" xp=%q yp=%q wp=%q sp=%q %s>%s</text>\n", s[2], s[3], s[4], s[5], fontColorOp(s[6:]), qesc(s[1]))
	return nil
}

// textcode generates markup for a block of code
func textcode(w io.Writer, s []string, linenumber int) error {
	switch len(s) {
	case 6:
		fmt.Fprintf(w, "<text type=\"code\" file=%s xp=%q yp=%q wp=%q sp=%q/>\n", s[1], s[2], s[3], s[4], s[5])
	case 7:
		fmt.Fprintf(w, "<text type=\"code\" file=%s xp=%q yp=%q wp=%q sp=%q color=%s/>\n", s[1], s[2], s[3], s[4], s[5], s[6])
	default:
		return fmt.Errorf("line %d: %s \"file\" x y width size [color]", linenumber, s[0])
	}
	return nil
}

// image generates markup for images (plain and captioned)
func image(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: [c]image \"image-file\" x y w h [scale] [link]", linenumber)

	switch n {
	case 6:
		fmt.Fprintf(w, "<image name=%s xp=%q yp=%q width=%q height=%q/>\n", s[1], s[2], s[3], s[4], s[5])
	case 7:
		fmt.Fprintf(w, "<image name=%s xp=%q yp=%q width=%q height=%q scale=%q/>\n", s[1], s[2], s[3], s[4], s[5], s[6])
	case 8:
		fmt.Fprintf(w, "<image name=%s xp=%q yp=%q width=%q height=%q scale=%q link=%s/>\n", s[1], s[2], s[3], s[4], s[5], s[6], s[7])
	default:
		return e
	}
	return nil
}

// cimage makes a captioned image
func cimage(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: cimage \"image-file\" \"caption\" x y w h [scale] [link]", linenumber)
	if n < 6 {
		return e
	}
	caption := xmlesc(s[2])
	switch n {
	case 7:
		fmt.Fprintf(w, "<image name=%s caption=%s xp=%q yp=%q width=%q height=%q/>\n", s[1], caption, s[3], s[4], s[5], s[6])
	case 8:
		fmt.Fprintf(w, "<image name=%s caption=%s xp=%q yp=%q width=%q height=%q scale=%q/>\n", s[1], caption, s[3], s[4], s[5], s[6], s[7])
	case 9:
		fmt.Fprintf(w, "<image name=%s caption=%s xp=%q yp=%q width=%q height=%q scale=%q link=%s/>\n", s[1], caption, s[3], s[4], s[5], s[6], s[7], s[8])
	default:
		return e
	}
	return nil
}

// list generates markup for lists
func list(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	if n < 4 {
		return fmt.Errorf("line %d: %s x y size [font] [color] [opacity]", linenumber, s[0])
	}
	var fco string
	if n > 4 {
		fco = fontColorOp(s[4:])
	}
	switch s[0] {
	case "list":
		fmt.Fprintf(w, "<list xp=%q yp=%q sp=%q %s>\n", s[1], s[2], s[3], fco)
	case "blist":
		fmt.Fprintf(w, "<list type=\"bullet\" xp=%q yp=%q sp=%q %s>\n", s[1], s[2], s[3], fco)
	case "nlist":
		fmt.Fprintf(w, "<list type=\"number\" xp=%q yp=%q sp=%q %s>\n", s[1], s[2], s[3], fco)
	}
	return nil
}

// listitem generates list items
func listitem(w io.Writer, s []string, linenumber int) error {
	if len(s) > 1 {
		fmt.Fprintf(w, "<li>%s</li>\n", qesc(s[1]))
	} else {
		fmt.Fprintln(w, "<li/>")
	}
	return nil
}

// shapes generates markup for rectangle and ellipse
func shapes(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: %s x y w h [color] [opacity]", linenumber, s[0])
	if n < 5 {
		return e
	}
	dim := fmt.Sprintf("xp=%q yp=%q wp=%q hp=%q", s[1], s[2], s[3], s[4])
	switch n {
	case 5:
		fmt.Fprintf(w, "<%s %s/>\n", s[0], dim)
	case 6:
		fmt.Fprintf(w, "<%s %s color=%s/>\n", s[0], dim, s[5])
	case 7:
		fmt.Fprintf(w, "<%s %s color=%s opacity=%q/>\n", s[0], dim, s[5], s[6])
	default:
		return e
	}
	return nil
}

// regshapes generates markup for square and circle
func regshapes(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: %s x y w [color] [opacity]", linenumber, s[0])
	if n < 4 {
		return e
	}
	switch s[0] {
	case "square":
		s[0] = "rect"
	case "circle":
		s[0] = "ellipse"
	}
	dim := fmt.Sprintf("xp=%q yp=%q wp=%q hr=\"100\"", s[1], s[2], s[3])
	switch n {
	case 4:
		fmt.Fprintf(w, "<%s %s/>\n", s[0], dim)
	case 5:
		fmt.Fprintf(w, "<%s %s color=%s/>\n", s[0], dim, s[4])
	case 6:
		fmt.Fprintf(w, "<%s %s color=%s opacity=%q/>\n", s[0], dim, s[4], s[5])
	default:
		return e
	}
	return nil
}

// polygon generates markup for polygons
func polygon(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: %s \"xcoord\" \"ycoord\" [color] [opacity]", linenumber, s[0])
	if n < 3 {
		return e
	}
	switch n {
	case 3:
		fmt.Fprintf(w, "<%s xc=%s yc=%s/>\n", s[0], s[1], s[2])
	case 4:
		fmt.Fprintf(w, "<%s xc=%s yc=%s color=%s/>\n", s[0], s[1], s[2], s[3])
	case 5:
		fmt.Fprintf(w, "<%s xc=%s yc=%s color=%s opacity=%q/>\n", s[0], s[1], s[2], s[3], s[4])
	default:
		return e
	}
	return nil
}

// line generates markup for lines
func line(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: %s x1 y1 x2 y2 [size] [color] [opacity]", linenumber, s[0])
	if n < 5 {
		return e
	}
	lc := fmt.Sprintf("xp1=%q yp1=%q xp2=%q yp2=%q", s[1], s[2], s[3], s[4])
	switch n {
	case 5:
		fmt.Fprintf(w, "<%s %s/>\n", s[0], lc)
	case 6:
		fmt.Fprintf(w, "<%s %s sp=%q/>\n", s[0], lc, s[5])
	case 7:
		fmt.Fprintf(w, "<%s %s sp=%q color=%s/>\n", s[0], lc, s[5], s[6])
	case 8:
		fmt.Fprintf(w, "<%s %s sp=%q color=%s opacity=%q/>\n", s[0], lc, s[5], s[6], s[7])
	default:
		return e
	}
	return nil
}

// arc makes the markup for arc
func arc(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: %s cx cy w h a1 a2 [size] [color] [opacity]", linenumber, s[0])
	if n < 7 {
		return e
	}
	ac := fmt.Sprintf("xp=%q yp=%q wp=%q hp=%q a1=%q a2=%q", s[1], s[2], s[3], s[4], s[5], s[6])
	switch n {
	case 7:
		fmt.Fprintf(w, "<%s %s/>\n", s[0], ac)
	case 8:
		fmt.Fprintf(w, "<%s %s sp=%q/>\n", s[0], ac, s[7])
	case 9:
		fmt.Fprintf(w, "<%s %s sp=%q color=%s/>\n", s[0], ac, s[7], s[8])
	case 10:
		fmt.Fprintf(w, "<%s %s sp=%q color=%s opacity=%q/>\n", s[0], ac, s[7], s[8], s[9])
	default:
		return e
	}
	return nil
}

// curve make quadratic Bezier curve
func curve(w io.Writer, s []string, linenumber int) error {
	n := len(s)
	e := fmt.Errorf("line %d: %s x1 y1 x2 y2 x3 y3 [size] [color] [opacity]", linenumber, s[0])
	if n < 7 {
		return e
	}
	ac := fmt.Sprintf("xp1=%q yp1=%q xp2=%q yp2=%q xp3=%q yp3=%q", s[1], s[2], s[3], s[4], s[5], s[6])
	switch n {
	case 7:
		fmt.Fprintf(w, "<%s %s/>\n", s[0], ac)
	case 8:
		fmt.Fprintf(w, "<%s %s sp=%q/>\n", s[0], ac, s[7])
	case 9:
		fmt.Fprintf(w, "<%s %s sp=%q color=%s/>\n", s[0], ac, s[7], s[8])
	case 10:
		fmt.Fprintf(w, "<%s %s sp=%q color=%s opacity=%q/>\n", s[0], ac, s[7], s[8], s[9])
	default:
		return e
	}
	return nil
}

// chart runs the chart command
func chart(w io.Writer, s string, linenumber int) error {
	// copy the command line into fields, evaluating as we go
	args := strings.Fields(s)
	for i := 1; i < len(args); i++ {
		args[i] = eval(args[i])
		// unquote substituted strings
		la := len(args[i])
		if la > 2 && args[i][0] == '"' && args[i][la-1] == '"' {
			args[i] = args[i][1 : la-1]
		}
	}
	// glue the arguments back into a single string
	s = args[0]
	for i := 1; i < len(args); i++ {
		s = s + " " + args[i]
	}
	out, err := exec.Command(shell, shellflag, s).Output()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", out)
	return err
}

// forloop processes for loops:
// for x=begin end [incr]
// ...
// efor
func forloop(w io.Writer, s []string, linenumber int) (float64, float64, float64, error) {
	var begin, end, incr float64
	if len(s) < 5 {
		return 0, 0, 0, fmt.Errorf("line %d: for begin end [incr] ... efor", linenumber)
	}
	fmt.Sscanf(s[3], "%f", &begin)
	fmt.Sscanf(s[4], "%f", &end)
	incr = 1.0
	if len(s) > 5 {
		fmt.Sscanf(s[5], "%f", &incr)
	}
	return begin, end, incr, nil
}

func parsefor(w io.Writer, s []string, linenumber int, scanner *bufio.Scanner) error {
	begin, end, incr, err := forloop(w, s, linenumber)
	if err != nil {
		return err
	}
	forvar := s[1]
	fmt.Fprintf(os.Stderr, "begin=%v end=%v incr=%v\n", begin, end, incr)
	for scanner.Scan() {
		t := scanner.Text()
		s = parse(t)
		if s[0] == "efor" {
			break
		}
		evaloop(w, forvar, s, begin, end, incr, scanner, linenumber)
	}
	return err
}

func evaloop(w io.Writer, forvar string, s []string, begin, end, incr float64, scanner *bufio.Scanner, linenumber int) {
		cp := make([]string, len(s))
		for i:=0; i < len(s); i++ {
			cp[i] = s[i]
		}
		
		for v := begin; v <= end; v += incr {
			fmt.Fprintf(os.Stderr, "forval=%v -> (%v)\t%v\n", forvar, v, s)
			emap[forvar] = fmt.Sprintf("%v", v)
			for i := 1; i < len(s); i++ {
				cp[i] = eval(s[i])
			}
			keyparse(w, cp, "", scanner, linenumber)
		}
}

// keyparse parses keywords and executes
func keyparse(w io.Writer, tokens []string, t string, sc *bufio.Scanner, n int) error {
	switch tokens[0] {
	case "deck":
		return deck(w, tokens, n)

	case "canvas":
		return canvas(w, tokens, n)

	case "slide":
		return slide(w, tokens, n)

	case "text", "ctext", "etext", "textfile":
		return text(w, tokens, n)

	case "textblock":
		return textblock(w, tokens, n)

	case "textcode":
		return textcode(w, tokens, n)

	case "image":
		return image(w, tokens, n)

	case "cimage":
		return cimage(w, tokens, n)

	case "list", "blist", "nlist":
		return list(w, tokens, n)

	case "elist", "eslide", "edeck":
		return endtag(w, tokens, n)

	case "li":
		return listitem(w, tokens, n)

	case "ellipse", "rect":
		return shapes(w, tokens, n)

	case "circle", "square":
		return regshapes(w, tokens, n)

	case "polygon", "poly":
		return polygon(w, tokens, n)

	case "line":
		return line(w, tokens, n)

	case "arc":
		return arc(w, tokens, n)

	case "curve":
		return curve(w, tokens, n)

	case "dchart", "chart":
		return chart(w, t, n)

	default:
		if len(tokens) > 1 && tokens[1] == "=" {
			return assign(tokens, n)
		}
	}
	return nil
}

// process reads input, parses, dispatches functions for code generation
func process(w io.Writer, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	errors := []error{}

	// For every line in the input, parse into tokens,
	// call the appropriate function, collecting errors as we go.
	// If any errors occurred, print them at the end, and return the latest
	for n := 1; scanner.Scan(); n++ {
		t := scanner.Text()
		tokens := parse(t)
		if len(tokens) < 1 || t[0] == '#' {
			continue
		}
		if tokens[0] == "for" {
			errors = append(errors, parsefor(w, tokens, n, scanner))
		}
		errors = append(errors, keyparse(w, tokens, t, scanner, n))
	}
	// report any collected errors
	nerrs := 0
	for _, e := range errors {
		if e != nil {
			nerrs++
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
	}

	// handle read errors from scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	// return the latest error
	if nerrs > 0 {
		return errors[nerrs-1]
	}

	// all is well, no errors
	return nil
}

// $ decksh                   # input from stdin, output to stdout
// $ decksh -o foo.xml        # input from stdin, output to foo.xml
// $ decksh foo.sh            # input from foo.sh output to stdout
// $ decksh -o foo.xml foo.sh # input from foo.sh output to foo.xml
func main() {
	var dest = flag.String("o", "", "output destination")
	var input io.ReadCloser = os.Stdin
	var output io.WriteCloser = os.Stdout
	var rerr, werr error

	flag.Parse()

	if len(flag.Args()) > 0 {
		input, rerr = os.Open(flag.Args()[0])
		if rerr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", rerr)
			os.Exit(1)
		}
	}

	if len(*dest) > 0 {
		output, werr = os.Create(*dest)
		if werr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", werr)
			os.Exit(2)
		}
	}

	err := process(output, input)
	if err != nil {
		os.Exit(3)
	}

	input.Close()
	output.Close()
	os.Exit(0)
}
