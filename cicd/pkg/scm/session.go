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

package scm

import (
	"time"
)

type SessionInterface interface {
	GetSession(baseUrl, username, password string) error
	GetToken() string
}

type Session struct {
	ID               int         `json:"id"`
	Username         string      `json:"username"`
	Email            string      `json:"email"`
	Name             string      `json:"name"`
	PrivateToken     string      `json:"private_token"`
	Blocked          bool        `json:"blocked"`
	CreatedAt        *time.Time  `json:"created_at"`
	IsAdmin          bool        `json:"is_admin"`
	CanCreateGroup   bool        `json:"can_create_group"`
	CanCreateTeam    bool        `json:"can_create_team"`
	CanCreateProject bool        `json:"can_create_project"`
}
