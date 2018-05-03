package web

import (
	"github.com/kataras/iris"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"os"
	"io/ioutil"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/kataras/iris/context"
	"reflect"
	"github.com/fatih/camelcase"
	"strings"
	"github.com/kataras/iris/core/router"
	"github.com/kataras/iris/httptest"
	"testing"
	"github.com/iris-contrib/httpexpect"
)

const (
	privateKeyPath = "/config/app.rsa"
	pubKeyPath     = "/config/app.rsa.pub"

	pathSep        = "/"
	AuthAnon       = "anon"
)

type ApplicationInterface interface {
	Init()
	Config() *system.Configuration
	GetSignKey() *rsa.PrivateKey
	Run()
	NewTestServer(t *testing.T) *httpexpect.Expect
}

type Application struct {
	app    *iris.Application
	config *system.Configuration
}

type Health struct {
	Status string `json:"status"`
}

type Controller struct{
	Name string
}

var (
	jwtHandler *jwtmiddleware.Middleware
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
)

const (
	application = "application"
	config      = "/config"
	yaml        = "yaml"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func (wa *Application) Init() {
	wd := utils.GetWorkingDir("")

	builder := &system.Builder{
		Path:       wd + config,
		Name:       application,
		FileType:   yaml,
		Profile:    os.Getenv("APP_PROFILES_ACTIVE"),
		ConfigType: system.Configuration{},
	}
	cp, err := builder.Build()
	if err == nil {
		wa.config = cp.(*system.Configuration)
		log.SetLevel(wa.config.Logging.Level)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	// check if key exist
	if utils.IsPathNotExist(wd + privateKeyPath) {
		return
	}

	signBytes, err := ioutil.ReadFile(wd + privateKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(wd + pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	jwtHandler = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//log.Debug(token)
			return verifyKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodRS256,
	})

	log.Debug("application init")
}

func (wa *Application) Config() *system.Configuration {
	return wa.config
}

func (wa *Application) GetSignKey() *rsa.PrivateKey {
	return signKey
}

func (wa *Application) Run() {
	serverPort := ":8080"
	if wa.config != nil {
		serverPort = fmt.Sprintf(":%v", wa.config.Server.Port)
	}
	// TODO: WithCharset should be configurable
	wa.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}

func (wa *Application) NewTestServer(t *testing.T) *httpexpect.Expect {
	return httptest.New(t, wa.app)
}

func healthHandler(app *iris.Application) *router.Route {
	return app.Get("/health", func(ctx context.Context) {
		health := Health{
			Status: "UP",
		}
		ctx.JSON(health)
	})
}

func (wa *Application) handle(method reflect.Method, object interface{}, ctx context.Context) {
	//log.Debug("NumIn: ", method.Type.NumIn())
	inputs := make([]reflect.Value, method.Type.NumIn())

	inputs[0] = reflect.ValueOf(object)
	inputs[1] = reflect.ValueOf(ctx)
	method.Func.Call(inputs)
}

func NewApplication(controllers interface{}) (*Application, error) {

	wa := &Application{}

	wa.Init()

	app := iris.New()


	// The only one Required:
	// here is how you define how your own context will
	// be created and acquired from the iris' generic context pool.
	app.ContextPool.Attach(func() context.Context {
		return &Context{
			// Optional Part 3:
			Context: context.NewContext(app),
		}
	})

	wa.app = app

	healthHandler(app)

	err := utils.ValidateReflectType(controllers, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
		appliedJwt := false
		for i := 0; i < fieldSize; i++ {
			for _, field := range utils.DeepFields(reflectType) {
				fieldName := field.Name
				fieldType := field.Type
				controller := value.FieldByName(fieldName).Interface()
				if fieldType.Elem().Kind() == reflect.Struct {
					log.Debug("name: ", fieldName)
					log.Debug("tag: ", field.Tag)

					controllerName := field.Tag.Get("controller")
					if controllerName == "" {
						controllerName = strings.ToLower(fieldName)
					}
					auth := field.Tag.Get("auth")
					log.Debug("controller: ", controllerName)
					log.Debug("auth: ", auth)
					if ! appliedJwt && auth != "" && ! (auth == AuthAnon) {
						appliedJwt = true
						app.Use(jwtHandler.Serve)
					}
					contextPath := pathSep + controllerName

					beforeMethod, ok := fieldType.MethodByName("Before")
					var party iris.Party
					if ok {
						log.Print(contextPath)
						log.Print(beforeMethod.Name)
						party = app.Party(contextPath, func(ctx context.Context) {
							wa.handle(beforeMethod, controller, ctx)
						})
					}

					numOfMethod := fieldType.NumMethod()
					log.Debug("methods: ", numOfMethod)
					for mi := 0; mi < numOfMethod; mi++ {
						method := fieldType.Method(mi)
						log.Debug("method: ", method.Name)
						methodName := method.Name

						// TODO: contextMapping naming can be configured in applicatino.yml, e.g camelCase or ...

						ctxMap := camelcase.Split(methodName)
						httpMethod := strings.ToUpper(ctxMap[0])
						contextMapping := strings.Replace(methodName, ctxMap[0], "", 1)
						contextMapping = pathSep + utils.LowerFirst(contextMapping)

						if party == nil {
							relativePath := contextPath + contextMapping
							log.Debug("relativePath: ", relativePath)
							app.Handle(httpMethod, relativePath, func(ctx context.Context) {
								wa.handle(method, controller, ctx)
							})
						} else {
							log.Debug("contextMapping: ", contextMapping)
							party.Handle(httpMethod, contextMapping, func(ctx context.Context) {
								wa.handle(method, controller, ctx)
							})
						}
					}
				}
			}
		}
		return nil
	})

	return wa, err
}

