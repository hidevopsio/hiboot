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

package main

import (
	"github.com/hidevopsio/hiboot/pkg/starter/web"
	// starter/actuator is response for health check etc.
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	_ "github.com/hidevopsio/hiboot/examples/web/jwt/controller"

)

func main()  {
	// create new web application and run it
	web.NewApplication().Run()
}
