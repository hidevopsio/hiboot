package main

import (
	"os"
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	_ "github.com/hidevopsio/hiboot/examples/web/jwt/controllers"
)

func main()  {

	app, err := web.NewApplication()
	if err != nil {

		log.Error(err)
		os.Exit(1)
	}

	app.Run()
}