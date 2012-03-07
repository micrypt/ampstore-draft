// Ampstore.go : Draft version of the Ampify data store
import (
  "./skiplist"
  "./serializer"
)

type Key []byte
type Value []byte

type Store struct {
  keyList *SkipList
  table *[Key] Value
}

type Datum struct {
  k Key
  v Value
}

func (s *Store) Get(k Key) Datum {
  datum := s.table.Get(k)
  return Datum
}

func (s *Store) Set(k Key, v Value) {
  s.keyList.Insert(k)
  s.table.Set(k, v)
}

func (s *Store) Scan(start Key, end Key) {
  vect = make(Vector)
  for k,v := range(s.keyList[start, end]) {
    vect.Append(new Datum{ k, v })
  }
  return Serializer.Write(vect)
}

func (s *Store) Delete(k Key) {
  s.keyList.Delete(k)
  s.table.Delete(k)
}

func Init(opts Opt) *Store {
    store = new Store
    return &store
}
