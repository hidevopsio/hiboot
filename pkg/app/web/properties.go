package web

import (
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
)

const (
	// ViewEnabled is the property for enabling web view
	ViewEnabled = "web.view.enabled"
	// ContextPath is the property for setting web view context path
	ContextPath = "web.view.contextPath"
	// DefaultPage is the property for setting default page
	DefaultPage = "web.view.defaultPage"
	// ResourcePath is the property for setting resource path
	ResourcePath = "web.view.resourcePath"
	// Extension is the property for setting extension
	Extension = "web.view.extension"
)

type view struct {
	// ViewEnabled is the property for enabling web view
	Enabled bool
	// ContextPath is the property for setting web view context path
	ContextPath string `default:"/"`
	// DefaultPage is the property for setting default page
	DefaultPage string `default:"index.html"`
	// ResourcePath is the property for setting resource path
	ResourcePath string `default:"./static"`
	// Extension is the property for setting extension
	Extension string `default:".html"`
}

type properties struct {
	at.ConfigurationProperties `value:"web"`

	// View is the properties for setting web view
	View view
}

func init() {
	app.Register(new(properties))
}
