package utils

import (
	"path/filepath"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	var files []string

	root := "./"
	err := filepath.Walk(root, Visit(&files))
	assert.Equal(t, nil, err)

	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}