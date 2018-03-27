package utils

import (
	"runtime"
	"strings"
	"github.com/hidevopsio/hi/boot/pkg/log"
)

func GetWorkingDir(file string) string {
	_, filename, _, _ := runtime.Caller(1)
	wd := strings.Replace(filename, file, "", -1)
	log.Debugf("working dir: %s", wd)

	return wd;
}