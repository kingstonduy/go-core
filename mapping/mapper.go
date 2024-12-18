package mapping

var (
	DefaultMapper = newNoopsMapper()
)

func SetDefaultMapper(mapper Mapper) {
	DefaultMapper = mapper
}

func Map(input interface{}, output interface{}, opts ...MapperOption) error {
	return DefaultMapper.Map(input, output, opts...)
}

type Mapper interface {
	Map(input interface{}, output interface{}, opts ...MapperOption) error
}
