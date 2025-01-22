// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package provider

import (
	"fmt"

	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"go.uber.org/zap"

	"github.com/liasica/autoacme/pkg/ali"
	"github.com/liasica/autoacme/pkg/tld"
)

type DnsProvider struct {
	Client *ali.DNSClient

	recordId *string
}

func NewDnsProvider(accessKeyId, accessKeySecret string) (p *DnsProvider, err error) {
	var client *ali.DNSClient
	client, err = ali.NewDNSClient(accessKeyId, accessKeySecret)
	if err != nil {
		return
	}
	p = &DnsProvider{
		Client: client,
	}
	return
}

func (p *DnsProvider) Present(domain, token, keyAuth string) (err error) {
	zap.S().Infof("DNS provider: present domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)

	info := dns01.GetChallengeInfo(domain, keyAuth)
	var l *tld.List
	l, err = tld.Parse(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("dns: failed to parse domain: %w", err)
	}

	var records []*dns.DescribeDomainRecordsResponseBodyDomainRecordsRecord
	records, err = p.Client.GetDomainRecords(l.Domain())
	if err != nil {
		return fmt.Errorf("dns: failed to get domain records: %w", err)
	}

	for _, record := range records {
		if *record.RR == l.SubDomain() {
			// 删除已存在的记录
			_, err = p.Client.DeleteResolve(record.RecordId)
			if err != nil {
				return fmt.Errorf("dns: failed to delete domain record: %w", err)
			}
		}
	}

	p.recordId, err = p.Client.AddResolve(l.Domain(), "TXT", l.SubDomain(), info.Value)
	if err != nil {
		return fmt.Errorf("dns: failed to add domain record: %w", err)
	}
	return
}

func (p *DnsProvider) CleanUp(domain, token, keyAuth string) (err error) {
	zap.S().Infof("DNS provider: cleanup domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)

	if p.recordId != nil {
		_, err = p.Client.DeleteResolve(p.recordId)
		if err != nil {
			return fmt.Errorf("dns: failed to delete domain record: %w", err)
		}
	}
	return
}
