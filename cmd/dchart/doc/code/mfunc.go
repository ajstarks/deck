package main

import (
	"flag"
	"fmt"
	"math"
)

type tfunc struct {
	label    string
	function func(float64) float64
}

func power(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
}

func pow2(n float64) float64 {
	return float64(power(2, int(n)))
}

func linear(x float64) float64 {
	return x
}

func squared(x float64) float64 {
	return x*x
}

func main() {
	fname := flag.String("f", "sine", "function name")
	xrange := flag.String("x", "0,6.283185,0.1", "x range")
	flag.Parse()

	var (
		f     tfunc
		xmin  = 0.0
		xmax  = 2 * math.Pi
		xstep = 0.1
	)
	fmt.Sscanf(*xrange, "%f,%f,%f", &xmin, &xmax, &xstep)
	switch *fname {
	case "square":
		f = tfunc{"y=x*x", squared}
	case "linear":
		f = tfunc{"y=x", linear}
	case "pow":
		f = tfunc{"y=2^n", pow2}
	case "e", "exp":
		f = tfunc{"y=e(x)", math.Exp}
	case "log":
		f = tfunc{"y=log(x)", math.Log10}
	case "sqrt":
		f = tfunc{"y=sqrt(x)", math.Sqrt}
	case "sine", "sin":
		f = tfunc{"y=sin(x)", math.Sin}
	case "cosine", "cos":
		f = tfunc{"y=cos(x)", math.Cos}
	case "sincos":
		f = tfunc{"y=sin(x) * cos(x)",
			func(x float64) float64 { return math.Sin(x) * math.Cos(x) }}
	default:
		f = tfunc{"y=1", func(float64) float64 { return 1 }}
	}
	fmt.Printf("# %s\n", f.label)
	for x := xmin; x <= xmax; x += xstep {
		fmt.Printf("%.2f\t%.4f\n", x, f.function(x))
	}
}
