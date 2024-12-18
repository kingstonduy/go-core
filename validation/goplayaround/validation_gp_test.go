package goplayaround

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/kingstonduy/go-core/errorx"
	"github.com/kingstonduy/go-core/validation"
	"github.com/stretchr/testify/assert"
)

type Product struct {
	Name        string `json:"name" validate:"required,email"`
	Price       int64  `json:"price" validate:"min=10,max=100"`
	Description string `json:"description" validate:"required"`
}

func TestValidation(t *testing.T) {
	data := &Product{
		Name:        "Product Name",
		Price:       10,
		Description: "Description",
	}

	validator := NewGpValidator()
	err := validator.Validate(data)
	if err != nil {
		errBytes, _ := json.Marshal(err)
		t.Logf("%v", string(errBytes))
	}
}

func TestDefaultValidation(t *testing.T) {
	validation.SetDefaultValidator(NewGpValidator())
	data := &Product{
		Name:        "Product Name",
		Price:       10,
		Description: "Description",
	}

	err := validation.Validate(data)
	if err != nil {
		errBytes, _ := json.MarshalIndent(err, "", "\t")
		t.Logf("%v", string(errBytes))
	}
}

func TestDefaultValidationWithNilValue(t *testing.T) {
	validation.SetDefaultValidator(NewGpValidator())
	err := validation.Validate(nil)
	if err != nil {
		errBytes, _ := json.MarshalIndent(err, "", "\t")
		t.Logf("%v", string(errBytes))
	}
}

func TestWithCustomValidatorInstance(t *testing.T) {
	internalVal := validator.New()

	type Request struct {
		Data string `customTag:"required"`
	}

	internalVal.SetTagName("customTag")

	val := NewGpValidator(
		WithValidator(internalVal),
	)

	err := val.Validate(Request{})

	assert.NotNil(t, err)

	var errx *errorx.Error
	assert.ErrorAs(t, err, &errx)
}
