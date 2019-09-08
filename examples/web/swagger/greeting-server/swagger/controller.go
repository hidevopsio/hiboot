package swagger

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"net/http"
	"path"
	"path/filepath"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/"`
	at.DisableSwagger
	openAPIDefinition *OpenAPIDefinition
}

func init() {
	app.Register(newController)
}

func newController(openAPIDefinition *OpenAPIDefinition) *controller {
	return &controller{openAPIDefinition: openAPIDefinition}
}

// TODO: add description 'Implemented by HiBoot Framework'
func (c controller) loadDoc() (retVal []byte, err error) {
	if c.openAPIDefinition != nil {
		swgSpec := &c.openAPIDefinition.Swagger

		retVal, err = json.MarshalIndent(swgSpec, "", "  ")
	}

	return
}

func (c *controller) serve(ctx context.Context, docsPath string) {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	basePath := filepath.Join(c.openAPIDefinition.BasePath, c.RequestMapping.Value)

	handler := middleware.Redoc(middleware.RedocOpts{
		BasePath: basePath,
		SpecURL:  path.Join(basePath, "swagger.json"),
		Path:     docsPath,
	}, http.NotFoundHandler())

	//visit := fmt.Sprintf("http://%s%s", ctx.Host(), ctx.Path())
	//log.Debugf("visit: %v", visit)

	handler = handlers.CORS()(middleware.Spec(basePath, b, handler))

	ctx.WrapHandler(handler)
}

// UI serve static resource via context StaticResource method
func (c *controller) Swagger(at struct{ at.GetMapping `value:"/swagger.json"` }) (response string) {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	response = string(b)
	return
}

// UI serve static resource via context StaticResource method
func (c *controller) SwaggerUI(at struct{ at.GetMapping `value:"/swagger-ui"` }, ctx context.Context) {
	c.serve(ctx, at.GetMapping.Value)
	return
}

