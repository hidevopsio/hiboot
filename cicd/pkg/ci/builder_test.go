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

package ci


import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hi/boot/pkg/system"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}

func TestPipelineBuilder(t *testing.T)  {

	log.Debug("Test Pipeline Builder")

	sysCfg := system.Build()
	log.Debug(sysCfg)
	assert.Equal(t, "hi", sysCfg.App.Name)

	cfg := Build("java")
	log.Debug(cfg)
	assert.Equal(t, "java", cfg.Pipeline.Name)
}
