// Copyright (C) autoacme. 2025-present.
//
// Created at 2025-01-07, by liasica

package storage

import (
	"os"
)

func CreateNonExistingFolder(path string) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0o700)
	}
	return
}
