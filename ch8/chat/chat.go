// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

//!+broadcaster
// type client chan<- string // an outgoing message channel

const clientBuffer = 5

type client struct {
	messageCh chan<- string
	name      string
}

type clientMessage struct {
	client  client
	message string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan clientMessage) // all incoming client messages
)

func formatClients(clients map[client]bool) string {
	w := strings.Builder{}
	w.WriteString("address\n")
	w.WriteString("----\n")

	for c, v := range clients {
		if !v {
			continue
		}

		w.WriteString(c.name)
		w.WriteString("\n")
	}

	return w.String()
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				if msg.client == cli {
					continue
				}
				go func(c client) { c.messageCh <- msg.message }(cli)
			}

		case cli := <-entering:
			clients[cli] = true
			go func() {
				messages <- clientMessage{cli, formatClients(clients)}
			}()
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.messageCh)
			go func() {
				messages <- clientMessage{cli, formatClients(clients)}
			}()
		}
	}
}

func clientReader(conn net.Conn) <-chan string {
	messages := make(chan string)
	returnChan := make(chan string)

	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			select {
			case messages <- input.Text():
			default:
			}
		}
		close(messages)
	}()

	go func() {
		for {
			select {
			case m, ok := <-messages:
				if !ok {
					return
				}
				returnChan <- m
			case <-time.After(5 * time.Minute):
				conn.Write([]byte("you are boring!\n"))
				conn.Close()
			}
		}

	}()

	return returnChan
}

func readName(inCh <-chan string, outCh chan<- string) string {
	outCh <- "Type your name:"

	return <-inCh
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	ch := make(chan string, clientBuffer) // outgoing client messages
	go clientWriter(conn, ch)

	reader := clientReader(conn)

	name := readName(reader, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- clientMessage{client{ch, who}, who + " has arrived"}

	c := client{messageCh: ch, name: who}
	entering <- c

	for m := range reader {
		messages <- clientMessage{client{ch, who}, name + ": " + m}
	}

	leaving <- c
	messages <- clientMessage{client{ch, who}, who + " has left"}
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
