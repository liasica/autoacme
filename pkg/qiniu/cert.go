// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-22, by liasica

package qiniu

import (
	"github.com/go-resty/resty/v2"
)

const (
	UrlSSLCert = "/sslcert"
)

// Cert 证书
// https://developer.qiniu.com/fusion/8593/interface-related-certificate
type Cert struct {
	Certid           string   `json:"certid,omitempty"` // 证书ID
	Name             string   `json:"name,omitempty"`   // 证书名称
	Uid              int      `json:"uid,omitempty"`
	CommonName       string   `json:"common_name,omitempty"` // 通用名称
	Dnsnames         []string `json:"dnsnames,omitempty"`    // DNS域名
	CreateTime       int      `json:"create_time,omitempty"` // 创建时间
	NotBefore        int      `json:"not_before,omitempty"`  // 生效时间
	NotAfter         int      `json:"not_after,omitempty"`   // 过期时间
	Orderid          string   `json:"orderid,omitempty"`
	ProductShortName string   `json:"product_short_name,omitempty"`
	ProductType      string   `json:"product_type,omitempty"`
	CertType         string   `json:"cert_type,omitempty"`
	Encrypt          string   `json:"encrypt,omitempty"`
	EncryptParameter string   `json:"encryptParameter,omitempty"`
	Enable           bool     `json:"enable,omitempty"`
	ChildOrderId     string   `json:"child_order_id,omitempty"`
	State            string   `json:"state,omitempty"`
	AutoRenew        bool     `json:"auto_renew,omitempty"`
	Renewable        bool     `json:"renewable,omitempty"`
	Pri              string   `json:"pri"` // 证书私钥
	Ca               string   `json:"ca"`  // 证书内容
}

// ListCertsResponse 获取证书列表返回
type ListCertsResponse struct {
	*ApiResponse
	Marker string  `json:"marker,omitempty"`
	Certs  []*Cert `json:"certs,omitempty"`
}

// ListCerts 获取证书列表
func (q *Qiniu) ListCerts() ([]*Cert, error) {
	var res ListCertsResponse
	err := request(&res, func() (*resty.Response, error) {
		return q.client.R().Get(UrlSSLCert)
	})
	if err != nil {
		return nil, err
	}
	return res.Certs, nil
}

// FindCert 根据证书ID查找证书
func (q *Qiniu) FindCert(certID string) (*Cert, error) {
	var res struct {
		*ApiResponse
		Cert *Cert `json:"cert"`
	}
	err := request(&res, func() (*resty.Response, error) {
		return q.client.R().SetResult(&res).Get(UrlSSLCert + "/" + certID)
	})
	if err != nil {
		return nil, err
	}
	return res.Cert, nil
}

// CreateCertRequest 创建证书请求
type CreateCertRequest struct {
	Name       string `json:"name"`        // 证书名称
	CommonName string `json:"common_name"` // 通用名称
	Pri        string `json:"pri"`         // 证书私钥
	Ca         string `json:"ca"`          // 证书内容
}

func (q *Qiniu) CreateCert(req *CreateCertRequest) (string, error) {
	var res struct {
		*ApiResponse
		CertID string `json:"certid"`
	}
	err := request(&res, func() (*resty.Response, error) {
		return q.client.R().SetBody(req).SetResult(&res).Post(UrlSSLCert)
	})
	if err != nil {
		return "", err
	}

	return res.CertID, nil
}
