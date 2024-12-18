// Package codec is an interface for encoding messages
package codec

var (
	DefaultCodec = newNoopsCodec()
)

func SetDefaultCodec(codec Codec) {
	DefaultCodec = codec
}

func Marshal(v interface{}) ([]byte, error) {
	return DefaultCodec.Marshal(v)
}

func Unmarshal(b []byte, v interface{}) error {
	return DefaultCodec.Unmarshal(b, v)
}

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(b []byte, v interface{}) error
	String() string
}
