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
	"fmt"
	"github.com/hidevopsio/golog"
	"github.com/hidevopsio/pio"
	"os"
	"testing"
	"time"
)

func init() {

}

//func TestSetOutput(t *testing.T) {
//	var w io.Writer
//	SetOutput(w)
//}

//func TestAddOutput(t *testing.T) {
//	var w io.Writer
//	AddOutput(w)
//}

func TestLogging(t *testing.T) {
	t.Run("should pass scan test", func(t *testing.T) {
		f := newLogFile("foo.log")
		defer f.Close()

		SetOutput(f)

		b := newLogFile("bar.log")
		defer f.Close()
		AddOutput(b)

		_ = Scan(os.Stdin)
		// type and enter one or more sentences to your console,
		// wait 10 seconds and open the .txt file.
		<-time.After(1 * time.Second)

		os.Remove("foo.log")
		os.Remove("bar.log")
	})

	t.Run("should pass new line test", func(t *testing.T) {
		NewLine("\n")
	})

	t.Run("should pass reset test", func(t *testing.T) {
		Reset()
	})

	t.Run("should pass set prefix test", func(t *testing.T) {
		SetPrefix("[TEST]")
	})

	t.Run("should pass log.Print() test", func(t *testing.T) {
		Print("testing ...")
	})

	t.Run("should pass log.SetTimeFormat() test", func(t *testing.T) {
		SetTimeFormat("[2006-01-02 15:04:05.000]")
		Info("TestSetTimeFormat")
		time.Sleep(2000000)
		Info("TestSetTimeFormat")
		time.Sleep(1234567)
		Info("TestSetTimeFormat")
	})

	t.Run("should pass log.Println() test", func(t *testing.T) {
		Println("testing ...")
	})

	t.Run("should pass log.Logf() test", func(t *testing.T) {
		Logf(golog.DebugLevel, "testing %v", "...")
	})

	t.Run("should pass log.Debug() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Debug("testing ...")
	})
	t.Run("should pass log.Debugf() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Debugf("testing %v", "...")
	})
	t.Run("should pass log.Info() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Info("testing ...")
	})
	t.Run("should pass log.Infof() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Infof("testing %v", "...")
	})
	t.Run("should pass log.Warn() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Warn("testing ...")
	})
	t.Run("should pass log.Warnf() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Warnf("testing %v", "...")
	})

	t.Run("should pass log.Error() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Error("testing ...")
	})
	t.Run("should pass log.Errorf() test", func(t *testing.T) {
		SetLevel(DebugLevel)
		Errorf("testing %v", "...")
	})

	t.Run("should pass log.() test", func(t *testing.T) {
		var l golog.ExternalLogger
		Install(l)
	})

	t.Run("should pass log.InstallStd() test", func(t *testing.T) {
		var l golog.StdLogger
		InstallStd(l)
	})

	t.Run("should pass log.Handle() test", func(t *testing.T) {
		Handle(func(value *golog.Log) (handled bool) {
			return true
		})
	})

	t.Run("should pass log.Hijack() test", func(t *testing.T) {
		Hijack(func(ctx *pio.Ctx) {
			Info("Hijack")
		})
	})

	t.Run("should pass log.Child() test", func(t *testing.T) {
		Child("TestChild")
	})

	withCaller = func(fn func(v ...interface{}), v ...interface{}) {
		fmt.Println("fake withCaller()")
	}
	t.Run("should pass log.Fatal() test", func(t *testing.T) {
		Fatal("test")
	})

	withCallerf = func(fn func(format string, v ...interface{}), format string, v ...interface{}) {
		fmt.Println("fake withCallerf()")
	}
	t.Run("should pass log.Fatalf() test", func(t *testing.T) {
		Fatalf("test: %v", "log")
	})
}

func newLogFile(filename string) *os.File {
	// open an output file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}
