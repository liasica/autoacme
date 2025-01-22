// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package g

import (
	"crypto"
	"crypto/ecdsa"

	"github.com/go-acme/lego/v4/registration"
)

type Account struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	Key          *ecdsa.PrivateKey      `json:"-"`
}

func (u *Account) GetEmail() string {
	return u.Email
}

func (u *Account) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *Account) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}
