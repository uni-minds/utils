/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: ffprobe_test.go
 */

package media

import (
	"testing"
)

const TESTMEDIA = "/Users/liuxy/Downloads/proj/B317-2022022814405265-2000010112000000-AA.mp4"

func Test_ffprobeExec(t *testing.T) {
	//info, _ := FfprobeMedia()
	info, err := FfprobeMedia(TESTMEDIA)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(info)
	}
}

func TestGetMediaDuration(t *testing.T) {
	t.Log(GetDuration(TESTMEDIA))
}

func TestGetFps(t *testing.T) {
	t.Log(GetFps(TESTMEDIA))
}
