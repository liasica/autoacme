// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package acme

import (
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"go.uber.org/zap"

	"github.com/liasica/autoacme/internal/acme/provider"
	"github.com/liasica/autoacme/internal/acme/storage"
	"github.com/liasica/autoacme/internal/g"
)

func Obtain(do *g.Domain, request certificate.ObtainRequest) (resource *certificate.Resource) {
	cfg := g.GetConfig()

	client, err := SetupClient()
	if err != nil {
		zap.S().Errorf("failed to create client: %v", err)
		return
	}

	switch do.Provider {
	case g.ProviderDNS:
		if do.DNSProvider == nil {
			zap.S().Error("DNS provider is not configured")
			return
		}
		var p *provider.DnsProvider
		p, err = provider.NewDnsProvider(do.DNSProvider.AccessKeyId, do.DNSProvider.AccessKeySecret, do.Timeout, do.Interval)
		if err != nil {
			zap.S().Errorf("failed to create DNS provider: %v", err)
			return
		}

		err = client.Challenge.SetDNS01Provider(p, dns01.AddRecursiveNameservers(cfg.Dns))
		if err != nil {
			zap.S().Errorf("failed to set DNS01 provider: %v", err)
			return
		}
	case g.ProviderOSS:
		if do.OSSProvider == nil {
			zap.S().Error("OSS provider is not configured")
			return
		}
		var p *provider.OssProvider
		p, err = provider.NewOssProvider(do.OSSProvider.Endpoint, do.OSSProvider.AccessKeyId, do.OSSProvider.AccessKeySecret, do.OSSProvider.Bucket, do.OSSProvider.Path)
		if err != nil {
			zap.S().Errorf("failed to create OSS provider: %v", err)
			return
		}

		err = client.Challenge.SetHTTP01Provider(p)
		if err != nil {
			zap.S().Errorf("failed to set HTTP01 provider: %v", err)
			return
		}
	case g.ProviderHTTP:
		// TODO: neet to implement
		panic("Not implemented")
	default:
		zap.S().Error("unknown provider")
		return
	}

	resource, err = client.Certificate.Obtain(request)
	if err != nil {
		zap.S().Errorf("failed to obtain certificate: %v", err)
		return
	}

	// Save the certificate resources to the archive
	var archiveStorage *storage.ArchiveStorage
	archiveStorage, err = storage.NewArchiveStorage()
	if err != nil {
		zap.S().Errorf("failed to create archive storage: %v", err)
		return
	}
	err = archiveStorage.SaveResource(resource)
	if err != nil {
		zap.S().Errorf("failed to save resource: %v", err)
		return
	}

	return
}
