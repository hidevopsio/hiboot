package cors

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FooController struct {
	at.RestController

	applicationContext app.ApplicationContext
}

func TestMyFunction(t *testing.T) {
	cfg := newConfiguration()
	assert.NotNil(t, cfg)

	testApp := web.NewTestApp(new(FooController)).
		Run(t)
	cfg.Properties = new(Properties)
	assert.NotNil(t, cfg.Middleware(testApp.(app.ApplicationContext)))
}
