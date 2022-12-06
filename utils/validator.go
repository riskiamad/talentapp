package utils

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func IsRequestValid(m interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getUnmarshallError(jeType reflect.Type) string {
	switch jeType.String() {
	case "float64", "int", "int64", "int8":
		return "Should be number"
	case "string":
		return "Should be string"
	}
	return "unknown type"
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "numeric":
		return "This field is numeric type"
	case "number":
		return "This field is number type"
	case "url":
		return "Should be url format"
	case "base64":
		return "Should be base64 string format"
	case "eq=REQUIRED|eq=OPTIONAL":
		return "Enum type should be REQUIRED or OPTIONAL"
	case "eq=ACTIVE|eq=INACTIVE":
		return "Enum type should be ACTIVE or INACTIVE"
	case "eq=true|eq=false":
		return "Enum type should be true or false"
	case "alphanum":
		return "Should be alphanumeric without space"
	case "alphanumunicode":
		return "Should be alphanumeric and unicode without space"
	case "required_without":
		return "This field is required without " + fe.Param()
	case "required_with":
		return "This field is required with " + fe.Param()
	case "required_with_all":
		return "This field is required with " + fe.Param()
	case "excluded_with":
		return "This field should empty if " + fe.Param() + " not empty"
	case "excludes":
		if fe.Param() == " " {
			return "Should be not contains whitespace"
		} else {
			return "Should be not contains " + fe.Param()
		}
	}
	return "Unknown error"
}

// CustomValidator : function to validate input request and return custom error message.
func CustomValidator(err error) []ErrorMsg {
	var (
		je *json.UnmarshalTypeError
		ve validator.ValidationErrors
	)

	var out []ErrorMsg

	if errors.As(err, &ve) {
		for _, fe := range ve {
			out = append(out, ErrorMsg{fe.Field(), getErrorMsg(fe)})
		}
	}

	if errors.As(err, &je) {
		out = append(out, ErrorMsg{je.Field, getUnmarshallError(je.Type)})
	}

	if err == nil {
		return nil
	}

	return out
}
