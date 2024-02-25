package wc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"
)

type stat struct {
	filename string
	bytes    uint64
	lines    uint64
	words    uint64
	chars    uint64
}

func (s *stat) addTo(d *stat) {
	s.bytes += d.bytes
	s.lines += d.lines
	s.words += d.words
	s.chars += d.chars
}

type WC struct {
	filenames []string
	printLine bool
	printWord bool
	printChar bool
	printByte bool
}

func NewWC(filenames []string, printLine, printWord, printChar, printByte bool) *WC {
	return &WC{
		filenames: filenames,
		printLine: printLine,
		printWord: printWord,
		printChar: printChar,
		printByte: printByte,
	}
}

func calculateStat(reader *bufio.Reader, filename string) (*stat, error) {
	s := &stat{filename: filename}
	var p rune
	for {
		c, b, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if c == rune(0) && !unicode.IsSpace(p) {
					s.words++
				}
				break
			}
			return nil, err
		}
		if c == '\n' {
			s.lines++
		}

		if unicode.IsSpace(c) && !unicode.IsSpace(p) {
			s.words++
		}

		s.bytes += uint64(b)
		s.chars++
		p = c

	}
	return s, nil
}

func (w *WC) printStat(tw *tabwriter.Writer, s *stat) {
	o := " " + s.filename
	count := 2
	if w.printByte {
		o = fmt.Sprintf("%s%d\t%s", strings.Repeat(" ", count), s.bytes, o)
		count++
	} else if w.printChar {
		o = fmt.Sprintf("%s%d\t%s", strings.Repeat(" ", count), s.chars, o)
		count++
	}
	if w.printWord {
		o = fmt.Sprintf("%s%d\t%s", strings.Repeat(" ", count), s.words, o)
		count++
	}
	if w.printLine {
		o = fmt.Sprintf("%s%d\t%s", strings.Repeat(" ", count), s.lines, o)
	}
	fmt.Fprintln(tw, o)
}

func (w *WC) GetStats() {
	totals := &stat{filename: "total"}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight)
	defer tw.Flush()

	if len(w.filenames) > 0 {
		for _, filename := range w.filenames {
			f, err := os.OpenFile(filename, os.O_RDONLY, 0)
			if err != nil {
				msg := fmt.Sprintf("ccwc: %s", err)
				fmt.Fprintln(tw, msg)
				continue
			}
			defer f.Close()
			reader := bufio.NewReader(f)
			st, err := calculateStat(reader, filename)
			if err != nil {
				fmt.Fprintln(tw, err)
			}
			w.printStat(tw, st)
			totals.addTo(st)
		}

		if len(w.filenames) > 1 {
			w.printStat(tw, totals)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		st, err := calculateStat(reader, "")
		if err != nil {
			fmt.Fprintln(tw, err)
		}
		w.printStat(tw, st)
	}
}
