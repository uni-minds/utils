/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: hash_test.go
 */

package tools

import "testing"

func TestGetStringMD5(t *testing.T) {
	str := "B317-2021123109292773-2000010112000000-AA.mp4"
	t.Log(GetStringChecksum(str, ModeChecksumMD5))
	t.Log(GetStringChecksum(str, ModeChecksumSHA256))
}
