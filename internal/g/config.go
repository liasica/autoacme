// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package g

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

var (
	k   = koanf.New(".")
	cfg *Config
)

type Oss struct {
	Bucket   string
	Endpoint string
}

// Config 配置
// 赋予权限列表:
//   - AliyunOSSFullAccess
//   - AliyunYundunCertFullAccess
//   - AliyunDNSFullAccess
//   - AliyunCDNFullAccess
type Config struct {
	Account string
	Dns     []string
	Domains []*Domain
}

func LoadConfig(path string) {
	cfg = &Config{}

	// Load the file.
	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		zap.L().Fatal("load config file failed", zap.Error(err))
	}

	err = k.Unmarshal("", cfg)
	if err != nil {
		zap.L().Fatal("unmarshal config file failed", zap.Error(err))
	}
}

func GetConfig() *Config {
	return cfg
}
