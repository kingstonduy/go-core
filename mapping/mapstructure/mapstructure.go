package mapstructure

import (
	"github.com/kingstonduy/go-core/mapping"
	"github.com/mitchellh/mapstructure"
)

type MapStructure struct {
}

func NewMapStructure() mapping.Mapper {
	return &MapStructure{}
}

func (m *MapStructure) Map(input interface{}, output interface{}, opts ...mapping.MapperOption) error {
	options := mapping.NewMapperOptions(opts...)
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: options.WeaklyTypedInput,
		Result:           &output,
	}

	fm := options.DecodeTimeFormat
	if len(fm) > 0 {
		if config.DecodeHook != nil {
			config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
				config.DecodeHook,
				mapstructure.StringToTimeHookFunc(fm),
			)
		} else {
			config.DecodeHook = mapstructure.StringToTimeHookFunc(fm)
		}
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(input)
	return err
}

type MapStructOption func(*mapstructure.DecoderConfig)
