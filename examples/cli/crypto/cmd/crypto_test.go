package cmd

import (
	"testing"
	"github.com/hidevopsio/hiboot/pkg/starter/cli"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
)


func init() {
	log.SetLevel(log.DebugLevel)
}


func TestCryptoCommands(t *testing.T) {
	cryptoCmd := new(CryptoCommand)
	testApp := cli.NewTestApplication(cryptoCmd)

	t.Run("should run crypto rsa -e", func(t *testing.T) {
		_, err := testApp.RunTest("rsa", "-e", "-s", "hello")
		assert.Equal(t, nil, err)
	})

	t.Run("should run crypto rsa -d ", func(t *testing.T) {
		_, err := testApp.RunTest("rsa", "-d", "-s", "Rprrfl5LX9NRmWKEqJW8ckObVjznnMmq8i7x6Pv6n1GSoEL9dUomNKOr6Pgj7RuVzCc/I7Hya20BZO1PbzTquBMp/G5rcF2Vy7HF1UKr8buHtppB+n3ycTxFvPxQB2vMvLyMtDBc29QtGe3HHD8TS+3h1pSK5WZS+CMKPHT4sho=")
		assert.Equal(t, nil, err)
	})
}


