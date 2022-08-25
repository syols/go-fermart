package pkg

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joeljunstrom/go-luhn"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("luhn", func(fl validator.FieldLevel) bool {
		if number, ok := fl.Field().Interface().(string); ok {
			return luhn.Valid(number)
		}
		return false
	})

	if err != nil {
		log.Fatalf(err.Error())
	}
}

func Validate(s interface{}) error {
	return validate.Struct(s)
}
