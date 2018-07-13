package mapstruct

import "github.com/mitchellh/mapstructure"

func Decode(to interface{}, from interface{}) error  {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           to,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(from)
	if err != nil {
		return err
	}

	return nil
}
