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

func TestGetByPublisher(t *testing.T){
	testApp := web.NewTestApp().
		SetProperty(web.ViewEnabled, true).
		Run(t)

	testApp.Get("/publisher/{publisher}").
		WithPath("publisher", "hello").
		Expect().Status(200)

}