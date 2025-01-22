// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-22, by liasica

package qiniu

import (
	"github.com/go-resty/resty/v2"
)

const (
	UrlDomain = "/domain"
)

type UpdateDomainHttpsRequest struct {
	CertID      string `json:"certId"`
	ForceHttps  bool   `json:"forceHttps"`
	Http2Enable bool   `json:"http2Enable"`
}

// UpdateDomainHttps 修改证书
// https://developer.qiniu.com/fusion/4246/the-domain-name#14
func (q *Qiniu) UpdateDomainHttps(name string, req *UpdateDomainHttpsRequest) error {
	var res ApiResponse
	return request(&res, func() (*resty.Response, error) {
		return q.client.R().SetBody(req).Put(UrlDomain + "/" + name + "/httpsconf")
	})
}
