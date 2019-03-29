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
	// set to true or false to filter in included profiles or not
	Filter bool `json:"filter" default:"false"`
	// included profiles
	Include []string `json:"include"`
	// active profile
	Active string `json:"active" default:"${APP_PROFILES_ACTIVE:default}"`
}

type banner struct {
	// disable banner
	Disabled bool `default:"false"`
}

// App is the properties of the application, it hold the base info of the application
type App struct {
	// project name
	Project string `json:"project" default:"hidevopsio"`
	// app name
	Name string `json:"name" default:"${APP_NAME:hiboot-app}"`
	// app description
	Description string `json:"description" default:"${app.name} is a Hiboot Application"`
	// profiles
	Profiles Profiles `json:"profiles"`
	// banner
	Banner banner
	// Version
	Version string `json:"version" default:"${APP_VERSION:v1}"`
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
	// env name
	Name string
	// env value
	Value string
}
