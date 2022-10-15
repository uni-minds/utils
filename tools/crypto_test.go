/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: crypto_test.go
 */

package tools

import (
	"testing"
	"time"
)

func TestRsaGenKeyPair(t *testing.T) {
	RsaGenKeypair()
	//t.Log(gotPriKey)
	//t.Log(gotPubKey)
	tmp := GetStringChecksum("abcd", "sha256")
	t.Log(tmp)
	s, p, err := RsaSha256SignHash(tmp)
	t.Log(s, p, err)
	t.Log(RsaSha256VerifyHash(tmp, s, p))
}

func TestRsaSavePrivateKey(t *testing.T) {
	RsaSavePrivateKey(nil, "./prikey.dat")
	t.Log("Publickey:\n", RsaGetPublicKey())
}

func TestRsaSha256Verify(t *testing.T) {
	d := time.Now()
	RsaGenKeypair()
	s, p, _ := RsaSha256Sign(d)
	got := RsaSha256Verify(d, s, p)
	t.Log(got)
}
