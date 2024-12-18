package mapping

import "log"

type noopsMapper struct{}

func newNoopsMapper() Mapper {
	return &noopsMapper{}
}

// Map implements Mapper.
func (n *noopsMapper) Map(input interface{}, output interface{}, opts ...MapperOption) error {
	n.noopsWarning()
	return nil
}

func (n *noopsMapper) noopsWarning() {
	log.Print("[WARN] No default mapper was set. Using noops mapper as default. Set the default mapper to do all functions\n")
}
