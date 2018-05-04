package system

import "fmt"

type InvalidControllerError struct {
	Name string
}

func (e *InvalidControllerError) Error() string {
	return fmt.Sprintf("the Controller %v is invalid, please add web.Controller as the struct member", e.Name)
}


type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v is not found", e.Name)
}
