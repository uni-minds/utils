/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: crypto.go
 */

package tools

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	priKeyData *rsa.PrivateKey
	pubKeyData string
)

func RsaGenKeypair() (priKey string, pubKey string) {
	priKeyData, _ = rsa.GenerateKey(rand.Reader, 2048)
	pk, _ := x509.MarshalPKIXPublicKey(priKeyData.Public())
	pubKeyData = base64.StdEncoding.EncodeToString(pk)

	derStream := x509.MarshalPKCS1PrivateKey(priKeyData)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	bs := pem.EncodeToMemory(block)
	priKey = string(bs)
	return priKey, pubKeyData
}

func RsaGetPublicKey() string {
	return pubKeyData
}

func RsaSavePrivateKey(priKey *rsa.PrivateKey, filepath string) error {
	if priKey == nil {
		fmt.Println("no private key, generated.")
		RsaGenKeypair()
		priKey = priKeyData
	}
	derStream := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	if file, err := os.Create(filepath); err != nil {
		return err
	} else if err = pem.Encode(file, block); err != nil {
		return err
	}
	return nil
}

func RsaLoadPrivateKey(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(bs)
	priKeyData, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	} else {
		pk, _ := x509.MarshalPKIXPublicKey(priKeyData.Public())
		pubKeyData = base64.StdEncoding.EncodeToString(pk)
		return nil
	}
}

func RsaSha256SignHash(hash string) (signature, pubKey string, err error) {
	if priKeyData == nil {
		return "", "", errors.New("no private key loaded")
	}
	h, err := hex.DecodeString(hash)
	if err != nil {
		return "", "", err
	}

	s, err := rsa.SignPKCS1v15(rand.Reader, priKeyData, crypto.SHA256, h)
	if err != nil {
		return "", pubKey, err
	} else {
		signature = base64.StdEncoding.EncodeToString(s)
		return signature, pubKeyData, err
	}
}

func RsaSha256Sign(i interface{}) (signature, pubKey string, err error) {
	bs, _ := json.Marshal(i)
	h := GetBytesChecksum(bs, "sha256")
	return RsaSha256SignHash(h)
}

func RsaSha256VerifyHash(hash, signature, pubKey string) bool {
	h, _ := hex.DecodeString(hash)
	s, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	pkbs, _ := base64.StdEncoding.DecodeString(pubKey)
	pk, err := x509.ParsePKIXPublicKey(pkbs)
	if err != nil {
		return false
	}

	publicKey := pk.(*rsa.PublicKey)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, h, s)
	return err == nil
}

func RsaSha256Verify(i interface{}, signature, pubKey string) bool {
	bs, _ := json.Marshal(i)
	h := GetBytesChecksum(bs, "sha256")
	return RsaSha256VerifyHash(h, signature, pubKey)
}
