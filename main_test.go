package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultEnvVar(t *testing.T) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	testEnv := []string{"FLATENV___id_hello___there__txt=data:,heyb"}
	root, _ := os.OpenRoot(dir)
	prefix = "FLATENV_"
	mainInRoot(root, func() []string { return testEnv })
	vals, _ := os.ReadFile(filepath.Join(dir, ".id_hello", "there.txt"))
	if string(vals) != "heyb" {
		t.Logf("Unexpected value %s", vals)
		t.Fatal("failed")
	}
}

func TestMainInRoot(t *testing.T) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	testEnv := []string{
		"BOOP___id_hello___there__txt=data:,heyb",
		"BOOP__hi__app=data:application/octet-stream;base64,aGV5YQ==",
		"BOOP__hi__txt=data:text/plain;base64,aGV5YQ=="}
	root, _ := os.OpenRoot(dir)
	prefix = "BOOP_"
	mainInRoot(root, func() []string { return testEnv })
	vals, _ := os.ReadFile(filepath.Join(dir, ".id_hello", "there.txt"))
	if string(vals) != "heyb" {
		t.Logf("Unexpected value %s", vals)
		t.Fatal("failed")
	}
	vals, _ = os.ReadFile(filepath.Join(dir, "_hi.txt"))
	if string(vals) != "heya" {
		t.Logf("Unexpected value %s", vals)
		t.Fatal("failed")
	}
	vals, _ = os.ReadFile(filepath.Join(dir, "_hi.app"))
	if string(vals) != "heya" {
		t.Logf("Unexpected value %s", vals)
		t.Fatal("failed")
	}
}
