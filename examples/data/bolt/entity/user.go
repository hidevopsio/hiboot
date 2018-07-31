package entity

import "github.com/hidevopsio/hiboot/pkg/model"


type User struct {
	model.RequestBody
	Id string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age int `json:"age"`
}

type UserRequestParams struct {
	model.RequestParams
	Id string `json:"id" validate:"required"`
}
