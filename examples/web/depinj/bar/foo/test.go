package foo

import "github.com/hidevopsio/hiboot/pkg/app"

type TestService struct {
	Name string `value:"${foo.bar.baz:.bar}"`
}

func newTestService() *TestService {
	return &TestService{}
}

func init() {
	app.Register(newTestService)
}
