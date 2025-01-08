// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package provider

import (
	"go.uber.org/zap"
)

type OssProvider struct {
}

func (p OssProvider) Present(domain, token, keyAuth string) error {
	zap.S().Infof("OSS Present, domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)
	return nil
}

func (p OssProvider) CleanUp(domain, token, keyAuth string) error {
	zap.S().Infof("OSS CleanUp, domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)
	return nil
}
