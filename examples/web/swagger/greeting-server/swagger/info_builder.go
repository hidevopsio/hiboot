package swagger

import "github.com/go-openapi/spec"

// TODO: no need to use builder anymore according to the performance ?
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

func OpenAPIDefinitionBuilder() OpenAPIDefinitionBuilderInterface {
	return &OpenAPIDefinition{
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

func (b *OpenAPIDefinition) Title(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Title = value
	return b
}

func (b *OpenAPIDefinition) Description(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Description = value
	return b
}

func (b *OpenAPIDefinition) TermsOfServiceUrl(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.TermsOfService = value
	return b
}

func (b *OpenAPIDefinition) Version(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Version = value
	return b
}

func (b *OpenAPIDefinition) Host(values string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Host = values
	return b
}

func (b *OpenAPIDefinition) BasePath(values string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.BasePath = values
	return b
}

func (b *OpenAPIDefinition) Schemes(values ...string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Schemes = values
	return b
}

func (b *OpenAPIDefinition) ContactName(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.Name = value
	return b
}

func (b *OpenAPIDefinition) ContactEmail(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.Email = value
	return b
}

func (b *OpenAPIDefinition) ContactURL(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.URL = value
	return b
}

func (b *OpenAPIDefinition) Contact(value Contact) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.Contact.Name = value.Name
	b.Swagger.SwaggerProps.Info.Contact.Email = value.Email
	b.Swagger.SwaggerProps.Info.Contact.URL = value.URL
	return b
}

func (b *OpenAPIDefinition) LicenseName(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.License.Name = value
	return b
}

func (b *OpenAPIDefinition) LicenseURL(value string) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.License.URL = value
	return b
}

func (b *OpenAPIDefinition) License(value License) OpenAPIDefinitionBuilderInterface {
	b.Swagger.SwaggerProps.Info.License.Name = value.Name
	b.Swagger.SwaggerProps.Info.License.URL = value.URL
	return b
}
