// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package storage

import (
	"crypto"
	"crypto/x509"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/internal/g"
)

const (
	baseArchiveFolderName = "archive"
	certFile              = "fullchain.pem"
	privateFile           = "privkey.pem"
	issuerFile            = "issuer.pem"
)

// ArchiveStorage a certificates' storage.
//
// rootPath:
//
//	./runtime/archive/
//	      │      └── archived certificates directory (rootPath)
//	      └── "path" option
//
// rootUserPath:
//
//	./runtime/archive/example.com/
//	      │      │        └── domain path ("domain" option)
//	      │      └── archived certificates directory (rootPath)
//	      └── "path" option
type ArchiveStorage struct {
	rootPath string
}

func NewArchiveStorage() (s *ArchiveStorage, err error) {
	rootPath := filepath.Join(g.StoragePath, baseArchiveFolderName)

	err = CreateNonExistingFolder(rootPath)
	if err != nil {
		return
	}

	s = &ArchiveStorage{
		rootPath: rootPath,
	}

	return
}

func (s *ArchiveStorage) SaveResource(resource *certificate.Resource) (err error) {
	err = CreateNonExistingFolder(filepath.Join(s.rootPath, resource.Domain))
	if err != nil {
		return
	}

	if resource.IssuerCertificate != nil {
		err = s.WriteFile(s.GetFileName(resource.Domain, issuerFile), resource.IssuerCertificate)
		if err != nil {
			zap.L().Error("failed to write issuer certificate", zap.Error(err))
			return
		}
	}

	if resource.PrivateKey != nil {
		err = s.WriteFile(s.GetFileName(resource.Domain, privateFile), resource.PrivateKey)
		if err != nil {
			zap.L().Error("failed to write private key", zap.Error(err))
			return
		}

		err = s.WriteFile(s.GetFileName(resource.Domain, certFile), resource.Certificate)
		if err != nil {
			zap.L().Error("failed to write certificate", zap.Error(err))
		}
	}

	return
}

func (s *ArchiveStorage) GetFileName(domain, filename string) string {
	return filepath.Join(s.rootPath, domain, filename)
}

func (s *ArchiveStorage) WriteFile(filename string, data []byte) error {
	return os.WriteFile(filepath.Join(s.rootPath, filename), data, 0o600)
}

func (s *ArchiveStorage) ReadFile(domain, filename string) ([]byte, error) {
	return os.ReadFile(s.GetFileName(domain, filename))
}

func (s *ArchiveStorage) ReadPrivateKey(domain string) (crypto.PrivateKey, error) {
	b, err := s.ReadFile(domain, privateFile)
	if err != nil {
		return nil, err
	}

	return certcrypto.ParsePEMPrivateKey(b)
}

func (s *ArchiveStorage) ReadCertificate(domain string) ([]*x509.Certificate, error) {
	content, err := s.ReadFile(domain, certFile)
	if err != nil {
		return nil, err
	}

	// The input may be a bundle or a single certificate.
	return certcrypto.ParsePEMBundle(content)
}
