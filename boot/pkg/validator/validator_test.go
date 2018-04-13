package validator


import (
	"gopkg.in/go-playground/validator.v8"
	"fmt"
	"testing"
)

// User contains user information
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"hexcolor|rgb|rgba"`
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func TestValidateStruct(t *testing.T) {

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000",
		Addresses:      []*Address{address},
	}

	// returns nil or ValidationErrors ( map[string]*FieldError )
	errs := Validate(user)

	if errs != nil {

		fmt.Println(errs) // output: Key: "User.Age" Error:Field validation for "Age" failed on the "lte" tag
		//	                         Key: "User.Addresses[0].City" Error:Field validation for "City" failed on the "required" tag
		err := errs.(validator.ValidationErrors)["User.Addresses[0].City"]
		fmt.Println(err.Field) // output: City
		fmt.Println(err.Tag)   // output: required
		fmt.Println(err.Kind)  // output: string
		fmt.Println(err.Type)  // output: string
		fmt.Println(err.Param) // output:
		fmt.Println(err.Value) // output:

		// from here you can create your own error messages in whatever language you wish
		return
	}

	// save user to database
}

func TestValidateField(t *testing.T) {
	myEmail := "joeybloggs.gmail.com"

	errs := ValidateField(myEmail, "required,email")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}

	// email ok, move on
}