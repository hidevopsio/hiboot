package main

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/inject/annotation"
	"hidevops.io/hiboot/pkg/log"
	"net/http"
	"testing"
	"time"
)


func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	time.Sleep(time.Second)
	testApp := web.NewTestApp(t, newEmployeeController).Run(t)

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/123").
			Expect().Status(http.StatusOK).
			Body().Contains("123")
	})

	t.Run("should list employee", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})

}

func TestAnnotation(t *testing.T) {
	type foo struct{
		at.GetMapping `value:"/{id:int}"`
		at.ApiOperation `value:"Get an employee"`
		at.ApiParam `value:"Path variable employee ID" required:"true"`
		at.ApiResponse200 `value:"Successfully get an employee"`
		at.ApiResponse404 `value:"The resource you were trying to reach is not found"`
	}
	var f foo
	err := annotation.InjectIntoObject(&f)
	assert.Equal(t, nil, err)
	log.Info(f)
}