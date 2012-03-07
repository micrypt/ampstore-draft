// Ampstore.go : Draft version of the Ampify data store
package ampstore

import (
  "fmt"
  "testing"
)

func Test_SkipList(t *testing.T) {
  sl := New(32)
  // Test Insertion
  sl.Insert([]byte{1, 2, 3})
  sl.Insert([]byte{3, 2, 3})
  sl.Insert([]byte{5, 2, 3})
  sl.Insert([]byte{4, 2, 3})
  sl.Show()
  // Test Deletion
  if sl.Delete([]byte{5, 2, 3}) {
    fmt.Print("Item deleted!\n")
  }
  if sl.Contains([]byte{5, 2, 3}) {
    fmt.Print("Item located!\n")
  } 
  sl.Show()
}
