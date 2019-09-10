// Copyright 2018~now John Deng (hi.devops.io@gmail.com).
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

package swagger

import (
	"github.com/go-openapi/spec"
	"hidevops.io/hiboot/pkg/at"
)

type OpenAPIDefinitionBuilderInterface interface {
	Title(value string) OpenAPIDefinitionBuilderInterface
	Description(value string) OpenAPIDefinitionBuilderInterface
	TermsOfServiceUrl(value string) OpenAPIDefinitionBuilderInterface
	Version(value string) OpenAPIDefinitionBuilderInterface
	ContactName(value string) OpenAPIDefinitionBuilderInterface
	ContactEmail(value string) OpenAPIDefinitionBuilderInterface
	ContactURL(value string) OpenAPIDefinitionBuilderInterface
	Contact(value Contact) OpenAPIDefinitionBuilderInterface
	LicenseName(value string) OpenAPIDefinitionBuilderInterface
	LicenseURL(value string) OpenAPIDefinitionBuilderInterface
	License(value License) OpenAPIDefinitionBuilderInterface
	Schemes(values ...string) OpenAPIDefinitionBuilderInterface
	Host(values string) OpenAPIDefinitionBuilderInterface
	BasePath(values string) OpenAPIDefinitionBuilderInterface
}

const (
	Profile = "swagger"
)

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type openAPIDefinition struct {
	at.ConfigurationProperties `value:"swagger"`
	spec.Swagger
}

func OpenAPIDefinitionBuilder() OpenAPIDefinitionBuilderInterface {
	return &openAPIDefinition{
		Swagger: spec.Swagger{
			VendorExtensible: spec.VendorExtensible{},
			SwaggerProps:     spec.SwaggerProps{
				Swagger:             "2.0",
				Info:                &spec.Info{},
				Paths:               &spec.Paths{
					VendorExtensible: spec.VendorExtensible{},
					Paths: make(map[string]spec.PathItem),
				},
			},
		},
	}
}

func (b *openAPIDefinition) Title(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Title = value
	return b
}

func (b *openAPIDefinition) Description(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Description = value
	return b
}

func (b *openAPIDefinition) TermsOfServiceUrl(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.TermsOfService = value
	return b
}

func (b *openAPIDefinition) Version(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Version = value
	return b
}

func (b *openAPIDefinition) Host(values string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Host = values
	return b
}

func (b *openAPIDefinition) BasePath(values string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.BasePath = values
	return b
}

func (b *openAPIDefinition) Schemes(values ...string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Schemes = values
	return b
}

func (b *openAPIDefinition) ContactName(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.Name = value
	return b
}

func (b *openAPIDefinition) ContactEmail(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.Email = value
	return b
}

func (b *openAPIDefinition) ContactURL(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.URL = value
	return b
}

func (b *openAPIDefinition) Contact(value Contact) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.Name = value.Name
	b.Swagger.SwaggerProps.Info.Contact.Email = value.Email
	b.Swagger.SwaggerProps.Info.Contact.URL = value.URL
	return b
}

func (b *openAPIDefinition) LicenseName(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.License.Name = value
	return b
}

func (b *openAPIDefinition) LicenseURL(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.License.URL = value
	return b
}

func (b *openAPIDefinition) License(value License) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.License.Name = value.Name
	b.Swagger.SwaggerProps.Info.License.URL = value.URL
	return b
}
