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

// Package log provides logging
package log

import (
	"fmt"
	"io"
	"runtime"

	"github.com/hidevopsio/golog"
	"github.com/hidevopsio/pio"
)

// Available level names are:
// "disable"
// "fatal"
// "error"
// "warn"
// "info"
// "debug"
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	Disable    = "disable"
)

func callerInfo(skip int) (file string, line int, fn string) {
	var pc uintptr
	var ok bool
	if pc, file, line, ok = runtime.Caller(skip); ok {
		fn = runtime.FuncForPC(pc).Name()
	}
	return
}

// TODO: logger should be able to filter out package by name

var withCaller = func(fn func(v ...interface{}), v ...interface{}) {
	argv := make([]interface{}, 1)
	_, line, fnName := callerInfo(3)
	argv[0] = fmt.Sprintf("[%v:%v] ", fnName, line)
	argv = append(argv, v...)

	fn(argv...)
}

var withCallerf = func(fn func(format string, v ...interface{}), format string, v ...interface{}) {
	_, line, fnName := callerInfo(3)
	f := fmt.Sprintf("[%v:%v] %v", fnName, line, format)

	fn(f, v...)
}

// NewLine can override the default package-level line breaker, "\n".
// It should be called (in-sync) before  the print or leveled functions.
//
// See `github.com/hidevopsio/pio#NewLine` and `Logger#NewLine` too.
func NewLine(newLineChar string) {
	golog.NewLine(newLineChar)
}

// Reset re-sets the default logger to an empty one.
func Reset() {
	golog.Reset()
}

// SetOutput overrides the golog.Logger's Printer's output with another `io.Writer`.
func SetOutput(w io.Writer) {
	golog.SetOutput(w)
}

// AddOutput adds one or more `io.Writer` to the golog.Logger's Printer.
//
// If one of the "writers" is not a terminal-based (i.e File)
// then colors will be disabled for all outputs.
func AddOutput(writers ...io.Writer) {
	golog.AddOutput(writers...)
}

// SetPrefix sets a prefix for the default package-level Logger.
//
// The prefix is the first space-separated
// word that is being presented to the output.
// It's written even before the log level text.
//
// Returns itself.
func SetPrefix(s string) *golog.Logger {
	return golog.SetPrefix(s)
}

// SetTimeFormat sets time format for logs,
// if "s" is empty then time representation will be off.
func SetTimeFormat(s string) {
	golog.SetTimeFormat(s)
}

// SetLevel accepts a string representation of
// a `Level` and returns a `Level` value based on that "levelName".
//
// Available level names are:
// "disable"
// "fatal"
// "error"
// "warn"
// "info"
// "debug"
//

// SetLevel alternatively you can use the exported `golog.Level` field, i.e `golog.Level = golog.ErrorLevel`
func SetLevel(levelName string) {
	golog.SetLevel(levelName)
}

// Print prints a log message without levels and colors.
func Print(v ...interface{}) {
	golog.Print(v...)
}

// Println prints a log message without levels and colors.
// It adds a new line at the end.
func Println(v ...interface{}) {
	golog.Println(v...)
}

// Logf prints a leveled log message to the output.
// This method can be used to use custom log levels if needed.
// It adds a new line in the end.
func Logf(level golog.Level, format string, args ...interface{}) {
	golog.Logf(level, format, args...)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(v ...interface{}) {
	withCaller(golog.Fatal, v...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	withCallerf(golog.Fatalf, format, args...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(v ...interface{}) {
	withCaller(golog.Error, v...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) {
	withCallerf(golog.Errorf, format, args...)
}

// Warn will print when logger's Level is warn, info or debug.
func Warn(v ...interface{}) {
	golog.Warn(v...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	golog.Warnf(format, args...)
}

// Info will print when logger's Level is info or debug.
func Info(v ...interface{}) {
	withCaller(golog.Info, v...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	withCallerf(golog.Infof, format, args...)
}

// Debug will print when logger's Level is debug.
func Debug(v ...interface{}) {
	withCaller(golog.Debug, v...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	withCallerf(golog.Debugf, format, args...)
}

// Install receives  an external logger
// and automatically adapts its print functions.
//
// Install adds a golog handler to support third-party integrations,
// it can be used only once per `golog#Logger` instance.
//
// For example, if you want to print using a logrus
// logger you can do the following:
// `golog.Install(logrus.StandardLogger())`
//
// Look `golog#Handle` for more.
func Install(logger golog.ExternalLogger) {
	golog.Install(logger)
}

// InstallStd receives  a standard logger
// and automatically adapts its print functions.
//
// Install adds a golog handler to support third-party integrations,
// it can be used only once per `golog#Logger` instance.
//
// Example Code:
//
//	import "log"
//	myLogger := log.New(os.Stdout, "", 0)
//	InstallStd(myLogger)
//
// Look `golog#Handle` for more.
func InstallStd(logger golog.StdLogger) {
	golog.InstallStd(logger)
}

// Handle adds a log handler to the default logger.
//
// Handlers can be used to intercept the message between a log value
// and the actual print operation, it's called
// when one of the print functions called.
// If it's return value is true then it means that the specific
// handler handled the log by itself therefore no need to
// proceed with the default behavior of printing the log
// to the specified logger's output.
//
// It stops on the handler which returns true firstly.
// The `Log` value holds the level of the print operation as well.
func Handle(handler golog.Handler) {
	golog.Handle(handler)
}

// Hijack adds a hijacker to the low-level logger's Printer.
// If you need to implement such as a low-level hijacker manually,
// then you have to make use of the pio library.
func Hijack(hijacker func(ctx *pio.Ctx)) {
	golog.Hijack(hijacker)
}

// Scan scans everything from "r" and prints
// its new contents to the logger's Printer's Output,
// forever or until the returning "cancel" is fired, once.
func Scan(r io.Reader) (cancel func()) {
	return golog.Scan(r)
}

// Child (creates if not exists and) returns a new child
// Logger based on the default package-level logger instance.
//
// Can be used to separate logs by category.
func Child(name string) *golog.Logger {
	return golog.Child(name)
}
