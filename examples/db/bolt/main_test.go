package main

import (
	"testing"
	"os"
)

func TestMain(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()
}