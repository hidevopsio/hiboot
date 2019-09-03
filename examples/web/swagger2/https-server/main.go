//go:generate statik -src=./public

package main

import (
	"github.com/rakyll/statik/fs"
	_ "hidevops.io/hiboot/examples/web/swagger2/https-server/statik"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"net/http"
)

// Before run go build, run go generate.
// Then, run the main program and visit http://localhost:8080/public/hello.txt
func main() {
	// static files path prefix
	stripPrefix := "/public/"

	staticFiles, err := fs.New()
	if err == nil {
		http.Handle(stripPrefix, http.StripPrefix(stripPrefix, http.FileServer(staticFiles)))

		// Run HiBoot Application
		web.NewApplication().SetProperty(app.ProfilesInclude, actuator.Profile).Run()
	}
}
