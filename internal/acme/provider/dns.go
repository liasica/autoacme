// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package provider

import (
	dns "github.com/alibabacloud-go/alidns-20150109/v2/client"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/pkg/ali"
	"github.com/liasica/aliacme/pkg/tld"
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
	zap.S().Infof("DNS Present, domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)

	info := dns01.GetChallengeInfo(domain, keyAuth)
	var l *tld.List
	l, err = tld.Parse(info.EffectiveFQDN)
	if err != nil {
		return
	}

	var records []*dns.DescribeDomainRecordsResponseBodyDomainRecordsRecord
	records, err = p.Client.GetDomainRecords(l.Domain())
	if err != nil {
		return
	}

	for _, record := range records {
		if *record.RR == l.SubDomain() {
			// 删除已存在的记录
			_, err = p.Client.DeleteResolve(record.RecordId)
			if err != nil {
				return
			}
		}
	}

	p.recordId, err = p.Client.AddResolve(l.Domain(), "TXT", l.SubDomain(), info.Value)
	return
}

func (p *DnsProvider) CleanUp(domain, token, keyAuth string) (err error) {
	zap.S().Infof("DNS CleanUp, domain = %s, token = %s, keyAuth = %s", domain, token, keyAuth)

	if p.recordId != nil {
		_, err = p.Client.DeleteResolve(p.recordId)
	}
	return
}
