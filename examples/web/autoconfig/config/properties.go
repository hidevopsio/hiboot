package config

import "hidevops.io/hiboot/pkg/at"

type properties struct {
	at.ConfigurationProperties `value:"config"`

	Name string `json:"name" default:"foo"`

	DefaultName string `json:"default_name" default:"foo"`
}