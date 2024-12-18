package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type source struct {
	Data   string
	Number int
}

type dest struct {
	Data   string
	Number int
}

func TestDefaultMapper(t *testing.T) {
	src := source{
		Data:   "test",
		Number: 9,
	}
	var des dest
	if err := Map(src, &des); err != nil {
		t.Error(err)
	}

	assert.Equal(t, des.Data, "")
	assert.Equal(t, des.Number, 0)
}
