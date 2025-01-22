// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package storage

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/lego"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/liasica/autoacme/internal/g"
	"github.com/liasica/autoacme/key"
)

const (
	baseAccountsRootFolderName = "accounts"
	accountFileName            = "account.json"
	accountPrivateKeyFileName  = "private.key"
)

// AccountsStorage A storage for account data.
//
// rootUserPath:
//
//	./runtime/hubert@hubert.com/
//	                 └── userID ("email" option)
//
// privateKeyFilePath:
//
//	./runtime/hubert@hubert.com/private.key
//	                 │           └── account private key
//	                 └── userID ("email" option)
//
// accountFilePath:
//
//	./runtime/hubert@hubert.com/account.json
//	                 │             └── account file
//	                 └── userID ("email" option)

// AccountsStorage A storage for account data.
//
// rootPath:
//
//	./runtime/accounts/
//	       │      └── root accounts directory
//	       └── "path" option
//
// rootUserPath:
//
//	./runtime/accounts/localhost_14000/hubert@hubert.com/
//	       │      │             │             └── userID ("email" option)
//	       │      │             └── CA server ("server" option)
//	       │      └── root accounts directory
//	       └── "path" option
//
// privateKeyPath:
//
//	./runtime/accounts/localhost_14000/hubert@hubert.com/private.key
//	       │      │             │             │             └── private key file (ECDSA)
//	       │      │             │             └── userID ("email" option)
//	       │      │             └── CA server ("server" option)
//	       │      └── root accounts directory
//	       └── "path" option
//
// accountFilePath:
//
//	./runtime/accounts/localhost_14000/hubert@hubert.com/account.json
//	       │      │             │             │             └── account file
//	       │      │             │             └── userID ("email" option)
//	       │      │             └── CA server ("server" option)
//	       │      └── root accounts directory
//	       └── "path" option
type AccountsStorage struct {
	userID          string
	rootPath        string
	rootUserPath    string
	privateKeyPath  string
	accountFilePath string
}

func NewAccountsStorage(email string) (s *AccountsStorage, err error) {
	rootPath := filepath.Join(g.StoragePath, baseAccountsRootFolderName)
	serverURL, _ := url.Parse(lego.LEDirectoryProduction)
	serverPath := serverURL.Host
	accountsPath := filepath.Join(rootPath, serverPath)
	rootUserPath := filepath.Join(accountsPath, email)

	err = CreateNonExistingFolder(rootUserPath)
	if err != nil {
		return
	}

	s = &AccountsStorage{
		userID:          email,
		rootPath:        rootPath,
		rootUserPath:    rootUserPath,
		privateKeyPath:  filepath.Join(rootUserPath, accountPrivateKeyFileName),
		accountFilePath: filepath.Join(rootUserPath, accountFileName),
	}

	return
}

func (s *AccountsStorage) AccountFileExists() bool {
	_, err := os.Stat(s.accountFilePath)
	return err == nil
}

func (s *AccountsStorage) PrivateKeyFileExists() bool {
	_, err := os.Stat(s.privateKeyPath)
	return err == nil
}

func (s *AccountsStorage) SaveAccountFile(data []byte) error {
	return os.WriteFile(s.accountFilePath, data, 0644)
}

func (s *AccountsStorage) SavePrivateKeyFile(privite *ecdsa.PrivateKey) error {
	b := key.EncodePrivateKey(privite)
	return os.WriteFile(s.privateKeyPath, b, 0644)
}

func (s *AccountsStorage) LoadPrivateKeyFile() (*ecdsa.PrivateKey, error) {
	b, err := os.ReadFile(s.privateKeyPath)
	if err != nil {
		return nil, err
	}
	return key.DecodePrivateKey(b)
}

// Save store account data to the storage.
func (s *AccountsStorage) Save(account *g.Account) {
	b, _ := jsoniter.MarshalIndent(account, "", "  ")
	err := os.WriteFile(s.accountFilePath, b, os.ModePerm)
	if err != nil {
		zap.S().Errorf("failed to save account file: %v", err)
	}

	err = s.SavePrivateKeyFile(account.Key)
	if err != nil {
		zap.S().Errorf("failed to save private key file: %v", err)
	}
}

// LoadAccount Load account data from the storage.
func (s *AccountsStorage) LoadAccount(email string) (account *g.Account, err error) {
	account = &g.Account{Email: email}

	if s.AccountFileExists() {
		b, _ := os.ReadFile(s.accountFilePath)
		_ = jsoniter.Unmarshal(b, account)
	}

	if s.PrivateKeyFileExists() {
		account.Key, err = s.LoadPrivateKeyFile()
		if err != nil {
			zap.S().Errorf("failed to load private key file: %v", err)
			return
		}
	}

	if account.Key == nil {
		account.Key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			zap.S().Error("failed to generate private key: %v", err)
			return
		}
	}

	return
}
