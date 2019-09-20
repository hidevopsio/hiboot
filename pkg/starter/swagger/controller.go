// Copyright 2018 ~ now John Deng (hi.devops.io@gmail.com).
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

package swagger

import (
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"hidevops.io/hiboot/pkg/app/web/context"
	"hidevops.io/hiboot/pkg/at"
	"net/http"
	"path"
	"path/filepath"
)

type controller struct {
	at.RestController
	at.RequestMapping `value:"/"`

	apiInfoBuilder *apiInfoBuilder
}

func newController(openAPIDefinition *apiInfoBuilder) *controller {
	return &controller{apiInfoBuilder: openAPIDefinition}
}

// TODO: add description 'Implemented by HiBoot Framework'
func (c *controller) loadDoc() (retVal []byte, err error) {
	retVal, err = json.MarshalIndent(c.apiInfoBuilder.Swagger, "", "  ")
	return
}

func (c *controller) serve(ctx context.Context, docsPath string) {
	b, err := c.loadDoc()
	if err == nil {
		basePath := filepath.Join(c.apiInfoBuilder.Swagger.BasePath, c.RequestMapping.Value)

		handler := middleware.Redoc(middleware.RedocOpts{
			BasePath: basePath,
			SpecURL:  path.Join(basePath, "swagger.json"),
			Path:     docsPath,
		}, http.NotFoundHandler())

		handler = handlers.CORS()(middleware.Spec(basePath, b, handler))

		ctx.WrapHandler(handler)
	}
}

// UI serve static resource via context StaticResource method
func (c *controller) Swagger(at struct{ at.GetMapping `value:"/swagger.json"` }) (response string) {
	b, err := c.loadDoc()
	if err == nil {
		response = string(b)
	}
	return
}

// UI serve static resource via context StaticResource method
func (c *controller) SwaggerUI(at struct{ at.GetMapping `value:"/swagger-ui"` }, ctx context.Context) {
	c.serve(ctx, at.GetMapping.Value)
	return
}

