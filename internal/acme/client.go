// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-08, by liasica

package acme

import (
	"time"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"go.uber.org/zap"

	"github.com/liasica/aliacme/internal/acme/storage"
	"github.com/liasica/aliacme/internal/g"
)

func SetupClient() (client *lego.Client, err error) {
	cfg := g.GetConfig()

	// Get accounts storage
	var accountsStorage *storage.AccountsStorage
	accountsStorage, err = storage.NewAccountsStorage(cfg.Account)
	if err != nil {
		return
	}

	// Load or create account
	var user *g.Account
	user, err = accountsStorage.LoadAccount(cfg.Account)
	if err != nil {
		zap.L().Error("failed to load account", zap.Error(err))
		return
	}

	config := lego.NewConfig(user)

	// TODO: 10分钟超时
	config.Certificate.Timeout = 1 * time.Minute

	// A client facilitates communication with the CA server.
	client, err = lego.NewClient(config)
	if err != nil {
		zap.L().Error("failed to create lego client", zap.Error(err))
		return
	}

	// needSave := false
	// New users will need to register
	var reg *registration.Resource
	if user.Registration == nil {
		// needSave = true
		reg, err = client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
		if err != nil {
			zap.L().Error("failed to register", zap.Error(err))
			return
		}

		user.Registration = reg
		accountsStorage.Save(user)
		// } else {
		// 	reg, err = client.Registration.QueryRegistration()
		// 	if err != nil {
		// 		zap.L().Error("failed to query registration", zap.Error(err))
		// 		return
		// 	}
	}

	// // Save user
	// user.Registration = reg
	// if needSave {
	// 	accountsStorage.Save(user)
	// }

	return
}
