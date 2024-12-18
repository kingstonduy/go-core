package mapstructure

import (
	"testing"

	"github.com/kingstonduy/go-core/mapping"
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

func TestMapper(t *testing.T) {
	source := source{
		Data:   "data",
		Number: 9,
	}

	var dest dest

	m := NewMapStructure()
	if err := m.Map(source, &dest); err != nil {
		t.Error(err)
	}

	assert.Equal(t, dest.Data, source.Data)
	assert.Equal(t, dest.Number, source.Number)
}

func TestDefaultMapstructure(t *testing.T) {
	source := source{
		Data:   "data",
		Number: 9,
	}

	var dest dest

	m := NewMapStructure()
	mapping.SetDefaultMapper(m)
	if err := mapping.Map(source, &dest); err != nil {
		t.Error(err)
	}

	assert.Equal(t, dest.Data, source.Data)
	assert.Equal(t, dest.Number, source.Number)
	assert.Equal(t, mapping.DefaultMapper, m)
}
