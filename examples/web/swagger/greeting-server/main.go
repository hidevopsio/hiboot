//go:generate statik -src=./dist

package main

import (
	"github.com/go-openapi/spec"
	_ "hidevops.io/hiboot/examples/web/swagger/greeting-server/controller"
	"hidevops.io/hiboot/examples/web/swagger/greeting-server/swagger"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/starter/actuator"
	"hidevops.io/hiboot/pkg/starter/logging"
)

//run http://localhost:8080/api/v1/greeting-server/swagger-ui to open swagger ui
func main() {
	web.NewApplication().
		// HiBoot profiles
		SetProperty(app.ProfilesInclude, actuator.Profile, logging.Profile).
		// server context path
		SetProperty("server.context_path", "/api/v1/greeting-server").
		// open api definitions
		SetProperty(swagger.Profile, &swagger.OpenAPIDefinition{
			Swagger: spec.Swagger{
				SwaggerProps: spec.SwaggerProps{
					Swagger: "2.0",
					Info: &spec.Info{
						InfoProps: spec.InfoProps{
							Description: "Greeting Server is an application that demonstrate the usage of Swagger Annotations",
							Title:       "HiBoot Swagger Demo Application - Greeting Server",
							Version:     "v1.0.1",
						},
					},
					Schemes:  []string{"http", "https"},
					Host:     "apps.hidevops.io",
					BasePath: "/api/v1/greeting-server",
					Paths: &spec.Paths{
						Paths: map[string]spec.PathItem{
							"/hello": spec.PathItem{
								PathItemProps: spec.PathItemProps{
									Get: &spec.Operation{
										OperationProps: spec.OperationProps{
											ID: "getGreeting",
											Produces: []string{
												"text/plain",
											},
											Parameters: []spec.Parameter{
												{
													SimpleSchema: spec.SimpleSchema{
														Type: "string",
													},
													ParamProps: spec.ParamProps{
														Description: "defaults to World if not given",
														Name:        "name",
														In:          "query",
														Required:    false,
													},
												},
											},

											Responses: &spec.Responses{
												VendorExtensible: spec.VendorExtensible{
													Extensions: nil,
												},
												ResponsesProps: spec.ResponsesProps{
													Default: nil,
													StatusCodeResponses: map[int]spec.Response{
														200: {
															ResponseProps: spec.ResponseProps{
																Description: "returns a greeting",
																Schema: &spec.Schema{
																	SchemaProps: spec.SchemaProps{
																		Type:        spec.StringOrArray{"string"},
																		Description: "contains the actual greeting as plain text",
																	},
																},
															},
														},
														404: {
															ResponseProps: spec.ResponseProps{
																Description: "greeter is not available",
																Schema: &spec.Schema{
																	SchemaProps: spec.SchemaProps{
																		Type:        spec.StringOrArray{"string"},
																		Description: "Report 'not found' error message",
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}).
		Run()
}
