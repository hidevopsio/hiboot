package utils

import (
	"runtime"
	"strings"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"os"
)

func GetWorkingDir(file string) string {
	_, filename, _, _ := runtime.Caller(1)
	wd := strings.Replace(filename, file, "", -1)
	wd2, _ := os.Getwd()
	log.Println("working dir: ", wd, " vs ", wd2)

	return wd;
}