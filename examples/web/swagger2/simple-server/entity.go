package main

import (
	"hidevops.io/hiboot/pkg/at"
)

type Employee struct {
	at.ApiModel     `description:"All details about the Employee. " json:"-"`
	Id        int    `api:"The database generated employee ID" json:"id"`
	FirstName string `api:"The employee first name" json:"first_name"`
	LastName  string `api:"The employee last name" json:"last_name"`
}

type EmployeeResponse struct {
	at.ResponseBody `json:"-"`

	Code            int      `json:"code"`
	Message         string   `json:"message"`
	Data            Employee `json:"data"`
}
