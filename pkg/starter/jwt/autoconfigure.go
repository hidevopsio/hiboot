package jwt

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/dgrijalva/jwt-go"
	mw "github.com/iris-contrib/middleware/jwt"
)

type configuration struct{
	app.Configuration

	Properties Properties `mapstructure:"jwt"`

	jwtHandler *JwtMiddleware

}

func init() {
	app.AutoConfiguration(new(configuration))
}

func (c *configuration ) JwtMiddleware() *JwtMiddleware {

	jt := c.JwtToken()

	return NewJwtMiddleware(mw.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//log.Debug(token)
			return jt.VerifyKey(), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func (c *configuration) JwtToken() Token  {
	jt := new(jwtToken)

	jt.Initialize(&c.Properties)

	return jt
}
