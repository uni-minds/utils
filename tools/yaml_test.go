/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: yaml_test.go
 */

package tools

import "testing"

func TestSaveYaml(t *testing.T) {
	//data := []string{"a", "b", "c"}
	data := map[string][]string{"1": {"a", "b", "c"}, "2": {"d", "e", "f"}}
	err := SaveYaml("./test.yaml", data)
	t.Log(err)
}
