// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package intset

import (
	"fmt"
	"testing"
)

func Example_one() {
	//!+main
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"

	fmt.Println(x.Has(9), x.Has(123)) // "true false"
	//!-main

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
}

func Example_two() {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(42)

	//!+note
	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536]}"
	//!-note

	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536]}
}

func TestLen(t *testing.T) {
	var x IntSet

	if l := x.Len(); l != 0 {
		t.Errorf("wrong lengt: got %d", l)
	}

	x.Add(1)
	x.Add(2)

	if l := x.Len(); l != 2 {
		t.Errorf("wrong length: got %d", l)
	}
}

func TestRemove(t *testing.T) {
	var x IntSet

	x.Add(1)
	x.Add(2)

	x.Remove(2)
	if l := x.Len(); l != 1 {
		t.Errorf("wrong length: got %d", l)
	}

	x.Remove(1)
	if l := x.Len(); l != 0 {
		t.Errorf("wrong length: got %d", l)
	}
}

func Example_Remove() {
	//!+main
	var x IntSet

	for i := 0; i < 5; i++ {
		x.Add(i)
		fmt.Println(&x)
		j := i
		defer func() {
			x.Remove(j)
			fmt.Println(&x)
		}()
	}
	//!-main

	// Output:
	// {0}
	// {0 1}
	// {0 1 2}
	// {0 1 2 3}
	// {0 1 2 3 4}
	// {0 1 2 3}
	// {0 1 2}
	// {0 1}
	// {0}
	// {}
}

func TestClear(t *testing.T) {
	var x IntSet

	for i := 0; i < 100; i++ {
		x.Add(i)
	}

	x.Clear()

	for i, word := range x.words {
		if word != 0 {
			t.Errorf("word %d is not 0", i)
		}
	}
}

func TestEquals(t *testing.T) {
	var x, y IntSet

	x.Add(1)
	x.Add(2)
	y.Add(1)
	y.Add(2)

	if !x.Equals(&y) {
		t.Fail()
	}

	y.Remove(2)
	if x.Equals(&y) {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	var x IntSet

	for i := 0; i < 20; i++ {
		x.Add(i)
	}

	newX := x.Copy()

	if !x.Equals(newX) {
		t.Fail()
	}
}

func Example_IntersectionWith() {
	var x, y, z IntSet

	x.Add(0)
	x.Add(1)
	y.Add(1)
	y.Add(2)

	fmt.Println(x.IntersectionWith(&y))
	fmt.Println(z.IntersectionWith(&x))

	// Output:
	// {1}
	// {}
}

func Example_DifferenceWith() {
	var x, y, z IntSet

	x.Add(0)
	x.Add(1)
	y.Add(1)
	y.Add(2)

	fmt.Println(x.DifferenceWith(&y))
	fmt.Println(x.DifferenceWith(&z))

	z.AddAll(3, 4, 5)
	fmt.Println(x.DifferenceWith(&z))

	// Output:
	// {0}
	// {0 1}
	// {0 1}
}

func Example_SymmetricDifference() {
	var x, y IntSet

	fmt.Println(x.SymmetricDifference(&y))

	x.AddAll(1, 2)
	y.AddAll(2, 3)

	fmt.Println(x.SymmetricDifference(&y))

	// Output:
	// {}
	// {1 3}
}
