package goplayaround

import (
	"github.com/go-playground/validator/v10"
	"github.com/kingstonduy/go-core/validation"
)

type ValidatorInstanceKey struct{}

func WithValidator(val *validator.Validate) validation.ValidationOption {
	return validation.SetOption(ValidatorInstanceKey{}, val)
}
