// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package internal

import (
	"crypto/x509"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/internal/acme"
	"github.com/liasica/aliacme/internal/acme/storage"
	"github.com/liasica/aliacme/internal/g"
	"github.com/liasica/aliacme/internal/hook"
)

type Hub struct {
	archive *storage.ArchiveStorage
}

func New() *Hub {
	archiveStorage, err := storage.NewArchiveStorage()
	if err != nil {
		zap.L().Fatal("Failed to create archive storage")
	}
	return &Hub{
		archive: archiveStorage,
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		h.run()
	}
}

func (h *Hub) run() {
	cfg := g.GetConfig()
	for _, do := range cfg.Domains {
		h.handle(do)
	}
}

func (h *Hub) handle(do *g.Domain) {
	certs, _ := h.archive.ReadCertificate(do.Domain)

	if len(certs) != 2 {
		h.handleHook(do, h.handleNew(do))
		return
	}

	cert := certs[0]
	if time.Now().After(cert.NotAfter.AddDate(0, 0, -7)) {
		zap.L().Info("renew certificate", zap.String("domain", do.Domain))
		h.handleHook(do, h.handleRenewwal(do, cert))
	} else {
		zap.L().Info("skip renew certificate", zap.String("domain", do.Domain))
	}
}

func (h *Hub) handleNew(do *g.Domain) (resource *certificate.Resource) {
	return acme.Obtain(do, certificate.ObtainRequest{
		Domains: []string{do.Domain},
		Bundle:  true,
	})
}

func (h *Hub) handleRenewwal(do *g.Domain, cert *x509.Certificate) *certificate.Resource {
	return nil
}

func (h *Hub) handleHook(do *g.Domain, resource *certificate.Resource) {
	if resource == nil || resource.Certificate == nil || resource.PrivateKey == nil {
		zap.L().Error("failed to obtain certificate")
		return
	}

	hook.NewHook(do, resource.PrivateKey, resource.Certificate).Run()
}
