// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package ali

import (
	"fmt"
	"time"

	cdn "github.com/alibabacloud-go/cdn-20180510/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

type CDNClient struct {
	Client *cdn.Client
}

// NewCDNClient 创建CDN客户端
func NewCDNClient(accessKeyId, accessKeySecret string) (client *CDNClient, err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	var c *cdn.Client
	c, err = cdn.NewClient(config)
	if err != nil {
		return
	}

	client = &CDNClient{
		Client: c,
	}
	return
}

// SetDomainServerCertificate 设置加速域名证书
func (c *CDNClient) SetDomainServerCertificate(domain string, private, public string) (response *cdn.SetDomainServerCertificateResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SetDomainServerCertificate panic: %v", r)
			return
		}
	}()

	request := &cdn.SetDomainServerCertificateRequest{
		CertName:                tea.String(domain + "-" + time.Now().Format("20060102150405")),
		DomainName:              &domain,
		ForceSet:                tea.String("1"),
		PrivateKey:              &private,
		ServerCertificate:       &public,
		ServerCertificateStatus: tea.String("on"),
	}

	return c.Client.SetDomainServerCertificate(request)
}
