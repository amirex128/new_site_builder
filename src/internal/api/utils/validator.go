package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		validate = validator.New(
			validator.WithRequiredStructEnabled(),
			validator.WithPrivateFieldValidation(),
		)
	}

	err := validate.Struct(s)

	return err
}

func Validate(field any, rule string) error {
	if validate == nil {
		validate = validator.New(
			validator.WithRequiredStructEnabled(),
			validator.WithPrivateFieldValidation(),
		)
	}

	return validate.Var(field, rule)
}

func GetValidationErrors(err error) *[]ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			var el ValidationError
			el.Property = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			validationErrors = append(validationErrors, el)
		}
		return &validationErrors
	}

	return nil
}

func IranianMobileNumberValidator(fld validator.FieldLevel) bool {
	iranianMobileNumberPattern := `^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}

	res, err := regexp.MatchString(iranianMobileNumberPattern, value)
	if err != nil {
		log.Print(err.Error())
	}

	return res
}
