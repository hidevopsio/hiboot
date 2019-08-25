package main

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

type swagger struct {
	at.Swagger `value:"2.0"`
	Info struct{
		Title string `value:"The demo application of a tutorial on https://hiboot.hidevopsio.io" json:"title"`
		Description string `value:"A demo application" json:"description"`
		Version string `value:"1.0.0" json:"version"`
	} `json:"info"`
	Schemes []string `value:"http" json:"schemes"`
}

func init() {
	app.Register(new(swagger))
}