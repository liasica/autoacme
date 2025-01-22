// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package logger

import "go.uber.org/zap"

type Logger struct {
	skip int
}

func NewLogger(skip int) *Logger {
	return &Logger{skip: skip}
}

func (l *Logger) Fatal(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(l.skip)).Fatal(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(l.skip)).Fatalln(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(l.skip)).Fatalf(format, args...)
}

func (l *Logger) Print(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(l.skip)).Info(args...)
}

func (l *Logger) Println(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(l.skip)).Infoln(args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(l.skip)).Infof(format, args...)
}
