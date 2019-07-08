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

// Package idgen provides unique id generation utilities that use twitter's snowflake algorithm
package idgen

import (
	"bytes"
	"fmt"
	"github.com/sony/sonyflake"
	"net"
	"strconv"
	"strings"
)

var sf *sonyflake.Sonyflake

func init() {
	macAddr := getMacAddr()
	st := sonyflake.Settings{
		MachineID: func() (uint16, error) {
			ma := strings.Split(macAddr, ":")
			mid, err := strconv.ParseInt(ma[0]+ma[1], 16, 64)
			return uint16(mid), err
		},
	}
	sf = sonyflake.NewSonyflake(st)
}

func getMacAddr() (addr string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				// Don't use random as we have a real address
				addr = i.HardwareAddr.String()
				break
			}
		}
	}
	return
}

// Next generates next id as an uint64
func Next() (id uint64, err error) {
	var i uint64
	if sf != nil {
		i, err = sf.NextID()
		if err == nil {
			id = i
		}
	}
	return
}

// NextString generates next id as a string
func NextString() (id string, err error) {
	var i uint64
	i, err = Next()
	if err == nil {
		id = fmt.Sprintf("%d", i)
	}
	return
}
