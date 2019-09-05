//go:generate statik -src=./public

package main

import (
	"github.com/rakyll/statik/fs"
	_ "hidevops.io/hiboot/examples/web/swagger/http-server/statik"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"net/http"
)

// Before run go build, run go generate.
// Then, run the main program and visit http://localhost:8080/public/hello.txt or http://localhost:8080/public/img/hiboot.png
func main() {
	// static files path prefix
	stripPrefix := "/public/"

	// create new static resources
	staticFiles, err := fs.New()
	if err == nil {
		http.Handle(stripPrefix, http.StripPrefix(stripPrefix, http.FileServer(staticFiles)))
	}
	// Run HiBoot Application
	web.NewApplication().SetProperty(app.ProfilesInclude, actuator.Profile).Run()
}
