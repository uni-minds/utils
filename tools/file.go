/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: file.go
 */

package tools

import (
	"fmt"
	"os"
)

func FileMove(src, dst string) error {
	if sourceFileStat, err := os.Stat(src); err != nil {
		return err

	} else if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)

	}

	if err := os.Rename(src, dst); err != nil {
		if err = Copy(src, dst); err != nil {
			return err
		} else if err := os.Remove(src); err != nil {
			return err
		}
	}
	return nil
}
