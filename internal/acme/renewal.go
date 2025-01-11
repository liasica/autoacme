// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package acme

import (
	"crypto/x509"
	"time"

	"github.com/go-acme/lego/v4/certificate"
	"go.uber.org/zap"
)

func GetRenewalTime(cert *x509.Certificate) (renewal bool) {
	client, err := SetupClient()
	if err != nil {
		zap.S().Errorf("failed to create client: %v", err)
		return
	}

	var info *certificate.RenewalInfoResponse
	info, err = client.Certificate.GetRenewalInfo(certificate.RenewalInfoRequest{Cert: cert})
	if err != nil {
		zap.S().Errorf("failed to get renewal info: %v", err)
		return
	}

	t := info.ShouldRenewAt(time.Now(), 0)
	if t == nil {
		return
	}
	renewal = time.Now().After(*t)

	return
}
