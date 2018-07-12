package gotest

import (
	"os"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"strings"
)

func IsRunning() bool  {

	args := os.Args

	//log.Println("args: ", args)
	//log.Println("args[0]: ", args[0])

	if utils.StringInSlice("-test.v", args) ||
		strings.Contains(args[0], ".test") {
		return true
	}

	return false
}