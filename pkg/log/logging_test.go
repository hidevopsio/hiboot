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
	"io"
	"github.com/kataras/pio"
)

func TestNewLine(t *testing.T) {
	NewLine("\n")
}

func TestReset(t *testing.T) {
	Reset()
}

func TestSetOutput(t *testing.T) {
	var w io.Writer
	SetOutput(w)
}

func TestAddOutput(t *testing.T) {
	var w io.Writer
	AddOutput(w)
}

func TestSetPrefix(t *testing.T) {
	SetPrefix("[TEST]")
}

func TestLogPrint(t *testing.T) {
	//Print("testing ...")
}

func TestSetTimeFormat(t *testing.T) {
	SetTimeFormat("[2006-01-02 15:04]")
	//Info("TestSetTimeFormat")
}

func TestLogPrintln(t *testing.T) {
	//Println("testing ...")
}

func TestLogf(t *testing.T) {
	//Logf(golog.DebugLevel,"testing %v", "...")
}

func TestLogDebug(t *testing.T)  {
	SetLevel(DebugLevel)
	//Debug("testing ...")
}

func TestLogDebugf(t *testing.T)  {
	SetLevel(DebugLevel)
	//Debugf("testing %v", "...")
}


func TestLogInfo(t *testing.T)  {
	SetLevel(DebugLevel)
	//Info("testing ...")
}

func TestLogInfof(t *testing.T)  {
	SetLevel(DebugLevel)
	//Infof("testing %v", "...")
}

func TestWarnInfo(t *testing.T)  {
	SetLevel(DebugLevel)
	//Warn("testing ...")
}

func TestLogWarnf(t *testing.T)  {
	SetLevel(DebugLevel)
	//Warnf("testing %v", "...")
}

func TestLogErrorInfo(t *testing.T)  {
	SetLevel(DebugLevel)
	//Error("testing ...")
}

func TestLogErrorf(t *testing.T)  {
	SetLevel(DebugLevel)
	//Errorf("testing %v", "...")
}

func TestInstall(t *testing.T) {
	var l golog.ExternalLogger
	Install(l)
}

func TestInstallStd(t *testing.T) {
	var l golog.StdLogger
	InstallStd(l)
}

func TestHandle(t *testing.T) {
	Handle(func(value *golog.Log) (handled bool) {
		return true
	})
}

func TestHijack(t *testing.T) {
	Hijack(func(ctx *pio.Ctx) {
		//Info("Hijack")
	})
}

func TestScan(t *testing.T) {
	//var r io.Reader
	//Scan(r)
}

func TestChild(t *testing.T) {
	//Child("TestChild")
}

func TestFatal(t *testing.T) {

}

func TestFatalf(t *testing.T) {

}