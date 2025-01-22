// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-22, by liasica

package qiniu

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomain(t *testing.T) {
	resp, err := New(os.Getenv("AK"), os.Getenv("SK")).client.R().Get(fmt.Sprintf("/domain/%s", "image.bijukeji.com"))
	require.NoError(t, err)
	t.Log(resp)
}

func TestUpdateDomainHttps(t *testing.T) {
	err := New(os.Getenv("AK"), os.Getenv("SK")).UpdateDomainHttps("image.bijukeji.com", &UpdateDomainHttpsRequest{
		CertID:      "671868a2b32fad84195b90a5",
		ForceHttps:  true,
		Http2Enable: true,
	})
	t.Log(err)
}
