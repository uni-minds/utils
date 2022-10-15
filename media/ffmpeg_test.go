/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: ffmpeg_test.go
 */

package media

import "testing"

func TestFfmpegCommandGenerator(t *testing.T) {
	srcinfo := MediaInfo{
		//Filepath:  "/Users/liuxy/Downloads/proj/B317-2022022814405265-2000010112000000-AA.mp4",
		Filepath:  "test/comb.nv12",
		MediaType: NV12,
		FPS:       5,
		Width:     800,
		Height:    600,
	}

	descInfo := MediaInfo{
		Filepath:  "test/2.mp4",
		MediaType: MP4V,
		Width:     300,
		Height:    240,
	}

	err := ConvertFormat(srcinfo, descInfo)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log("OK")
	}
}
