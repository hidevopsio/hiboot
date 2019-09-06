//go:generate statik -src=./dist

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	_ "hidevops.io/hiboot/examples/web/swagger/http-server/statik"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
	"hidevops.io/hiboot/pkg/utils/io"
	"net/http"
	"path"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/greeting-server"`
}

func init() {
	app.Register(newController)
}

func newController() *controller {
	return &controller{}
}

func (c controller) loadDoc() (retVal []byte, err error) {
	var specDoc *loads.Document
	io.EnsureWorkDir(1, "")
	wd := io.GetWorkDir()
	log.Debug(wd)

	specDoc, err = loads.Spec("./swagger.yml")
	if err != nil {
		return
	}
	retVal, err = json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return
	}
	return
}

func (c *controller) serve(ctx context.Context, docsPath string)  {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	basePath := c.RequestMapping.Value

	handler := middleware.Redoc(middleware.RedocOpts{
		BasePath: basePath,
		SpecURL:  path.Join(basePath, "swagger.json"),
		Path:     docsPath,
	}, http.NotFoundHandler())

	visit := fmt.Sprintf("http://%s%s", ctx.Host(), ctx.Path())

	log.Debugf("visit: %v", visit)

	handler = handlers.CORS()(middleware.Spec(basePath, b, handler))

	ctx.WrapHandler(handler)
}

// UI serve static resource via context StaticResource method
func (c *controller) Swagger(at struct{ at.GetMapping `value:"/swagger.json"`}) ( response string) {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	response = string(b)
	return
}

// UI serve static resource via context StaticResource method
func (c *controller) SwaggerUI(at struct{ at.GetMapping `value:"/swagger-ui"`}, ctx context.Context) {
	c.serve(ctx, at.GetMapping.Value)
	return
}

//run http://localhost:8080/greeting-server/swagger-ui to open swagger ui
func main() {
	web.NewApplication(newController).
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile).
		Run()
}
