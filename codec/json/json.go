package json

import (
	"encoding/json"

	"github.com/kingstonduy/go-core/codec"
)

type JsonCodec struct{}

func NewJsonCodec() codec.Codec {
	return &JsonCodec{}
}

func (j *JsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonCodec) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (j *JsonCodec) String() string {
	return "json"
}
