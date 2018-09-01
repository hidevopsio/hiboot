// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mapstruct

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

// Decode decode (convert) map to struct
func Decode(to interface{}, from interface{}) error {

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
