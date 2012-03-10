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
  fmt.Print("\n")
  // Test Deletion
  if sl.Delete([]byte{5, 2, 3}) {
    fmt.Print("Item deleted!\n")
  }
  if sl.Contains([]byte{5, 2, 3}) {
    fmt.Print("Item located!\n")
  } 
  fmt.Print("\n")
  sl.Show()
  fmt.Print("\n")
  sl.Insert([]byte{9, 2, 3})
  sl.Insert([]byte{2, 2, 3})
  sl.Insert([]byte{6, 2, 3})
  sl.Insert([]byte{7, 2, 3})
  sl.Show()
  fmt.Print("\n")
  fmt.Println(sl.FindElement([]byte{6, 2, 3}).Value)
}
