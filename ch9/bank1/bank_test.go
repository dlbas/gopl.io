// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package bank_test

import (
	"fmt"
	"testing"

	bank "gopl.io/ch9/bank1"
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		done <- struct{}{}
	}()

	// Wait for both transactions.
	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}

	withdrawResults := make(chan bool)

	go func() {
		res := bank.Withdraw(150)
		done <- struct{}{}
		withdrawResults <- res
	}()

	<-done

	if got, want := bank.Balance(), 150; got != want {
		t.Errorf("Balance = %d, want %d", got, want)
	}

	if got, want := <-withdrawResults, true; got != want {
		t.Errorf("Transaction returned wrong status %v", got)
	}
}
