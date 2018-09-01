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

package entity

import "github.com/hidevopsio/hiboot/pkg/model"

type User struct {
	model.RequestBody
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

type UserRequestParams struct {
	model.RequestParams
	Id string `json:"id" validate:"required"`
}
