package utils

import (
	"regexp"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate
var once sync.Once

func NewValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		_ = validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
			email := fl.Field().String()

			regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			match, _ := regexp.MatchString(regex, email)
			return match
		})
	})

	return validate
}

func ValidatorErrors(err error) map[string]error {
	fields := map[string]error{}

	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err
	}

	return fields
}
