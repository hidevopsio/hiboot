package web

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/hidevopsio/hiboot/pkg/log"
	mwjwt "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// JwtMiddleware derrived from github.com/iris-contrib/middleware/jwt/Middleware
type JwtMiddleware struct {
	mwjwt.Middleware
}

// Serve the middleware's action
func (m *JwtMiddleware) Serve(ctx context.Context) {
	if err := m.CheckJWT(ctx); err != nil {
		c := ctx.(*Context)
		c.ResponseError(err.Error(), http.StatusUnauthorized)
		c.StopExecution()
		return
	}
	// If everything ok then call next.
	ctx.Next()
}

// CheckJWT the main functionality, checks for token
func (m *JwtMiddleware) CheckJWT(ctx context.Context) error {
	if !m.Config.EnableAuthOnOptions {
		if ctx.Method() == iris.MethodOptions {
			return nil
		}
	}

	// Use the specified token extractor to extract a token from the request
	token, err := m.Config.Extractor(ctx)

	// If debugging is turned on, log the outcome
	if err != nil {
		log.Errorf("Error extracting JWT: %v", err)
	} else {
		log.Debugf("Token extracted: %s", token)
	}

	// If an error occurs, call the error handler and return an error
	if err != nil {
		return fmt.Errorf("Error extracting token: %v", err)
	}

	// If the token is empty...
	if token == "" {
		// Check if it was required
		if m.Config.CredentialsOptional {
			log.Debug("  No credentials found (CredentialsOptional=true)")
			// No error, just no token (and that is ok given that CredentialsOptional is true)
			return nil
		}

		// If we get here, the required token is missing
		errorMsg := "Required authorization token not found"
		log.Debug("  Error: No credentials found (CredentialsOptional=false)")
		return fmt.Errorf(errorMsg)
	}

	// Now parse the token
	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	// Check if there was an error in parsing...
	if err != nil {
		log.Debugf("Error parsing token: %v", err)
		return fmt.Errorf("Error parsing token: %v", err)
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		log.Debugf("Error validating token algorithm: %s", message)
		return fmt.Errorf("Error validating token algorithm: %s", message)
	}

	// Check if the parsed token is valid...
	if !parsedToken.Valid {
		log.Debug("Token is invalid")
		return fmt.Errorf("Token is invalid")
	}

	log.Debugf("JWT: %v", parsedToken)

	// If we get here, everything worked and we can set the
	// user property in context.
	ctx.Values().Set(m.Config.ContextKey, parsedToken)

	return nil
}

// NewJwtMiddleware New constructs a new Secure instance with supplied options.
func NewJwtMiddleware(cfg ...mwjwt.Config) *JwtMiddleware {

	var c mwjwt.Config
	if len(cfg) == 0 {
		c = mwjwt.Config{}
	} else {
		c = cfg[0]
	}

	if c.ContextKey == "" {
		c.ContextKey = mwjwt.DefaultContextKey
	}

	if c.ErrorHandler == nil {
		c.ErrorHandler = mwjwt.OnError
	}

	if c.Extractor == nil {
		c.Extractor = mwjwt.FromAuthHeader
	}

	return &JwtMiddleware{mwjwt.Middleware{Config: c}}
}
