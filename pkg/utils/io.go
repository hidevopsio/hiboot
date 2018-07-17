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
	"runtime"
	"os"
	"path/filepath"
	"strings"
)

func ChangeWorkDir(workdir string) {

	os.Chdir(workdir)

}

func GetWorkDir() string {

	wd, _ := os.Getwd()

	return wd
}


func GetRelativePath(level int) string {
	_, path, _, _ := runtime.Caller(level)

	return filepath.Dir(path)
}

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

	f, _ := os.OpenFile(filepath.Join(path, filename), os.O_RDWR | os.O_CREATE, 0666)
	defer f.Close()
	if cb != nil {
		return cb(f)
	}
	return 0, err
}

func CreateFile(path, filename string) (error) {
	_, err := write(path, filename, nil)
	return err
}

func WriterFile(path, filename string, in []byte) (int, error) {
	return write(path, filename, func(f *os.File) (int, error) {
		return f.Write(in)
	})
}

func Visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		*files = append(*files, path)
		return nil
	}
}

func Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

func Filename(s string) string {
	n := strings.LastIndexByte(s, filepath.Separator)
	if n >= 0 {
		return s[n + 1:]
	}
	return s
}

func BaseDir(s string) string {
	n := strings.LastIndexByte(s, filepath.Separator)
	if n > 0 {
		return s[:n]
	} else if n == 0 {
		return s[:n + 1]
	}
	return s
}


func DirName(s string) string {
	n := strings.LastIndexByte(s, filepath.Separator)
	if n >= 0 {
		return s[n + 1:]
	}
	return s
}


func EnsureWorkDir(path string) string  {
	wd := GetWorkDir()
	if ! strings.Contains(wd, path) {
		ChangeWorkDir(path)
		wd = GetWorkDir()
	}
	return wd
}