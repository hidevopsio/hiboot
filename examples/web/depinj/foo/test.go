package foo

import "github.com/hidevopsio/hiboot/pkg/app"

type TestService struct {
	Name string `value:"${foo.bar.foz:.foo}"`
}

func newTestService() *TestService {
	return &TestService{}
}

func init() {
	app.Register(newTestService)
}
