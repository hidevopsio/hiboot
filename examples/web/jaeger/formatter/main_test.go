package main

import (
	"hidevops.io/hiboot/pkg/app/web"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()
	time.Sleep(time.Second)
}

func TestGetByFormatter(t *testing.T){
	testApp := web.NewTestApp().
		SetProperty(web.ViewEnabled, true).
		Run(t)

	testApp.Get("/formatter/{format}").
		WithPath("format", "hello").
		Expect().Status(200)

}