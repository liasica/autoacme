// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package acme

import (
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/internal/acme/provider"
	"github.com/liasica/aliacme/internal/acme/storage"
	"github.com/liasica/aliacme/internal/g"
)

func Obtain(do *g.Domain, request certificate.ObtainRequest) (resource *certificate.Resource) {
	cfg := g.GetConfig()

	client, err := SetupClient()
	if err != nil {
		zap.L().Error("failed to create client", zap.Error(err))
		return
	}

	switch do.Provider {
	case g.ProviderDNS:
		if do.DNSProvider == nil {
			zap.L().Error("DNS provider is not configured")
			return
		}
		var p *provider.DnsProvider
		p, err = provider.NewDnsProvider(do.DNSProvider.AccessKeyId, do.DNSProvider.AccessKeySecret)
		if err != nil {
			zap.L().Error("failed to create DNS provider", zap.Error(err))
			return
		}

		err = client.Challenge.SetDNS01Provider(p, dns01.AddRecursiveNameservers(cfg.Dns))
		if err != nil {
			zap.L().Error("failed to set DNS01 provider", zap.Error(err))
			return
		}
	case g.ProviderOSS:
		if do.OSSProvider == nil {
			zap.L().Error("OSS provider is not configured")
			return
		}
		var p *provider.OssProvider
		p, err = provider.NewOssProvider(do.OSSProvider.Endpoint, do.OSSProvider.AccessKeyId, do.OSSProvider.AccessKeySecret, do.OSSProvider.Bucket, do.OSSProvider.Path)
		if err != nil {
			zap.L().Error("failed to create OSS provider", zap.Error(err))
			return
		}

		err = client.Challenge.SetHTTP01Provider(p)
		if err != nil {
			zap.L().Error("failed to set HTTP01 provider", zap.Error(err))
			return
		}
	case g.ProviderHTTP:
		// TODO: neet to implement
		panic("Not implemented")
	default:
		zap.L().Error("unknown provider")
		return
	}

	resource, err = client.Certificate.Obtain(request)
	if err != nil {
		zap.L().Error("failed to obtain certificate", zap.Error(err))
		return
	}

	// Save the certificate resources to the archive
	var archiveStorage *storage.ArchiveStorage
	archiveStorage, err = storage.NewArchiveStorage()
	if err != nil {
		zap.L().Error("failed to create archive storage", zap.Error(err))
		return
	}
	err = archiveStorage.SaveResource(resource)
	if err != nil {
		zap.L().Error("failed to save resource", zap.Error(err))
		return
	}

	return
}
