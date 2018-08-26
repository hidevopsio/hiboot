package jwt

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/dgrijalva/jwt-go"
	mw "github.com/iris-contrib/middleware/jwt"
	"github.com/hidevopsio/hiboot/pkg/factory/instance"
)

type configuration struct{
	app.Configuration

	Properties Properties `mapstructure:"jwt"`

	jwtHandler *JwtMiddleware
	instanceFactory instance.Factory

}

func init() {
	app.AutoConfiguration(new(configuration))
}

func (c *configuration ) Init(instanceFactory instance.Factory) {
	c.instanceFactory = instanceFactory
}

func (c *configuration ) jwtMiddleware(jt *jwtToken) *JwtMiddleware {
	return NewJwtMiddleware(mw.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//log.Debug(token)
			return jt.verifyKey, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func (c *configuration) JwtToken() Token  {
	jt := new(jwtToken)
	//wd := io.GetWorkDir()

	jt.Initialize(&c.Properties)

	jwtMiddleware := c.jwtMiddleware(jt)
	c.instanceFactory.SetInstance("jwtMiddleware", jwtMiddleware)
	jt.jwtEnabled = true

	return jt
}
