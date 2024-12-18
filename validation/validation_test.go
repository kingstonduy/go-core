package validation

import "testing"

type data struct {
}

func TestDefaultValidator(t *testing.T) {
	if err := Validate(data{}); err != nil {
		t.Error(err)
	}
}
