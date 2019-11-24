package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	f, err := os.Open("test.dsh")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	process(os.Stdout, f)
	os.Exit(m.Run())
}
