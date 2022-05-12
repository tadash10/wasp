package main

import (
	"fmt" //nolint:gofumpt
	"github.com/99designs/keyring"
)

func main() {
	ring, _ := keyring.Open(keyring.Config{
		ServiceName: "IOTA Foundation/wasp-cli",
	})

	_ = ring.Set(keyring.Item{
		Key:  "foo",
		Data: []byte("secret-bar"),
	})

	i, _ := ring.Get("foo")

	fmt.Printf("%s", i.Data)
}
