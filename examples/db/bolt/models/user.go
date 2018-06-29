package models


type User struct {
	Id string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age int `json:"age"`
}
