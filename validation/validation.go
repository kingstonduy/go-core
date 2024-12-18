package validation

import (
	"encoding/json"
)

// Default
var (
	DefaultValidator = newNoopsValidator()
)

func SetDefaultValidator(val Validator) {
	DefaultValidator = val
}

func Validate(obj interface{}) error {
	return DefaultValidator.Validate(obj)
}

// interface for validation
type Validator interface {
	Validate(obj interface{}) error
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func (err *ValidationError) Error() string {
	b, _ := json.Marshal(err)
	return string(b)
}
