// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package internal

import (
	"crypto/x509"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"go.uber.org/zap"

	"github.com/liasica/autoacme/internal/acme"
	"github.com/liasica/autoacme/internal/acme/storage"
	"github.com/liasica/autoacme/internal/g"
	"github.com/liasica/autoacme/internal/hook"
)

type Hub struct {
	archive *storage.ArchiveStorage
}

func New() *Hub {
	archiveStorage, err := storage.NewArchiveStorage()
	if err != nil {
		zap.S().Fatal("failed to create archive storage")
	}
	return &Hub{
		archive: archiveStorage,
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(6 * time.Hour)
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
	if time.Now().After(cert.NotAfter.AddDate(0, 0, -15)) {
		zap.S().Infof("renew certificate: %s", do.Domain)
		h.handleHook(do, h.handleRenewal(do, cert))
	} else {
		zap.S().Infof("skip renew certificate: %s", do.Domain)
	}
}

func (h *Hub) handleNew(do *g.Domain) (resource *certificate.Resource) {
	return acme.Obtain(do, certificate.ObtainRequest{
		Domains: []string{do.Domain},
		Bundle:  true,
	})
}

func (h *Hub) handleRenewal(do *g.Domain, cert *x509.Certificate) *certificate.Resource {
	privateKey, err := h.archive.ReadPrivateKey(do.Domain)
	if err != nil {
		zap.S().Error("failed to read private key", zap.String("domain", do.Domain))
		return nil
	}

	var replacesCertID string
	replacesCertID, err = certificate.MakeARICertID(cert)
	if err != nil {
		zap.S().Error("failed to make certificate ID", zap.String("domain", do.Domain))
		return nil
	}

	return acme.Obtain(do, certificate.ObtainRequest{
		Domains:        []string{do.Domain},
		Bundle:         true,
		PrivateKey:     privateKey,
		ReplacesCertID: replacesCertID,
	})
}

func (h *Hub) handleHook(do *g.Domain, resource *certificate.Resource) {
	if resource == nil || resource.Certificate == nil || resource.PrivateKey == nil {
		zap.S().Error("failed to obtain certificate")
		return
	}

	hook.NewHook(do, resource.PrivateKey, resource.Certificate).Run()
}
