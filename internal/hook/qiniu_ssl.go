// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-22, by liasica

package hook

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/liasica/autoacme/internal/g"
	"github.com/liasica/autoacme/pkg/qiniu"
)

func (h *Hook) QiniuSSL(cfg *g.DomainHookQiniuSSLConfig) {
	defer func() {
		h.wg.Done()
	}()

	q := qiniu.New(cfg.AK, cfg.SK)

	certID, err := q.CreateCert(&qiniu.CreateCertRequest{
		Name:       fmt.Sprintf("[autoacme] %s", time.Now().Format("20060102150405")),
		CommonName: h.do.Domain,
		Pri:        string(h.privateKey),
		Ca:         string(h.certificate),
	})
	if err != nil {
		zap.S().Errorf("failed to create Qiniu SSL certificate: %v", err)
		return
	}
	zap.S().Infof("Qiniu SSL certificate created successfully, certID: %s", certID)

	err = q.UpdateDomainHttps(h.do.Domain, &qiniu.UpdateDomainHttpsRequest{
		CertID:      certID,
		ForceHttps:  true,
		Http2Enable: true,
	})

	if err != nil {
		zap.S().Errorf("failed to update Qiniu domain HTTPS: %v", err)
		return
	}

	zap.S().Infof("Qiniu domain HTTPS updated successfully for domain: %s", h.do.Domain)
}
