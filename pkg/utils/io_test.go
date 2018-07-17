// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
	"github.com/hidevopsio/hiboot/pkg/log"
	"path/filepath"
)

const testPath = "/a/b/c.txt"

func init() {
	log.SetLevel(log.DebugLevel)
}


func TestChangeWorkDir(t *testing.T) {
	wd1 := GetWorkDir()
	ChangeWorkDir("..")
	wd2 := GetWorkDir()
	ChangeWorkDir(wd1)

	assert.NotEqual(t, wd1, wd2)
}

func TestGetWorkingDir(t *testing.T) {
	wd := GetWorkDir()
	expected, err := os.Getwd()
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, wd)
}

func TestGetRelativePath(t *testing.T) {
	path := GetRelativePath(1)

	assert.Equal(t, "utils", DirName(path))
}

func TestIsPathNotExist(t *testing.T) {
	assert.Equal(t, true, IsPathNotExist("/TestNotExistPath"))
}

func TestCreateFile(t *testing.T) {
	err := CreateFile(os.TempDir(), "test.txt")
	assert.Equal(t, nil, err)

	bytesToWrite := "hello"
	b, err := WriterFile(os.TempDir(), "test.txt", []byte(bytesToWrite))
	assert.Equal(t, nil, err)
	assert.Equal(t, len(bytesToWrite), b)
}

func TestWriteFile(t *testing.T) {
	in := "hello, world"
	buf := bytes.NewBufferString(in)
	path := filepath.Join(os.TempDir(), "foo")
	err := os.RemoveAll(path) // remove it first
	assert.Equal(t, nil, err)
	log.Println("path: ", path)
	n, err := WriterFile(path, "test.txt", buf.Bytes())
	assert.Equal(t, nil, err)
	assert.Equal(t, len(in), n)

	n, err = WriterFile("/should-not-have-access-permission", "test.txt", buf.Bytes())
	assert.Equal(t, "mkdir /should-not-have-access-permission: permission denied", err.Error())
}

func TestListFiles(t *testing.T) {
	var files []string

	root := "./"
	err := filepath.Walk(root, Visit(&files))
	assert.Equal(t, nil, err)

	for _, file := range files {
		log.Debug(file)
	}

	err = filepath.Walk("/dir-does-not-exist", Visit(&files))
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestBasename(t *testing.T) {
	b := Basename(testPath)
	assert.Equal(t, "/a/b/c", b)

	b = Basename(".a/b/c")
	assert.Equal(t, ".a/b/c", b)

	b = Basename(".a")
	assert.Equal(t, ".a", b)
}

func TestFilename(t *testing.T) {
	b := Filename(testPath)
	assert.Equal(t, "c.txt", b)

	b = Filename("test.txt")
	assert.Equal(t, "test.txt", b)

	b = Filename("/test.txt")
	assert.Equal(t, "test.txt", b)
}

func TestBaseDir(t *testing.T) {
	b := BaseDir(testPath)
	assert.Equal(t, "/a/b", b)

	b = BaseDir("/a/b/c")
	assert.Equal(t, "/a/b", b)

	b = BaseDir("/a")
	assert.Equal(t, "/", b)

	b = BaseDir("a/b")
	assert.Equal(t, "a", b)

	b = BaseDir("a")
	assert.Equal(t, "a", b)
}

func TestDirName(t *testing.T) {
	d := DirName("/a/b")
	assert.Equal(t, "b", d)

	d = DirName("/a")
	assert.Equal(t, "a", d)

	d = DirName("a")
	assert.Equal(t, "a", d)
}

func TestEnsureWorkDir(t *testing.T) {
	wd := GetWorkDir()
	d := EnsureWorkDir("")
	assert.Equal(t, wd, d)

	d = EnsureWorkDir(wd + "/..")
	assert.NotEqual(t, wd, d)
}