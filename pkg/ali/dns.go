// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package ali

import (
	"fmt"

	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

type DNSClient struct {
	Client dns.Client
}

// NewDNSClient 创建DNS客户端
func NewDNSClient(accessKeyId, accessKeySecret string) (*DNSClient, error) {
	client, err := dns.NewClient(&openapi.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
	})
	if err != nil {
		return nil, err
	}

	return &DNSClient{
		Client: *client,
	}, nil
}

func (c *DNSClient) FindTextRecords(fqdn string) {

}

// GetDomainRecords 获取域名解析记录
// domain: 域名
func (c *DNSClient) GetDomainRecords(domain string) (records []*dns.DescribeDomainRecordsResponseBodyDomainRecordsRecord, err error) {
	var size int64 = 500
	req := &dns.DescribeDomainRecordsRequest{
		DomainName: &domain,
		PageNumber: tea.Int64(1),
		PageSize:   &size,
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("GetDomainRecords panic: %v", r)
			return
		}
	}()

	for {
		var resp *dns.DescribeDomainRecordsResponse
		resp, err = c.Client.DescribeDomainRecords(req)
		if err != nil {
			return
		}

		body := resp.Body
		if body != nil {
			records = append(records, body.DomainRecords.Record...)
			if *req.PageNumber*size >= *body.TotalCount {
				break
			}
		} else {
			break
		}

		*req.PageNumber += 1
	}
	return
}

// AddResolve 添加解析记录
// domain: 域名
// recordType: 记录类型
// hostname: 主机记录
// value: 记录值
func (c *DNSClient) AddResolve(domain, recordType, hostname, value string) (recordId *string, err error) {
	req := &dns.AddDomainRecordRequest{
		DomainName: &domain,
		Type:       &recordType,
		RR:         &hostname,
		Value:      &value,
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("AddResolve panic: %v", r)
			return
		}
	}()

	var resp *dns.AddDomainRecordResponse
	resp, err = c.Client.AddDomainRecord(req)
	if err != nil {
		return
	}

	if resp.Body != nil {
		recordId = resp.Body.RecordId
	}
	return
}

// DeleteResolve 删除解析记录
// domain: 域名
// rr: 主机记录
func (c *DNSClient) DeleteResolve(recordId *string) (resp *dns.DeleteDomainRecordResponse, err error) {
	req := &dns.DeleteDomainRecordRequest{
		RecordId: recordId,
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("DeleteResolve panic: %v", r)
			return
		}
	}()

	resp, err = c.Client.DeleteDomainRecord(req)
	if err != nil {
		return
	}

	return
}
