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

package system

// Profiles is app profiles
// .include auto configuration starter should be included inside this slide
// .active active profile
type Profiles struct {
	Include []string `json:"include"`
	Active  string   `json:"active" default:"${APP_PROFILES_ACTIVE:dev}"`
}

// App is the properties of the application, it hold the base info of the application
type App struct {
	Project  string   `json:"project" default:"hidevopsio"`
	Name     string   `json:"name" default:"hiboot-app"`
	Profiles Profiles `json:"profiles"`
	// TODO: should defined in application-version.yml
	//Version        string   `json:"version" default:"0.0.1"`
}

// Server is the properties of http server
type Server struct {
	Port string `json:"port" default:"8080"`
}

// Logging is the properties of logging
type Logging struct {
	Level string `json:"level" default:"info"`
}

// Env is the name value pair of environment variable
type Env struct {
	Name  string
	Value string
}
