// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 165.

// Package intset provides a set of integers based on a bit vector.
package intset

import (
	"bytes"
	"fmt"
)

//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string

// Len computes the length of the IntSet
func (s *IntSet) Len() (len int) {
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for i := 0; i < 64; i++ {
			if word&(1<<uint(i)) != 0 {
				len++
			}
		}
	}
	return
}

func (s *IntSet) Remove(x int) {
	if s.Len() == 0 {
		return
	}

	word, bit := x/64, uint64(x%64)
	var mask uint64
	mask--
	allOnes := mask
	mask <<= bit + 1
	mask |= (allOnes >> (64 - bit))
	s.words[word] &= mask
}

func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
}

func (s *IntSet) Copy() *IntSet {
	var newSet IntSet
	for _, oldWord := range s.words {
		newSet.words = append(newSet.words, oldWord)
	}
	return &newSet
}

func (s *IntSet) Equals(another *IntSet) (equals bool) {
	if len(s.words) != len(another.words) {
		return
	}

	for i, word := range s.words {
		if another.words[i] != word {
			return
		}
	}
	equals = true
	return
}

func (s *IntSet) AddAll(values ...int) {
	for _, value := range values {
		s.Add(value)
	}
}

func (s *IntSet) IntersectionWith(another *IntSet) *IntSet {
	var result IntSet

	for i := 0; i < len(s.words) && i < len(another.words); i++ {
		result.words = append(result.words, s.words[i]&another.words[i])
	}

	return &result
}

func (s *IntSet) DifferenceWith(another *IntSet) *IntSet {
	var result IntSet

	for i := 0; i < len(s.words); i++ {
		if i < len(another.words) {
			result.words = append(result.words, s.words[i]&^another.words[i])
		} else {
			result.words = append(result.words, s.words[i])
		}
	}
	return &result
}

func (s *IntSet) SymmetricDifference(another *IntSet) *IntSet {
	result := s.DifferenceWith(another)
	result.UnionWith(another.DifferenceWith(s))
	return result
}
