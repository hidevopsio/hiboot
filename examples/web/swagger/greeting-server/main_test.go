package main

import (
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	time.Sleep(1)
	go main()
}

