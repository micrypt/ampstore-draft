// Ampstore.go : Draft version of the Ampify data store
// Skiplist.go: Skiplist indexed by byte slices
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

type Element struct {
	Next  []*Element
	Value []byte
}

type SkipList struct {
	Head   *Element
	Level  int
	Length int
}

// Create new skiplist
func New(maxLevel int) *SkipList {
	head := &Element{make([]*Element, maxLevel), []byte{}}
	return &SkipList{Level: 0, Head: head, Length: 0}
}

// Insert element containing value into skiplist
func (s *SkipList) Insert(value []byte) {
	level := 0
	t := time.Now()
	rand.Seed((int64)(t.Nanosecond()))
	for rand.Float64() < 0.5 {
		level++
	}
	if level+1 > s.Level {
		s.Level = level + 1
	}

	newElement := &Element{make([]*Element, level+1), value}

	for i := s.Level - 1; i >= 0; i-- {
		var ptr *Element
		for ptr = s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			if bytes.Compare(ptr.Next[i].Value, value) > 0 {
				break
			}
		}

		if i <= level {
			newElement.Next[i] = ptr.Next[i]
			ptr.Next[i] = newElement
		}
	}
	s.Length++
}

// Retrieve base list entry containing element for iteration
func (s *SkipList) FindElement(value []byte) *Element {
	for i := s.Level - 1; i >= 0; i-- {
		for ptr := s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			if bytes.Compare(ptr.Next[i].Value, value) == 0 {
				elPtr := ptr.Next[i]
				currLevel := i
				if currLevel > 0 {
					for ; currLevel > 0; currLevel-- {
						if elPtr.Next[1] != nil {
							elPtr = elPtr.Next[1]
						}
					}
				}
				return elPtr
			}
			if bytes.Compare(ptr.Next[i].Value, value) > 0 {
				break
			}
		}
	}
	return nil
}

// Delete elements containing value from the skiplist
func (s *SkipList) Delete(value []byte) bool {
	deleted := false
	for i := s.Level - 1; i >= 0; i-- {
		for ptr := s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			if bytes.Compare(ptr.Next[i].Value, value) == 0 {
				ptr.Next[i] = ptr.Next[i].Next[i]
				deleted = true
				fmt.Printf("Deleted!\n")
				break
			}
			if bytes.Compare(ptr.Next[i].Value, value) > 0 {
				fmt.Printf("Next level >>\n")
				break
			}
		}
	}
	return deleted
}

// Test if skiplist contains value
func (s *SkipList) Contains(value []byte) bool {
	for i := s.Level - 1; i >= 0; i-- {
		for ptr := s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			if bytes.Compare(ptr.Next[i].Value, value) == 0 {
				return true
			}
			if bytes.Compare(ptr.Next[i].Value, value) > 0 {
				break
			}
		}
	}
	return false
}

// Retrieve length of skiplist
func (s *SkipList) Len() int {
	return s.Length
}

//  For debugging: Show elements in skiplist
func (s *SkipList) Show() {
	for i := s.Level - 1; i >= 0; i-- {
		fmt.Printf("Level %d: ", i)
		for ptr := s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			fmt.Printf("%v ", ptr.Next[i].Value)
		}
		fmt.Printf("\n\n")
	}
}
