package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/httpclient/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRunMain(t *testing.T) {
	go main()
	time.Sleep(time.Second)
}

type mockReadCloser struct {
	mock.Mock
}

func (m mockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestController(t *testing.T) {

	t.Run("should get http status 200", func(t *testing.T) {
		r := mux.NewRouter()
		r.HandleFunc("/formatter/{name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("hi"))
		}))
		r.HandleFunc("/publisher/{name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("hi"))
		}))

		server := httptest.NewServer(r)
		defer server.Close()

		testApp := web.NewTestApp().
			SetProperty(web.ViewEnabled, true).
			SetProperty("provider.formatter", server.URL).
			SetProperty("provider.publisher", server.URL).
			Run(t)

		testApp.Get("/greeting/{greeting}/name/{name}/").
			WithPath("greeting", "hello").
			WithPath("name", "world").
			Expect().Status(http.StatusOK)

		server = httptest.NewServer(r)
		testApp.Get("/greeting/{greeting}/name/{name}/").
			WithPath("greeting", "hello").
			WithPath("name", "world").
			Expect().Status(http.StatusOK)
	})
	t.Run("should get http status 200 when client.Get return error", func(t *testing.T) {
		testApp := web.NewTestApp().
			SetProperty(web.ViewEnabled, true).
			Run(t)

		testApp.Get("/greeting/{greeting}/name/{name}/").
			WithPath("greeting", "hello").
			WithPath("name", "world").
			Expect().Status(200)

	})
	t.Run("should get http status 200 when get an error reading the body from formatter in formatString", func(t *testing.T) {

		mockReadCloser := mockReadCloser{}
		// if Read is called, it will return error
		mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
		// if Close is called, it will return error
		mockReadCloser.On("Close").Return(fmt.Errorf("error closing"))

		stringReadCloser := ioutil.NopCloser(mockReadCloser)
		resp := &http.Response{Status: "200", Body: stringReadCloser}

		mockClient := new(mocks.Client)
		mockClient.On("Get",
			"http://localhost:8081/formatter/world",
			http.Header(nil), mock.AnythingOfType("func(*http.Request)")).
			Return(resp, nil)

		controller := newController(mockClient)

		testApp := web.NewTestApp(controller).
			SetProperty(web.ViewEnabled, true).
			Run(t)

		testApp.Get("/greeting/{greeting}/name/{name}/").
			WithPath("greeting", "hello").
			WithPath("name", "world").
			Expect().Status(200)
	})
}
