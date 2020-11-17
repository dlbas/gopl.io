// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 173.

// Bytecounter demonstrates an implementation of io.Writer that counts bytes.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

//!+bytecounter

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

//!-bytecounter

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	words := 0
	advance, t, err := bufio.ScanWords(p, true)
	for len(t) > 0 {
		if err != nil {
			return words, err
		}
		words++
		p = p[advance:]
		advance, t, err = bufio.ScanWords(p, true)
	}
	*c += WordCounter(words)
	return words, nil
}

type LineCounter int

func (c *LineCounter) Write(p []byte) (int, error) {
	lines := 0

	for len(p) > 0 {
		r, size := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			return lines, errors.New("rune error")
		}

		switch r {
		case '\n', '\r':
			lines++
		}
		p = p[size:]
	}

	*c += LineCounter(lines)
	return lines, nil
}

type countingWriter struct {
	w     io.Writer
	count int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	c.count += int64(len(p))
	return c.w.Write(p)
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	newW := &countingWriter{w: w}
	return newW, &newW.count
}

func main() {
	//!+main
	var c ByteCounter
	c.Write([]byte("hello"))
	fmt.Println(c) // "5", = len("hello")

	c = 0 // reset the counter
	var name = "Dolly"
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // "12", = len("hello, Dolly")
	//!-main

	var w WordCounter
	w.Write([]byte("polly wants a cracker"))
	fmt.Println(w)

	var l LineCounter
	l.Write([]byte(`mary
		had
		a
		little
		lamb
	`))
	fmt.Println(l)

	counting, count := CountingWriter(&l)

	fmt.Println(*count)
	counting.Write([]byte("фыдвлаофждыаожфдлофывжалофывджаофывадasdlfkjasdlfkjadslf"))
	fmt.Println(*count)
}
