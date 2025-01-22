// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-22, by liasica

package qiniu

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListCerts(t *testing.T) {
	certs, err := New(os.Getenv("AK"), os.Getenv("SK")).ListCerts()
	require.NoError(t, err)
	t.Log(certs)
}

func TestFindCert(t *testing.T) {
	cert, err := New(os.Getenv("AK"), os.Getenv("SK")).FindCert("671868a2b32fad84195b90a5")
	require.NoError(t, err)
	t.Logf("%#v", cert)
}

func TestCreateCert(t *testing.T) {
	pri, _ := os.ReadFile("../../runtime/archive/demo.bijuzaixian.com/privkey.pem")
	ca, _ := os.ReadFile("../../runtime/archive/demo.bijuzaixian.com/fullchain.pem")
	cert, err := New(os.Getenv("AK"), os.Getenv("SK")).CreateCert(&CreateCertRequest{
		Name:       "TEST",
		CommonName: "TEST-COMMON",
		Pri:        string(pri),
		Ca:         string(ca),
	})
	require.NoError(t, err)
	t.Log(cert)
}
