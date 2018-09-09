package fake

import "github.com/kataras/iris/context"

type ApplicationContext struct {
}


func (a *ApplicationContext) RegisterController(controller interface{}) error  {
	return nil
}

func (a *ApplicationContext) Use(handlers ...context.Handler)  {

}

func (a *ApplicationContext) GetProperty(name string) (value interface{}, ok bool) {
	return
}