// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"github.com/kataras/iris/context"
	"reflect"
)

// should always print "($PATH) Handler is executing from 'Context'"
func recordWhichContextJustForProofOfConcept(ctx context.Context) {
	ctx.Application().Logger().Infof("(%s) Handler is executing from: '%s'", ctx.Path(), reflect.TypeOf(ctx).Elem().Name())
	ctx.Next()
}

//
//func TestCustomContext(t *testing.T) {
//	log.Debug("TestCustomContext")
//	app := iris.New()
//	// app.Logger().SetLevel("debug")
//
//	// The only one Required:
//	// here is how you define how your own context will
//	// be created and acquired from the iris' generic context pool.
//	app.ContextPool.Attach(func() context.Context {
//		return &Context{
//			// Optional Part 3:
//			Context: context.NewContext(app),
//		}
//	})
//
//	e := httptest.New(t, app)
//
//	e.Request("GET", "/health").Expect().Status(http.StatusOK).Body()
//
//}


