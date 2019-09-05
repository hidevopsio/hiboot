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

// Package swagger provides the hiboot starter for swagger
package swagger

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

const (
	Profile = "swagger"
)

type configuration struct {
	at.AutoConfiguration

	Properties *properties
}

func newConfiguration(properties *properties) *configuration {
	return &configuration{Properties: properties}
}

func init() {
	app.Register(newConfiguration, new(properties))
}
