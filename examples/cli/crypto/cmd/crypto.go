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

package cmd

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto/rsa"
	"github.com/hidevopsio/hiboot/pkg/app"
)

// define the command
type rootCommand struct {
	// embedding cli.BaseCommand in each command
	cli.BaseCommand
	// inject (bind) flag to field 'Source', 'Encrypt', and 'Decrypt', so that it can be used on Run method, please note that the data type must be pointer
	Source  string
	Encrypt bool
	Decrypt bool
	Key     string
}

func init() {
	app.Component(newRootCommand)
}

func newRootCommand() *rootCommand {
	c := new(rootCommand)
	c.Use = "crypto"
	c.Short = "crypto command"
	c.Long = "run crypto command to encrypt/decrypt "
	c.Example = `
crypto rsa -h
crypto rsa -e -s "text to encrypt"
crypto rsa -d -s "text to decrypt"
`
	pflags := c.PersistentFlags()
	pflags.StringVarP(&c.Source, "source", "s", "", "run with option --source=source text to encrypt or encrypt")
	pflags.StringVarP(&c.Key, "key", "k", "", "run with option --key or -k for rsa key")
	pflags.BoolVarP(&c.Encrypt, "encrypt", "e", false, "run with option --encrypt or -e for text encryption")
	pflags.BoolVarP(&c.Decrypt, "decrypt", "d", false, "run with option --decrypt or -d for text encryption")
	return c
}

// Run OnRsa for crypto command rsa
func (c *rootCommand) OnRsa(args []string) bool {
	if c.Decrypt {
		res, err := rsa.DecryptBase64([]byte(c.Source), []byte(c.Key))
		if err == nil {
			fmt.Println(string(res))
		}
	} else {
		res, err := rsa.EncryptBase64([]byte(c.Source), []byte(c.Key))
		if err == nil {
			fmt.Println(string(res))
		}
	}
	return true
}
