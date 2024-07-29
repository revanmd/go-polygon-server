package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) (bool, string) {
	err := validate.Struct(s)
	if err != nil {
		var sb strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString(err.StructField())
			sb.WriteString(" is ")
			sb.WriteString(err.Tag())
			sb.WriteString(", ")
		}
		return false, sb.String()
	}
	return true, ""
}
