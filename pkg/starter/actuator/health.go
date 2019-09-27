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

// Package controller provide the controller for health check
package actuator

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/factory"
)

// HealthService is the interface for health check
type HealthService interface {
	Name() string
	Status() bool
}

// Health is the health check struct
type Health struct {
	at.Schema `json:"-"`
	Status string `schema:"The status of health check" json:"status"`
}

type healthController struct {
	at.RestController
	at.GetMapping `value:"/health"`

	configurableFactory factory.ConfigurableFactory
}

func init() {
	app.Register(newHealthController)
}

func newHealthController(configurableFactory factory.ConfigurableFactory) *healthController {
	return &healthController{configurableFactory: configurableFactory}
}

// GET /health
func (c *healthController) Get( struct {
	at.GetMapping `value:"/"`
	at.Operation  `id:"health" description:"health check endpoint"`
	at.Produces   `values:"application/json"`
	Responses struct {
		StatusOK struct {
			at.Response `code:"200" description:"Returns the status of health check"`
			Health
		}
	}
}) map[string]interface{} {
	healthServices := c.configurableFactory.GetInstances(at.HealthCheckService{})
	healthCheckProfiles := make(map[string]interface{})

	healthCheckProfiles["status"] = "Up"

	if healthServices != nil {
		for _, md := range healthServices {
			metaData := factory.CastMetaData(md)
			if metaData.Instance != nil {
				healthService := metaData.Instance.(HealthService)
				status := "Down"
				if healthService.Status() {
					status = "Up"
				}
				healthCheckProfiles[healthService.Name()] = Health{
					Status: status,
				}
			}
		}
	}

	return healthCheckProfiles
}
