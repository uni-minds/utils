/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: base.go
 */

package tools

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))
var mu sync.Mutex

const ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func InterfaceExpand(msg []interface{}) string {
	builder := strings.Builder{}
	for argNum, m := range msg {
		if m != "" {
			if argNum > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprint(m))
		}
	}
	return builder.String()
}

func RandString0f(len int) string {
	return RandStringFromAlphabet(len, ALPHABET[:16])
}

func RandStringFromAlphabet(length int, alphabet string) string {
	alphaLen := 62
	if alphabet == "" {
		alphabet = ALPHABET
	} else {
		alphaLen = len(alphabet)
	}

	bs := make([]byte, length)
	for i := 0; i < length; i++ {
		mu.Lock()
		b := r.Intn(alphaLen)
		mu.Unlock()
		bs[i] = alphabet[b]
	}
	return string(bs)
}
