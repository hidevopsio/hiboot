package jwt

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/utils/io"
	"github.com/hidevopsio/hiboot/pkg/log"
)

func init() {
	log.SetLevel(log.DebugLevel)
	wd := io.EnsureWorkDir("../../..")
	log.Debugf("wd: %v", wd)
}

func TestAutoConfigure(t *testing.T) {

	config := &configuration{
		Properties: Properties{
			PrivateKeyPath: "config/ssl/app.rsa",
			PublicKeyPath: "config/ssl/app.rsa.pub",
		},
	}

	token := config.JwtToken()
	assert.NotEqual(t, nil, token)
	mw := config.JwtMiddleware()
	assert.NotEqual(t, nil, mw)
}
