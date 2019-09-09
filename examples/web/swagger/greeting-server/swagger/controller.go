package swagger

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/spec"
	"github.com/gorilla/handlers"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/system"
	"hidevops.io/hiboot/pkg/utils/mapstruct"
	"net/http"
	"path"
	"path/filepath"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/"`

	openAPIDefinition *OpenAPIDefinition
	builder system.Builder

	// should inject api builder
}

func init() {
	app.Register(newController)
}

func newController(builder system.Builder) *controller {
	c := &controller{builder: builder}

	c.openAPIDefinition = OpenAPIDefinitionBuilder().(*OpenAPIDefinition)
	_ = c.builder.Load(c.openAPIDefinition, mapstruct.WithSquash)

	return c
}

// TODO: add description 'Implemented by HiBoot Framework'
func (c *controller) loadDoc() (retVal []byte, err error) {
	// TODO: move to api builder
	paths := c.builder.GetProperty("swagger.paths").(map[string]interface{})
	for k, p := range paths {
		pi := spec.PathItem{}
		err = mapstruct.Decode(&pi, p, mapstruct.WithSquash)
		c.openAPIDefinition.Paths.Paths[k] = pi
	}
	retVal, err = json.MarshalIndent(c.openAPIDefinition.Swagger, "", "  ")
	return
}

func (c *controller) serve(ctx context.Context, docsPath string) {
	b, err := c.loadDoc()
	if err != nil {
		return
	}
	basePath := filepath.Join(c.openAPIDefinition.Swagger.BasePath, c.RequestMapping.Value)

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

