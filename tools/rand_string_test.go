/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: base_test.go
 * Date: 2022/09/08 11:39:08
 */

package tools

import (
	"testing"
	"time"
)

func TestExpandInterface(t *testing.T) {
	data := []interface{}{"%s", "a", "b", "c"}
	str := InterfaceExpand(data)
	t.Log(str)
}

func TestRandString0f(t *testing.T) {
	for g := 0; g < 50; g++ {
		go func(w int) {
			for i := 0; i < 100; i++ {
				t.Log(w, RandString0f(4+w*2))
				//time.Sleep(10*time.Millisecond)
			}
		}(g)
	}

	time.Sleep(3 * time.Second)
	t.Log("Stop")
}
