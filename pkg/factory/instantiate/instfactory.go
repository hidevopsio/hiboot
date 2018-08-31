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

package instantiate

import (
	"github.com/hidevopsio/hiboot/pkg/utils/cmap"
	"fmt"
	"errors"
)

var (
	NotInitializedError = errors.New("InstantiateFactory is not initialized")
)

type InstantiateFactory struct {
	instanceMap cmap.ConcurrentMap
}

func (f *InstantiateFactory) Initialize(instanceMap cmap.ConcurrentMap)  {
	f.instanceMap = instanceMap
}

func (f *InstantiateFactory) Initialized() bool  {
	return f.instanceMap != nil
}

func (f *InstantiateFactory) SetInstance(name string, instance interface{}) (err error) {
	if !f.Initialized() {
		return NotInitializedError
	}

	if _, ok := f.instanceMap.Get(name); ok {
		return fmt.Errorf("instance name %v is already taken", name)
	}
	f.instanceMap.Set(name, instance)
	return
}

func (f *InstantiateFactory) GetInstance(name string) (inst interface{}) {
	if !f.Initialized() {
		return nil
	}
	var ok bool
	if inst, ok = f.instanceMap.Get(name); !ok {
		return nil
	}
	return
}

func (f *InstantiateFactory) Items() map[string]interface{} {
	if !f.Initialized() {
		return nil
	}
	return f.instanceMap.Items()
}