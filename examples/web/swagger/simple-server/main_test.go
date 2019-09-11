package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app/web"
	"net/http"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()
}

func TestController(t *testing.T) {
	time.Sleep(time.Second)
	testApp := web.NewTestApp(t).Run(t)

	t.Run("should get employee ", func(t *testing.T) {
		testApp.Get("/employee/123").
			Expect().Status(http.StatusOK)
	})

	t.Run("should report 404 when employee does not exist", func(t *testing.T) {
		testApp.Get("/employee/100").
			Expect().Status(http.StatusNotFound)
	})

	t.Run("should list employee", func(t *testing.T) {
		testApp.Get("/employee").
			Expect().Status(http.StatusOK)
	})

}

func TestCloneRef(t *testing.T) {
	var b bytes.Buffer
	src := spec.MustCreateRef("#/definitions/test")
	err := gob.NewEncoder(&b).Encode(&src)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	var dst spec.Ref
	err = gob.NewDecoder(&b).Decode(&dst)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	jazon, err := json.Marshal(dst)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, `{"$ref":"#/definitions/test"}`, string(jazon))
}
