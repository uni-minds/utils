/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: screen.go
 */

package tools

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"os"
)

func ScreenClear() {
	print("\033[H\033[2J")
}

func PauseExit() {
	fmt.Println("Press any to exit")
	keyboard.GetSingleKey()
	os.Exit(0)
}
