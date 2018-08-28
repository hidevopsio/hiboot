package jwt


import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/app"
)

type postProcessor struct {
	jwtMiddleware *JwtMiddleware
	application app.Application
}

func init() {
	// register postProcessor
	app.RegisterPostProcessor(new(postProcessor))
}

func (p *postProcessor) Init(application app.Application, jwtMiddleware *JwtMiddleware)  {
	p.application = application
	p.jwtMiddleware = jwtMiddleware
}

func (p *postProcessor) BeforeInitialization(factory interface{})  {
	log.Debug("[jwt] BeforeInitialization")
}

func (p *postProcessor) AfterInitialization(factory interface{})  {
	log.Debug("[jwt] AfterInitialization")

	// use jwt
	p.application.Use(p.jwtMiddleware.Serve)

	// finally register jwt controllers
	err := p.application.RegisterController(new(JwtController))
	if err != nil {
		log.Warnf("[jwt] %v", err)
	}
}
