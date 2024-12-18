package goplayaround

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/validation"
)

type _validator struct {
	validator *validator.Validate
}

func NewGpValidator(opts ...validation.ValidationOption) validation.Validator {

	options := validation.NewOptions(opts...)

	var val *validator.Validate
	if v := options.Context.Value(ValidatorInstanceKey{}); v != nil {
		if castedVal, ok := v.(*validator.Validate); ok {
			val = castedVal
		}
	}

	if val == nil {
		val = validator.New()
	}

	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &_validator{
		validator: val,
	}
}

func (v *_validator) Validate(i interface{}) error {
	var errs []error
	err := v.validator.Struct(i)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var el validation.ValidationError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			el.Message = err.Error()
			errs = append(errs, &el)
		}
		// err = fmt.Errorf("%v", errs)
	}

	return errorx.ValidationErrorWithDetails(errs, "")
}
