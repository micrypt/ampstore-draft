// Ampstore.go : Draft version of the Ampify data store
package main

import (
	"errors"
	"sync"
	//"time"
)

type Store interface {
	Set(key string, value *[]byte) error
	Get(key string, value *[]byte) error
}

type KVStore struct {
	//	keyList SkipList
	mu    sync.RWMutex
	table map[string][]byte
	//  save chan record
}

type record struct {
	Key, Value string
}

func (s *KVStore) Get(k string, v *[]byte) error {
	//  s.mu.RLock()
	//  defer s.mu.Unlock()
	//    s.keyList.Insert([]byte(k))
	val, ok := s.table[k]
	if ok {
		*v = val
		return nil
	}
	return errors.New("key not found")
}

func (s *KVStore) Set(k string, v *[]byte, resp *[]byte) error {
	//  s.mu.Lock()
	//  defer s.mu.Unlock()
	s.table[k] = *v
	*resp = []byte(OK_STRING)
	return nil
}

//func (s *KVStore) Scan(start string, end string, resp *[]byte ) {
//  el := s.keyList.FindElement([]byte(start))
//  for ; el.Next[0] != nil {
//    append(resp, el.Value)
//    el = el.Next[0]
//  }
//  return nil
//}

func (s *KVStore) Delete(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.table, k)
}

func NewKVStore() *KVStore {
	kvs := &KVStore{table: make(map[string][]byte)}
	return kvs
}
