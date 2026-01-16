package main

import (
	"os"
	"testing"
)

func TestMainInRoot(*t testing.T) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	testEnv := map[string]string{"BOOP___id_hello___there__txt": "data:,heya"}
	mainInRoot(os.OpenRoot(dir), func() map[string]string { return testEnv })
}
