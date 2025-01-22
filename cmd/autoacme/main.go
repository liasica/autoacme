// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package main

import (
	"flag"

	"github.com/liasica/autoacme/internal"
)

var Version = "v1.0.0"

func main() {
	var cfg string
	var storage string
	flag.StringVar(&cfg, "config", "/etc/autoacme/config.yaml", "Config file")
	flag.StringVar(&storage, "storage", "/etc/autoacme", "Storage file")
	flag.Parse()

	internal.Boot(cfg, storage, Version)
	internal.New().Run()
}
