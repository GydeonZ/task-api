package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FieldError holds a single validation error for a specific field.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator wraps go-playground/validator and provides a convenient Validate method.
type Validator struct {
	v *validator.Validate
}

// New returns a configured Validator instance.
func New() *Validator {
	return &Validator{v: validator.New()}
}

// Validate validates the given struct and returns a slice of FieldErrors when
// there are validation failures, or nil when the struct is valid.
func (val *Validator) Validate(s interface{}) []FieldError {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	var errs []FieldError
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, FieldError{
			Field:   e.Field(),
			Message: fieldMessage(e),
		})
	}
	return errs
}

func fieldMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	case "url":
		return "must be a valid URL"
	default:
		return fmt.Sprintf("failed on '%s' validation", e.Tag())
	}
}
