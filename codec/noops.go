package codec

import "log"

type noopsCodec struct{}

// Marshal implements Codec.
func (n *noopsCodec) Marshal(v interface{}) ([]byte, error) {
	n.noopsWarning()
	return nil, nil
}

// String implements Codec.
func (n *noopsCodec) String() string {
	n.noopsWarning()
	return "noops"
}

// Unmarshal implements Codec.
func (n *noopsCodec) Unmarshal(b []byte, v interface{}) error {
	n.noopsWarning()
	return nil
}

func (n *noopsCodec) noopsWarning() {
	log.Print("[WARN] No default codec was set. Using noops codec as default. Set the default codec to do all functions\n")
}

func newNoopsCodec() Codec {
	return &noopsCodec{}
}
