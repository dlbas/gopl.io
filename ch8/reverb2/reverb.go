// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	idleSeconds = 10 * time.Second
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func scanInput(scanner *bufio.Scanner) chan string {
	ch := make(chan string)

	go func() {
		for scanner.Scan() { ch <- scanner.Text() }
	}()

	return ch
}

func closeConnection(c net.Conn) {
	fmt.Fprintf(c, "closing connection after timeout %s\n", idleSeconds)
	c.Close() // ignoring error
}

func stopAbort(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <- timer.C:
		default:
		}
	}
}

//!+
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)

	text := scanInput(input)
	abort := time.NewTimer(idleSeconds)

	for {
		select {
		case s := <- text:
			go echo(c, s, 1 * time.Second)
			stopAbort(abort) // read abort.Reset() docs on why we should stop the timer before .Reset()
			abort.Reset(idleSeconds)
		case <- abort.C:
			closeConnection(c)
			stopAbort(abort)
			return
		}
	}
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
