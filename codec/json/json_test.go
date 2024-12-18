package json

import (
	"testing"

	"github.com/kingstonduy/go-core/codec"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Data string
}

func TestMarshal(t *testing.T) {
	data := Data{
		Data: "test",
	}

	codec.SetDefaultCodec(NewJsonCodec())
	bytes, err := codec.Marshal(data)

	assert.Nil(t, err)
	assert.Greater(t, len(bytes), 0)
}

func TestUnmarshal(t *testing.T) {
	data := Data{
		Data: "test",
	}

	codec.SetDefaultCodec(NewJsonCodec())
	bytes, err := codec.Marshal(data)

	assert.Nil(t, err)
	assert.Greater(t, len(bytes), 0)

	var result Data
	err = codec.Unmarshal(bytes, &result)
	assert.Nil(t, err)
}
