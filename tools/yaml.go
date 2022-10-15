/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: yaml.go
 */

package tools

import (
	yaml "gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func LoadYaml(path string, i interface{}) error {
	if f, err := os.Open(path); err != nil {
		return err
	} else {
		return yaml.NewDecoder(f).Decode(i)
	}
}

func SaveYaml(path string, i interface{}) error {
	os.MkdirAll(filepath.Dir(path), 0644)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	} else {
		yaml.NewEncoder(f).Encode(i)
		return f.Close()
	}
}
