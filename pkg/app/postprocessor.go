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

package app

import "github.com/hidevopsio/hiboot/pkg/inject"

type PostProcessor interface {
	AfterInitialization(factory interface{})
}

type postProcessor struct {
	subscribes []PostProcessor
}

func newPostProcessor() *postProcessor {
	return new(postProcessor)
}

var (
	postProcessors []interface{}
)

func init() {

}

func RegisterPostProcessor(p ...interface{}) {
	postProcessors = append(postProcessors, p...)
}

func (p *postProcessor) Init() {
	for _, processor := range postProcessors {
		ss, err := inject.IntoFunc(processor)
		if err == nil {
			p.subscribes = append(p.subscribes, ss.(PostProcessor))
		}
	}
}

func (p *postProcessor) AfterInitialization(factory interface{}) {
	for _, processor := range p.subscribes {
		inject.IntoFunc(processor)
		processor.AfterInitialization(factory)
	}
}
