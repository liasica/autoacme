// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package storage

import (
	"crypto/x509"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"go.uber.org/zap"
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
	rootPath := filepath.Join(baseStoragePath, baseArchiveFolderName)

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
	path := filepath.Join(s.rootPath, resource.Domain)
	err = CreateNonExistingFolder(path)
	if err != nil {
		return
	}

	if resource.IssuerCertificate != nil {
		err = s.WriteFile(s.GetFileName(path, issuerFile), resource.IssuerCertificate)
		if err != nil {
			zap.L().Error("Failed to write issuer certificate", zap.Error(err))
			return
		}
	}

	if resource.PrivateKey != nil {
		err = s.WriteFile(s.GetFileName(path, privateFile), resource.PrivateKey)
		if err != nil {
			zap.L().Error("Failed to write private key", zap.Error(err))
			return
		}

		err = s.WriteFile(s.GetFileName(path, certFile), resource.Certificate)
		if err != nil {
			zap.L().Error("Failed to write certificate", zap.Error(err))
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

func (s *ArchiveStorage) ReadFile(domain, extension string) ([]byte, error) {
	return os.ReadFile(s.GetFileName(domain, extension))
}

func (s *ArchiveStorage) ReadCertificate(domain string) ([]*x509.Certificate, error) {
	content, err := s.ReadFile(domain, certFile)
	if err != nil {
		return nil, err
	}

	// The input may be a bundle or a single certificate.
	return certcrypto.ParsePEMBundle(content)
}
