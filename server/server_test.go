package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerStart(t *testing.T) {
	wrapper := NewServerWrapper()
	err := wrapper.Start(context.TODO())
	assert.Nil(t, err)
}

func TestServerStop(t *testing.T) {
	wrapper := NewServerWrapper()
	err := wrapper.Stop(context.TODO())
	assert.Nil(t, err)
}
