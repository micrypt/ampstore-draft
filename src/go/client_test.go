// Ampstore.go : Draft version of the Ampify data store
package main

import (
	"fmt"
	"testing"
)

func Test_CreatClient(t *testing.T) {
	client := NewClient()
	client.Connect()
	fmt.Println("Connection created")
}
