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

// Package io provides file or directory io utilities.
package io

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ChangeWorkDir change current working dir
func ChangeWorkDir(workDir string) error {
	return os.Chdir(workDir)
}

// GetWorkDir get current working dir
func GetWorkDir() string {

	wd, _ := os.Getwd()

	return wd
}

// EnsureWorkDir ensure the working dir is set the correct dir specified by user
func EnsureWorkDir(skip int, dir string) (ok bool) {
	var path string
	if _, file, _, ok := runtime.Caller(skip); ok && strings.Contains(os.Args[0], "go_build_") {
		path = BaseDir(file)
	} else {
		path = GetWorkDir()
	}
	lastPath := ""

	for {
		//log.Debugf("%v", path)
		configPath := filepath.Join(path, dir)
		if !IsPathNotExist(configPath) {
			ChangeWorkDir(path)
			ok = true
			break
		}

		path = BaseDir(path)
		if lastPath == path {
			break
		}
		lastPath = path
	}
	return
}

// GetRelativePath get relative path
func GetRelativePath(level int) string {
	_, path, _, _ := runtime.Caller(level)

	return filepath.Dir(path)
}

// IsPathNotExist check if path is not exist
func IsPathNotExist(path string) bool {
	_, err := os.Stat(path)
	isNotExist := os.IsNotExist(err)
	return isNotExist
}

func write(path, filename string, cb func(f *os.File) (n int, err error)) (int, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
	}
	if err != nil {
		return 0, err
	}

	f, _ := os.OpenFile(filepath.Join(path, filename), os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	if cb != nil {
		return cb(f)
	}
	return 0, err
}

// CreateFile create file
func CreateFile(path, filename string) error {
	_, err := write(path, filename, nil)
	return err
}

// WriterFile write bytes to file
func WriterFile(path, filename string, in []byte) (int, error) {
	return write(path, filename, func(f *os.File) (int, error) {
		return f.Write(in)
	})
}

// Visit get files in a specific dir
func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		*files = append(*files, path)
		return nil
	}
}

// Basename get the base name from a path
// /a/b/c.ext => /a/b/c
func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

// Filename get file name from a path
// /a/b/c.ext => c.ext
func Filename(s string) string {
	n := strings.LastIndexByte(s, filepath.Separator)
	if n >= 0 {
		return s[n+1:]
	}
	return s
}

// BaseDir get base dir from a path
// /a/b/c.ext => /a/b
// /a/b/c => /a/b
func BaseDir(s string) string {
	n := strings.LastIndexByte(s, filepath.Separator)
	if n > 0 {
		return s[:n]
	} else if n == 0 {
		return s[:n+1]
	}
	return s
}

// DirName get dir name from a path
// /a/b/c => c
func DirName(s string) string {
	n := strings.LastIndexByte(s, filepath.Separator)
	if n >= 0 {
		return s[n+1:]
	}
	return s
}

// CallerInfo get call info, include filename, line number or function name
func CallerInfo(skip int) (file string, line int, fn string) {
	var pc uintptr
	var ok bool
	if pc, file, line, ok = runtime.Caller(skip); ok {
		fn = runtime.FuncForPC(pc).Name()
	}
	return
}
