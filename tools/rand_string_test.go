/*
 * Copyright (c) 2022
 * Author: LIU Xiangyu
 * File: base_test.go
 * Date: 2022/09/08 11:39:08
 */

package tools

import "testing"

func TestExpandInterface(t *testing.T) {
	data := []interface{}{"%s", "a", "b", "c"}
	str := InterfaceExpand(data)
	t.Log(str)
}
