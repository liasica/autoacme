// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package ali

import (
	"bytes"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSClient struct {
	*oss.Client
	Bucket *oss.Bucket
}

func NewOSS(endpoint, accessKeyId, accessKeySecret, bucket string) (*OSSClient, error) {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	c := &OSSClient{
		Client: client,
	}

	c.Bucket, err = client.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// UploadBytes 上传文件
func (c *OSSClient) UploadBytes(name string, b []byte) error {
	return c.Bucket.PutObject(name, bytes.NewReader(b))
}
