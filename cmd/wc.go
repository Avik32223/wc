package main

import (
	"flag"

	"github.com/Avik32223/wc/internal/wc"
)

func main() {
	var c, l, w, m bool
	flag.BoolVar(&c, "c", false, "byte count")
	flag.BoolVar(&l, "l", false, "line count")
	flag.BoolVar(&w, "w", false, "word count")
	flag.BoolVar(&m, "m", false, "character count")
	flag.Parse()
	if !c && !l && !w && !m {
		c = true
		w = true
		l = true
	}
	flag.Parse()
	args := flag.Args()
	ccwc := wc.NewWC(args, l, w, m, c)
	ccwc.GetStats()
}
