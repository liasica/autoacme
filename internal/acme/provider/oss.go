// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package provider

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-acme/lego/v4/challenge/http01"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/pkg/ali"
)

type OssProvider struct {
	client *ali.OSSClient
	path   string
}

func NewOssProvider(endpoint, accessKeyId, accessKeySecret, bucket, path string) (*OssProvider, error) {
	client, err := ali.NewOSS(endpoint, accessKeyId, accessKeySecret, bucket)
	if err != nil {
		return nil, fmt.Errorf("oss: failed to create oss client: %w", err)
	}

	return &OssProvider{
		client: client,
		path:   path,
	}, nil
}

func (p *OssProvider) Present(domain, token, keyAuth string) error {
	zap.S().Infof("OSS Present, domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)

	key := strings.Trim(filepath.Join(p.path, http01.ChallengePath(token)), "/")
	err := p.client.PutObject(key, []byte(keyAuth))
	if err != nil {
		return fmt.Errorf("oss: failed to upload bytes: %w", err)
	}
	return nil
}

func (p *OssProvider) CleanUp(domain, token, keyAuth string) error {
	zap.S().Infof("OSS CleanUp, domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)
	key := strings.Trim(filepath.Join(p.path, http01.ChallengePath(token)), "/")
	err := p.client.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("oss: failed to delete object: %w", err)
	}
	return nil
}
