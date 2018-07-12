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


package log

import (
	"testing"
	"github.com/kataras/golog"
)

func TestLogPrint(t *testing.T) {
	Print("testing ...")
}

func TestLogPrintln(t *testing.T) {
	Println("testing ...")
}

func TestLogf(t *testing.T) {
	Logf(golog.DebugLevel,"testing %v", "...")
}

func TestLogDebug(t *testing.T)  {
	SetLevel(DebugLevel)
	Debug("testing ...")
}

func TestLogDebugf(t *testing.T)  {
	SetLevel(DebugLevel)
	Debugf("testing %v", "...")
}


func TestLogInfo(t *testing.T)  {
	SetLevel(DebugLevel)
	Info("testing ...")
}

func TestLogInfof(t *testing.T)  {
	SetLevel(DebugLevel)
	Infof("testing %v", "...")
}

func TestWarnInfo(t *testing.T)  {
	SetLevel(DebugLevel)
	Warn("testing ...")
}

func TestLogWarnf(t *testing.T)  {
	SetLevel(DebugLevel)
	Warnf("testing %v", "...")
}

func TestLogErrorInfo(t *testing.T)  {
	SetLevel(DebugLevel)
	Error("testing ...")
}

func TestLogErrorf(t *testing.T)  {
	SetLevel(DebugLevel)
	Errorf("testing %v", "...")
}