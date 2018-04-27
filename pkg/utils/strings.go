package utils

import (
	"unicode"
	"strings"
)

func UpperFirst(str string) string {
	return strings.Title(str)
}

func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
