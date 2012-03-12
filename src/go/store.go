// Ampstore.go : Draft version of the Ampify data store
package main

import (
	"errors"
	"sync"
	//"time"
)

type Store interface {
	Set(key, value *string) error
	Get(key, value *string) error
}

type KVStore struct {
	//keyList *SkipList
	mu    sync.RWMutex
	table map[string][]byte
	//  save chan record
}

type record struct {
	Key, Value string
}

func (s *KVStore) Get(k *string, v *[]byte) error {
	//  s.mu.RLock()
	//  defer s.mu.Unlock()
	if val, ok := s.table[*k]; ok {
		*v = val
		return nil
	}
	return errors.New("key not found")
}

func (s *KVStore) Set(k *string, v *[]byte) error {
	//  s.mu.Lock()
	//  defer s.mu.Unlock()
	s.table[*k] = *v
	return nil
}

//func (s *KVStore) Scan(start string, end string) {
//  vect = make(Vector)
//  for k,v := range(s.keyList[start, end]) {
//    vect.Append(new Datum{ k, v })
//  }
//  return Serializer.Write(vect)
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
