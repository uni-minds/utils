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
	"io"
	"os"
	"strings"
)

func FileGetMD5(path string) (checksum string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	h := md5.New()
	if _, err = io.Copy(h, file); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	}
}

func FileCheckMD5(fullpath string, checksum string) (ok bool) {
	target, err := FileGetMD5(fullpath)
	if err != nil {
		fmt.Println("e", "get md5:", fullpath)
		return false
	} else {
		fmt.Println("d", "target checksum:", target)
		return target == checksum
	}
}

func GetBytesChecksum(bs []byte, m string) string {
	switch strings.ToLower(m) {
	case "sha256":
		h := sha256.New()
		h.Write(bs)
		return fmt.Sprintf("%x", h.Sum(nil))

	default:
		h := md5.New()
		h.Write(bs)
		return fmt.Sprintf("%x", h.Sum(nil))

	}
}

func GetStringChecksum(str string, m string) string {
	return GetBytesChecksum([]byte(str), m)
}

func GetStringMD5(str string) string {
	return GetBytesChecksum([]byte(str), "md5")
}

func GetBytesMD5(bs []byte) string {
	return GetBytesChecksum(bs, "md5")
}
