package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	r, err := os.Open("test.dsh")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	w, err := os.Create("test.xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	err = process(w, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
	os.Exit(m.Run())
}
