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

// Package base64 provides base64 encryption/decryption utilities
package base64

import "encoding/base64"

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

// EncodeToString encode src string to base64 format string
func EncodeToString(src string) string {
	return coder.EncodeToString([]byte(src))
}

// Encode encode src bytes to base64 format bytes
func Encode(src []byte) (dst []byte) {
	dst = make([]byte, coder.EncodedLen(len(src)))
	coder.Encode(dst, src)
	return
}

// DecodeToString decode string from base64 string
func DecodeToString(src string) (retVal string, err error) {
	retBytes, err := coder.DecodeString(src)
	retVal = string(retBytes)
	return
}

// Decode decode base64 bytes
func Decode(src []byte) (dst []byte, err error) {
	size := coder.DecodedLen(len(src))
	dst = make([]byte, size)
	_, err = coder.Decode(dst, src)
	return
}
