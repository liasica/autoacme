// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package provider

import "fmt"

type OssProvider struct {
}

func (p OssProvider) Present(domain, token, keyAuth string) error {
	fmt.Println("Oss Present", domain, token, keyAuth)
	return nil
}

func (p OssProvider) CleanUp(domain, token, keyAuth string) error {
	fmt.Println("Oss CleanUp", domain, token, keyAuth)
	return nil
}
