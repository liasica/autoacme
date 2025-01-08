// Copyright (C) aliacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package main

import (
	"flag"

	"github.com/liasica/aliacme/internal"
)

var Version = "v1.0.0"

func main() {
	var cfg string
	var storage string
	flag.StringVar(&cfg, "config", "/etc/aliacme/config.yaml", "Config file")
	flag.StringVar(&storage, "storage", "/etc/aliacme", "Storage file")
	flag.Parse()

	internal.Boot(cfg, storage, Version)
	internal.New().Run()
}
