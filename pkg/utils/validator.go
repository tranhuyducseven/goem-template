package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator is type setting of third party validator
type CustomValidator struct {
	validator *validator.Validate
}

type ValidationRule struct {
	Tag string
	Fn  validator.Func
}

func (cv *CustomValidator) AddRules(rules ...ValidationRule) error {
	for _, r := range rules {
		if err := cv.validator.RegisterValidation(r.Tag, r.Fn); err != nil {
			return err
		}
	}
	return nil
}

// Init validator
func NewValidator() *CustomValidator {
	cv := &CustomValidator{
		validator: validator.New(),
	}
	return cv
}

type ValidationError struct {
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (ve ValidationError) Error() string {
	return ve.Message
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		validateErrs, ok := err.(validator.ValidationErrors)
		if !ok {
			ve := ValidationError{
				Message: err.Error(),
				Metadata: map[string]interface{}{
					"error": "Invalid input",
				},
			}
			return ve
		}
		firstErr := validateErrs[0]
		errMsg := translateError(firstErr)

		validationErr := &ValidationError{
			Message: errMsg,
			Metadata: map[string]interface{}{
				"field":     firstErr.Field(),
				"namespace": firstErr.Namespace(),
				"type":      firstErr.Type().String(),
				"tag":       firstErr.Tag(),
				"param":     firstErr.Param(),
			},
		}

		return validationErr
	}
	return nil
}

// See list of tag here https://github.com/go-playground/validator?tab=readme-ov-file#baked-in-validations
// TODO: Add more validation error message
func translateError(validateErr validator.FieldError) string {
	tag := validateErr.Tag()
	field := pascalCaseToWords(validateErr.Field())

	var msg string = ""
	switch tag {
	case "alpha":
		msg = fmt.Sprintf("%s must contain only letters", field)
	case "alphanum":
		msg = fmt.Sprintf("%s must contain only letters and numbers", field)
	case "alphanumunicode":
		msg = fmt.Sprintf("%s must contain only letters, numbers, and unicode characters", field)
	case "alphaunicode":
		msg = fmt.Sprintf("%s must contain only letters and unicode characters", field)
	case "ascii":
		msg = fmt.Sprintf("%s must contain only ascii characters", field)
	case "contains":
		msg = fmt.Sprintf("%s must contain %s", field, validateErr.Param())
	case "required":
		msg = fmt.Sprintf("%s is required", field)
	case "number", "cidr":
		msg = fmt.Sprintf("%s must be a %s", field, tag)
	case "email", "url", "uri", "uuid":
		msg = fmt.Sprintf("%s is invalid", field)
	case "min":
		switch validateErr.Type().Name() {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			msg = fmt.Sprintf("%s must be greater than or equal to %s", field, validateErr.Param())
		case "string":
			msg = fmt.Sprintf("%s must be at least %s characters long", field, validateErr.Param())
		case "slice", "array":
			msg = fmt.Sprintf("%s must have at least %s items", field, validateErr.Param())
		}
	case "max":
		switch validateErr.Type().Name() {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			msg = fmt.Sprintf("%s must be less than or equal to %s", field, validateErr.Param())
		case "string":
			msg = fmt.Sprintf("%s must be at most %s characters long", field, validateErr.Param())
		case "slice", "array":
			msg = fmt.Sprintf("%s must have at most %s items", field, validateErr.Param())
		}
	case "oneof":
		params := strings.Split(validateErr.Param(), " ")
		paramMsg := ""
		for i, param := range params {
			if i == len(params)-1 {
				paramMsg += fmt.Sprintf("or %s", param)
			} else {
				paramMsg += fmt.Sprintf("%s, ", param)
			}
		}
		msg = fmt.Sprintf("%s must be one of %s", field, paramMsg)
	default:
		msg = fmt.Sprintf("%s is invalid", field)
	}

	if len(msg) == 0 {
		msg = validateErr.Error()
	}

	return msg
}

// Pa
func pascalCaseToWords(s string) string {
	var words []rune
	var prevChar rune
	for i, char := range s {
		tmp := char
		if i > 0 && ('A' <= char && char <= 'Z') {
			if 'a' <= prevChar && prevChar <= 'z' {
				words = append(words, ' ')
			}
			char += 'a' - 'A'
		}
		words = append(words, char)
		prevChar = tmp
	}

	return string(words)
}

func BindAndValidate(c echo.Context, i interface{}) error {
	if err := c.Bind(i); err != nil {
		err, ok := err.(*echo.HTTPError)
		if ok {
			msg := err.Internal.Error()

			flag := strings.Compare(msg, "the provided hex string is not a valid ObjectID") == 0 || strings.Contains(msg, "cannot unmarshal into an ObjectID, the length must be 24")

			if flag {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
			}

			return err
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if err := c.Validate(i); err != nil {
		// not wrapped in echo.NewHTTPError because we want to return the error as is custom ValidationError
		return err
	}
	return nil
}
