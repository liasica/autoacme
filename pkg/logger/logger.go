// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package logger

import "go.uber.org/zap"

type Logger struct {
}

func (l Logger) Fatal(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(1)).Fatal(args...)
}

func (l Logger) Fatalln(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(1)).Fatalln(args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(1)).Fatalf(format, args...)
}

func (l Logger) Print(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(1)).Info(args...)
}

func (l Logger) Println(args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(1)).Infoln(args...)
}

func (l Logger) Printf(format string, args ...interface{}) {
	zap.S().WithOptions(zap.AddCallerSkip(1)).Infof(format, args...)
}
