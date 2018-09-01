package cmd

import (
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/app/cli"
	"github.com/hidevopsio/hiboot/pkg/utils/crypto/rsa"
)

// define the command
type CryptoCommand struct {
	// embedding cli.BaseCommand in each command
	cli.BaseCommand
	// inject (bind) flag to field 'Source', 'Encrypt', and 'Decrypt', so that it can be used on Run method, please note that the data type must be pointer
	Source  *string `flag:"shorthand=s,usage=run with option --source=source text to encrypt or encrypt"`
	Encrypt *bool   `flag:"shorthand=e,usage=run with option --encrypt or -e for text encryption"`
	Decrypt *bool   `flag:"shorthand=d,usage=run with option --decrypt or -d for text decryption"`
	Key     *string `flag:"shorthand=k,usage=run with option --key or -k for rsa key"`
}

// Init constructor
func (c *CryptoCommand) Init() {
	c.Use = "crypto"
	c.Short = "crypto command"
	c.Long = "run crypto command to encrypt/decrypt "
	c.Example = `
crypto rsa -h
crypto rsa -e -s "text to encrypt"
crypto rsa -d -s "text to decrypt"
`
}

// Run OnRsa for crypto command rsa
func (c *CryptoCommand) OnRsa(args []string) bool {
	if *c.Decrypt {
		res, err := rsa.DecryptBase64([]byte(*c.Source), []byte(*c.Key))
		if err == nil {
			fmt.Println(string(res))
		}
	} else {
		res, err := rsa.EncryptBase64([]byte(*c.Source), []byte(*c.Key))
		if err == nil {
			fmt.Println(string(res))
		}
	}
	return true
}
