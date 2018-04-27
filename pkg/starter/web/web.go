package web

import (
	"github.com/kataras/iris"
	"crypto/rsa"
	"time"
	"github.com/dgrijalva/jwt-go"
	"os"
	"io/ioutil"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"fmt"
	"github.com/hidevopsio/hiboot/pkg/system"
	"github.com/hidevopsio/hiboot/pkg/model"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/kataras/iris/context"
	"reflect"
	"net/http"
)

var (
	Application iris.Application
	Context iris.Context
)


const (
	privateKeyPath = "/config/app.rsa"
	pubKeyPath     = "/config/app.rsa.pub"

	pathSep        = "/"
	AuthAnon       = "anon"
)

type Boot struct {
	app    *iris.Application
	config *system.Configuration
}

type MapJwt map[string]interface{}

type JwtToken string

type Health struct {
	Status string `json:"status"`
}

type Controller struct{
	Name string
}

type controllerMethod func(*Controller, context.Context)

var (
	jwtHandler *jwtmiddleware.Middleware
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
	sysCfg     *system.Configuration
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

func init() {

}

func GenerateJwtToken(payload MapJwt, expired int64, unit time.Duration) (JwtToken, error) {

	claim := jwt.MapClaims{
		"exp": time.Now().Add(unit * time.Duration(expired)).Unix(),
		"iat": time.Now().Unix(),
	}

	for k, v := range payload {
		claim[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signKey)

	jwtToken := JwtToken(tokenString)

	return jwtToken, err
}

// set response
func Response(ctx iris.Context, message string, data interface{}) {
	response := &model.Response{
		Code:    ctx.GetStatusCode(),
		Message: message,
		Data:    data,
	}

	// just for debug now
	ctx.JSON(response)
}

// set response
func ResponseError(ctx iris.Context, message string, code int) {
	response := &model.Response{
		Code:    code,
		Message: message,
	}

	// just for debug now
	ctx.StatusCode(code)
	ctx.JSON(response)
}

func (b *Boot) Init() {
	wd := utils.GetWorkingDir("")

	builder := &system.Builder{
		Path:       wd + config,
		Name:       application,
		FileType:   yaml,
		Profile:    os.Getenv("APP_PROFILES_ACTIVE"),
		ConfigType: system.Configuration{},
	}
	cp, err := builder.Build()
	b.config = cp.(*system.Configuration)
	log.SetLevel(b.config.Logging.Level)

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

	b.app = iris.New()

	b.app.Get("/health", func(ctx context.Context) {
		health := Health{
			Status: "UP",
		}
		ctx.JSON(health)
	})
}

func (b *Boot) App() *iris.Application {
	return b.app
}

func (b *Boot) Config() *system.Configuration {
	return b.config
}

func (b *Boot) GetSignKey() *rsa.PrivateKey {
	return signKey
}

func (b *Boot) ApplyJwt() {
	b.app.Use(jwtHandler.Serve)
}

func (b *Boot) Run() {
	serverPort := fmt.Sprintf(":%v", b.config.Server.Port)
	b.app.Run(iris.Addr(fmt.Sprintf(serverPort)), iris.WithCharset("UTF-8"), iris.WithoutVersionChecker)
}

func NewApplication(controllers interface{}) (*Boot, error) {
	// iris app
	boot := new(Boot)

	boot.Init()

	app := boot.App()

	log.Print(app)

	err := utils.ValidateReflectType(controllers, func(value *reflect.Value, reflectType reflect.Type, fieldSize int, isSlice bool) error {
		appliedJwt := false
		for i := 0; i < fieldSize; i++ {
			for _, field := range utils.DeepFields(reflectType) {
				fieldName := field.Name
				fieldType := field.Type
				if fieldType.Elem().Kind() == reflect.Struct{
					log.Print("name: ", fieldName)
					log.Print("tag: ", field.Tag)

					controller := field.Tag.Get("controller")
					auth := field.Tag.Get("auth")
					log.Print("controller: ", controller)
					log.Print("auth: ", auth)
					if ! appliedJwt && ! (auth == AuthAnon) {
						appliedJwt = true
						boot.ApplyJwt()
					}
					contextPath := pathSep + controller

					beforeMethod, ok := fieldType.MethodByName("Before")
					var party iris.Party
					if ok {
						log.Print(contextPath)
						log.Print(beforeMethod.Name)
						//party = app.Party(contextPath, beforeMethod.Func.Interface().(context.Handler))
					}

					numOfMethod := fieldType.NumMethod()
					log.Print("methods: ", numOfMethod)
					for mi := 0; mi < numOfMethod; mi++{
						method := fieldType.Method(mi)
						log.Print("method: ", method.Name)
						methodName := method.Name

						// TODO: contextMapping naming can be configured in applicatino.yml, e.g camelCase or ...

						contextMapping := pathSep + utils.LowerFirst(methodName)
						methodInterface := method.Func.Interface()
						log.Print(methodInterface)

						if party == nil {
							log.Print(contextPath + contextMapping)
							app.Handle(http.MethodPost, contextPath + contextMapping, methodInterface.(context.Handler))
						} else {
							log.Print(contextMapping)
							//party.Post(contextMapping, methodInterface.(context.Handler))
						}
					}
				}
			}
		}
		return nil
	})

	return boot, err
}