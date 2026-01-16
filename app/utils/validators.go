package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using the validator package
func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// ValidateVar validates a single variable using the validator package
func ValidateVar(field any, tag string) error {
	return validate.Var(field, tag)
}

