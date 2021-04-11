// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

// Package bank provides a concurrency-safe bank with one account.
package bank

var deposits = make(chan int)                  // send amount to deposit
var balances = make(chan int)                  // receive balance
var withdrawals = make(chan withdrawalMessage) // withdraw balance

type withdrawalMessage struct {
	amount     int
	responseCh chan bool
}

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func Withdraw(amount int) (result bool) {
	m := withdrawalMessage{amount: amount, responseCh: make(chan bool)}
	withdrawals <- m
	return <-m.responseCh
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case w := <-withdrawals:
			if balance-w.amount < 0 {
				w.responseCh <- false // potentially may block if client is unable to read
			} else {
				balance -= w.amount
				w.responseCh <- true // potentially may block if client is unable to read
			}
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
