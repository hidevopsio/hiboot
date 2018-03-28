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

package impl

import (
	"testing"
	"github.com/hidevopsio/hi/boot/pkg/log"
	"github.com/hidevopsio/hi/cicd/pkg/ci"
	"os"
	"github.com/stretchr/testify/assert"
	"reflect"
	"fmt"
)

func init()  {
	log.SetLevel(log.DebugLevel)
}

func TestJavaPipeline(t *testing.T)  {

	log.Debug("Test Java Pipeline")

	javaPipeline := &JavaPipeline{
		ci.Pipeline{
			App: "test",
			Project: "demo",
		},
	}

	username := os.Getenv("SCM_USERNAME")
	password := os.Getenv("SCM_PASSWORD")
	javaPipeline.Init(&ci.Pipeline{Name: "java", GitUrl: os.Getenv("SCM_URL")})
	err := javaPipeline.Run(username, password, false)
	assert.Equal(t, nil, err)
}


type Book struct {
	Id    int
	Title string
	Price float32
	Authors []string
}

func TestIterateStruct(t *testing.T) {
	book := Book{Id: 12, Title: "test"}
	e := reflect.ValueOf(&book).Elem()

	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i).Interface()
		fmt.Printf("%v %v %v\n", varName,varType,varValue)
	}
}
