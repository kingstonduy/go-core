package validation

import (
	"log"
)

type noopsValidator struct{}

func newNoopsValidator() Validator {
	return &noopsValidator{}
}

// Validate implements Validator.
func (n *noopsValidator) Validate(obj interface{}) error {
	n.noopsWarning()
	return nil
}

func (n *noopsValidator) noopsWarning() {
	log.Print("[WARN] No default validator was set. Using noops validator as default. Set the default validator to do all functions\n")
}
