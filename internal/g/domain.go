// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package g

import "time"

type DomainProvider string

const (
	ProviderDNS  DomainProvider = "DNS"
	ProviderHTTP DomainProvider = "HTTP"
	ProviderOSS  DomainProvider = "OSS"
)

type DomainHookName string

const (
	DomainHookNameCDN      DomainHookName = "CDN"       // 上传至CDN
	DomainHookNameQiniuSSL DomainHookName = "QINIU-SSL" // 上传至七牛云SSL
)

// 赋予权限列表:
//   - AliyunOSSFullAccess
//   - AliyunYundunCertFullAccess

type Domain struct {
	Domain string // 域名

	Timeout  time.Duration
	Interval time.Duration

	Provider     DomainProvider // 申请证书方式, DNS / HTTP / OSS
	DNSProvider  *DomainProviderDNSConfig
	OSSProvider  *DomainProviderOSSConfig
	HTTPProvider *DomainProviderHTTPConfig

	Hooks []*DomainHook
}

// DomainProviderDNSConfig DNS配置
// 需要分配权限: AliyunDNSFullAccess <dns:DescribeDomains, dns:DeleteDomainRecord, dns:AddDomainRecord>
type DomainProviderDNSConfig struct {
	AccessKeyId     string
	AccessKeySecret string
}

// DomainProviderOSSConfig OSS配置
// 需要分配权限: AliyunOSSFullAccess
type DomainProviderOSSConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string // bucket
	Endpoint        string // endpoint
	Path            string // 校验文件存放路径
}

// DomainProviderHTTPConfig HTTP配置
type DomainProviderHTTPConfig struct {
	Path string // 校验文件存放路径
}

type DomainHook struct {
	Name         DomainHookName            // 钩子类型
	CDNHook      *DomainHookCDNConfig      // CDN配置
	QiniuSSLHook *DomainHookQiniuSSLConfig // 七牛云SSL配置
}

// DomainHookCDNConfig CDN配置
// 需要分配权限: AliyunCDNFullAccess <cdn:SetDomainServerCertificate>
type DomainHookCDNConfig struct {
	AccessKeyId     string
	AccessKeySecret string
}

// DomainHookQiniuSSLConfig 七牛云SSL配置
type DomainHookQiniuSSLConfig struct {
	AK string
	SK string
}
