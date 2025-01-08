// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package acme

import (
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/internal/acme/provider"
	"github.com/liasica/aliacme/internal/acme/storage"
	"github.com/liasica/aliacme/internal/g"
)

func Obtain(do *g.Domain, request certificate.ObtainRequest) (resource *certificate.Resource) {
	cfg := g.GetConfig()

	// Get accounts storage
	accountsStorage, err := storage.NewAccountsStorage(cfg.Account)
	if err != nil {
		zap.L().Error("Failed to create accounts storage", zap.Error(err))
		return
	}

	// Load or create account
	var user *g.Account
	user, err = accountsStorage.LoadAccount(cfg.Account)
	if err != nil {
		zap.L().Error("Failed to load account", zap.Error(err))
		return
	}

	config := lego.NewConfig(user)

	// TODO: 10分钟超时
	config.Certificate.Timeout = 1 * time.Minute

	// A client facilitates communication with the CA server.
	var client *lego.Client
	client, err = lego.NewClient(config)
	if err != nil {
		zap.L().Error("Failed to create lego client", zap.Error(err))
		return
	}

	needSave := false
	// New users will need to register
	var reg *registration.Resource
	if user.Registration == nil {
		needSave = true
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			zap.L().Error("Failed to register", zap.Error(err))
			return
		}
	} else {
		reg, err = client.Registration.QueryRegistration()
		if err != nil {
			zap.L().Error("Failed to query registration", zap.Error(err))
			return
		}
	}

	// Save user
	user.Registration = reg
	if needSave {
		accountsStorage.Save(user)
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
			zap.L().Error("Failed to create DNS provider", zap.Error(err))
			return
		}

		err = client.Challenge.SetDNS01Provider(p, dns01.AddRecursiveNameservers(cfg.Dns))
		if err != nil {
			zap.L().Error("Failed to set DNS01 provider", zap.Error(err))
			return
		}
	case g.ProviderOSS:
		err = client.Challenge.SetHTTP01Provider(&provider.OssProvider{})
		if err != nil {
			zap.L().Error("Failed to set HTTP01 provider", zap.Error(err))
			return
		}

		err = client.Challenge.SetTLSALPN01Provider(&provider.OssProvider{})
		if err != nil {
			zap.L().Error("Failed to set TLSALPN01 provider", zap.Error(err))
			return
		}
	case g.ProviderHTTP:
		// TODO: neet to implement
		panic("Not implemented")
	default:
		zap.L().Error("Unknown provider")
		return
	}

	resource, err = client.Certificate.Obtain(request)
	if err != nil {
		zap.L().Error("Failed to obtain certificate", zap.Error(err))
		return
	}

	// Save the certificate resources to the archive
	var archiveStorage *storage.ArchiveStorage
	archiveStorage, err = storage.NewArchiveStorage()
	if err != nil {
		zap.L().Error("Failed to create archive storage", zap.Error(err))
		return
	}
	err = archiveStorage.SaveResource(resource)
	if err != nil {
		zap.L().Error("Failed to save resource", zap.Error(err))
		return
	}

	return
}
