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

// Package md5 provides md5 encryption utilities
package md5

import (
	"crypto/md5"
	"encoding/hex"
)

// Encrypt md5 encryption util
func Encrypt(in string) (out string) {
	md5Ctx := md5.New()                // md5 init
	n, err := md5Ctx.Write([]byte(in)) // md5 update
	if err == nil && n != 0 {
		cipherStr := md5Ctx.Sum(nil)        // md5 final
		out = hex.EncodeToString(cipherStr) // hex digest
	}
	return
}
