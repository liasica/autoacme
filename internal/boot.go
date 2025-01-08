// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package internal

import (
	"os"
	"time"

	"github.com/go-acme/lego/v4/log"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/internal/g"
	"github.com/liasica/aliacme/pkg/logger"
)

var Version string

func Boot(path, ver string) {
	Version = ver

	// 设置全局时区
	tz := "Asia/Shanghai"
	_ = os.Setenv("TZ", tz)
	loc, _ := time.LoadLocation(tz)
	time.Local = loc

	// 设置zap
	l, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(l)
	log.Logger = &logger.Logger{}

	// 创建runtime目录
	if _, err := os.Stat("runtime"); os.IsNotExist(err) {
		_ = os.MkdirAll("runtime", 0755)
	}

	// 打印版本号
	zap.S().Info("version: " + Version)

	// 读取配置文件
	g.LoadConfig(path)
}
