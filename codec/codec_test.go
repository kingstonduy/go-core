package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	Data string
}

func TestMarshal(t *testing.T) {
	data := Data{
		Data: "test",
	}

	bytes, err := Marshal(data)

	assert.Nil(t, err)
	assert.Equal(t, len(bytes), 0)
}

func TestUnmarshal(t *testing.T) {
	bytes := []byte("test")
	var result string
	err := Unmarshal(bytes, &result)
	assert.Nil(t, err)
	assert.Equal(t, result, "")
}
