package jaeger

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go/log"
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/at"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestConfiguration(t *testing.T) {
	c := newConfiguration()
	assert.NotEqual(t, nil, c)

	tracer := c.Tracer()
	assert.NotEqual(t, reflect.Struct, tracer)

	//ctx:=context.NewContext(fake.ApplicationContext{})
}

// PATH /foo
type Controller struct {
	at.RestController
}

func newController() *Controller {
	return &Controller{
	}
}

// Get GET /foo/{foo}
func (c *Controller) GetByFoo(foo string, span *Span) string {
	defer span.Finish()
	//span.SetTag("hello-to", foo)
	// response

	r := mux.NewRouter()
	r.HandleFunc("/formatter/{name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hi"))
	}))

	server := httptest.NewServer(r)
	defer server.Close()

	req := new(http.Request)
	req, err := http.NewRequest("GET", server.URL+"/formatter/bar", nil)
	if err != nil {

	}
	newSpan := span.Inject(context.Background(), "GET",
		server.URL+"/formatter/bar", req)
	var l Logger
	l.Error("foobar")
	l.Infof("foobar")

	newSpan.LogFields(
		log.String("event", "string-format"),
		log.String("value", "helloStr"),
	)

	return "bar"
}

// Get GET /formatter/{format}
func (c *Controller) GetByFormatter(formatter string, span *ChildSpan) string {
	defer span.Finish()
	greeting := span.BaggageItem("greeting")
	if greeting == "" {
		greeting = "Hello"
	}

	helloStr := fmt.Sprintf("[%s] %s, %s", time.Now().Format(time.Stamp), greeting, formatter)

	url := "http://a.b"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
	}

	_ = span.Inject(context.Background(), "GET", url, req)

	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	// response
	return helloStr
}

func init() {
	app.Register(newController)
}

func TestController(t *testing.T) {

	testApp := web.NewTestApp().
		SetProperty("jaeger.config.serviceName", "test").
		Run(t)

	t.Run("should response 200 when GET /foo/{foo}", func(t *testing.T) {
		testApp.
			Request(http.MethodGet, "/foo/{foo}").
			WithPath("foo", "bar").
			Expect().Status(http.StatusOK)
	})
	t.Run("should response 200 when GET /formatter/{format}", func(t *testing.T) {
		testApp.
			Request(http.MethodGet, "/formatter/{format}").
			WithPath("format", "bar").
			Expect().Status(http.StatusOK)
	})
}

