/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: hash.go
 */

package tools

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"syscall"
)

type ModeChecksum uint

const (
	ModeChecksumMD5    ModeChecksum = iota
	ModeChecksumSHA256 ModeChecksum = iota
)

func GetChecksum(i interface{}, mode ModeChecksum) string {
	switch i.(type) {
	case string:
		return GetStringChecksum(i.(string), mode)
	case []byte:
		return GetBytesChecksum(i.([]byte), mode)
	default:
		return ""
	}
}

func GetStringChecksum(str string, mode ModeChecksum) string {
	return GetBytesChecksum([]byte(str), mode)
}

func GetBytesChecksum(bs []byte, m ModeChecksum) string {
	switch m {
	case ModeChecksumSHA256:
		h := sha256.New()
		h.Write(bs)
		return fmt.Sprintf("%x", h.Sum(nil))

	case ModeChecksumMD5:
		h := md5.New()
		h.Write(bs)
		return fmt.Sprintf("%x", h.Sum(nil))

	default:
		return ""
	}
}

func GetFileChecksum(path string, m ModeChecksum) (string, error) {
	file, err := os.OpenFile(path, syscall.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var h hash.Hash
	switch m {
	case ModeChecksumMD5:
		h = md5.New()
	case ModeChecksumSHA256:
		h = sha256.New()
	}

	if _, err = io.Copy(h, file); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	}
}
