// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package hook

import (
	cdn "github.com/alibabacloud-go/cdn-20180510/client"
	"go.uber.org/zap"

	"github.com/liasica/autoacme/internal/g"
	"github.com/liasica/autoacme/pkg/ali"
)

func (h *Hook) RunCDN(cfg *g.DomainHookCDNConfig) {
	defer func() {
		h.wg.Done()
	}()

	cdnClient, err := ali.NewCDNClient(cfg.AccessKeyId, cfg.AccessKeySecret)
	if err != nil {
		zap.S().Errorf("failed to create CDN client: %v", err)
	}

	var cdnResponse *cdn.SetDomainServerCertificateResponse
	cdnResponse, err = cdnClient.SetDomainServerCertificate(h.do.Domain, string(h.privateKey), string(h.certificate))
	if err != nil {
		zap.S().Errorf("failed to set domain server certificate: %v", err)
		return
	}

	zap.L().Info("set domain server certificate response", zap.String("domain", h.do.Domain), zap.Reflect("response", cdnResponse))
}
