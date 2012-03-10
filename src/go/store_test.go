// Ampstore.go : Draft version of the Ampify data store
package ampstore 

import (
  "fmt"
  "testing"
)

func Test_Store_Set(t *testing.T) {
  fmt.Println("Begin Store tests")
  store := NewKVStore()
  key := "name"
  value := []byte{1, 2, 3, 4, 5, 6, 7}
  if err := store.Set(&key, &value); err != nil {
    fmt.Println("Error occured")
  }
  fmt.Printf("SET: %v\n", value)
  var val []byte
  store.Get(&key, &val)
  fmt.Printf("GET: %v\n", val)
}
