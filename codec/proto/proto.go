package proto

import (
	"github.com/kingstonduy/go-core/codec"
	"google.golang.org/protobuf/proto"
)

type ProtoCodec struct{}

func NewProtoCodec() codec.Codec {
	return &ProtoCodec{}
}

func (p *ProtoCodec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (p *ProtoCodec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (p *ProtoCodec) String() string {
	return "proto"
}
