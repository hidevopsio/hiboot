package swagger

import "hidevops.io/hiboot/pkg/at"

type properties struct {
	at.ConfigurationProperties `value:"swagger"`

	UI struct{
		PathPrefix string `json:"path_prefix" default:"/swagger-ui"`
	} `json:"ui"`
}