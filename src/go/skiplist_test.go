// Ampstore.go : Draft version of the Ampify data store
package ampstore

import (
  "testing"
)

func Test_SkipList(t *testing.T) {
  sl := New(32)
  sl.Insert([]byte{1, 2, 3})
  sl.Insert([]byte{3, 2, 3})
  sl.Insert([]byte{5, 2, 3})
  sl.Insert([]byte{7, 2, 3})
  sl.Show()
}
