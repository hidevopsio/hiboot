package mapstruct

import (
	"github.com/mitchellh/mapstructure"
	"fmt"
)

func Decode(to interface{}, from interface{}) error  {

	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           to,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if from == nil {
		return fmt.Errorf("parameters of mapstruct.Decode must not be nil")
	}

	return decoder.Decode(from)
}
