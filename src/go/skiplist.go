// Ampstore.go : Draft version of the Ampify data store
package ampstore

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

func New(maxLevel int) *SkipList {
	head := &Element{make([]*Element, maxLevel), []byte{}}
	return &SkipList{Level: 0, Head: head, Length: 0}
}

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

func (s *SkipList) Delete(value []byte) bool {
	deleted := false
	for i := s.Level - 1; i >= 0; i-- {
      for ptr := s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			if bytes.Compare(ptr.Next[i].Value, value) == 0 {
				ptr.Next[i] = ptr.Next[i].Next[i]
				deleted = true
                fmt.Printf("Deleted")
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

func (s *SkipList) Len() int {
	return s.Length
}

func (s *SkipList) Show() {
	for i := s.Level - 1; i >= 0; i-- {
		fmt.Printf("Level %d: ", i)
		for ptr := s.Head; ptr.Next[i] != nil; ptr = ptr.Next[i] {
			fmt.Printf("%v ", ptr.Next[i].Value)
		}
		fmt.Printf("\n\n")
	}
}
