// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-22, by liasica

package qiniu

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/qiniu/go-sdk/v7/auth"
)

type Qiniu struct {
	ak, sk string
	client *resty.Client
}

type Response interface {
	HasError() bool
	GetError() error
}

type ApiResponse struct {
	Code      int    `json:"code,omitempty"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"error_code"`
}

func (r *ApiResponse) HasError() bool {
	return (r.Code != 0 && r.Code != 200) || r.ErrorCode != ""
}

func (r *ApiResponse) GetError() error {
	if r.HasError() {
		return fmt.Errorf("code = %d, error = %s, error_code = %s", r.Code, r.Error, r.ErrorCode)
	}
	return nil
}

func request(p Response, fn func() (*resty.Response, error)) (err error) {
	var resp *resty.Response
	resp, err = fn()
	if err != nil {
		return
	}
	err = jsoniter.Unmarshal(resp.Body(), p)
	if err != nil {
		return
	}
	if p.HasError() {
		return p.GetError()
	}
	return
}

// New 创建七牛请求实例
func New(ak, sk string) *Qiniu {
	return &Qiniu{
		ak: ak,
		sk: sk,
		client: resty.New().
			SetBaseURL("https://api.qiniu.com").
			SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
				return auth.New(ak, sk).AddToken(auth.TokenQBox, request)
			}),
	}
}
