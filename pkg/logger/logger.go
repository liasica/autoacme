// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package logger

import "go.uber.org/zap"

type Logger struct {
}

func (l Logger) Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

func (l Logger) Fatalln(args ...interface{}) {
	zap.S().Fatalln(args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

func (l Logger) Print(args ...interface{}) {
	zap.S().Info(args...)
}

func (l Logger) Println(args ...interface{}) {
	zap.S().Infoln(args...)
}

func (l Logger) Printf(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}
