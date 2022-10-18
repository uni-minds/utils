/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: ffmpeg.go
 */

package media

import (
	"errors"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type MediaInfo struct {
	Filepath  string
	MediaType PixelFormat
	FPS       int
	Width     int
	Height    int
	Bitrate   string
}

func ConvertFormat(srcInfo, descInfo MediaInfo) error {
	if srcInfo.MediaType == descInfo.MediaType && srcInfo.MediaType > 0 {
		fmt.Println("ignore convert to same format")
		return errors.New("ignore convert to same format")
	}

	kv1 := make(map[string]interface{})
	kv2 := make(map[string]interface{})
	switch srcInfo.MediaType {
	case NV12:
		if srcInfo.Width*srcInfo.Height <= 0 {
			return errors.New("src height or width not given")
		}
		if srcInfo.FPS > 0 {

			kv1["r"] = srcInfo.FPS
		} else if descInfo.FPS > 0 {
			kv1["r"] = descInfo.FPS
		} else {
			return errors.New("nv12 need fps info")
		}

		kv1["pix_fmt"] = "nv12"
		kv1["f"] = "rawvideo"
		kv1["s"] = fmt.Sprintf("%dx%d", srcInfo.Width, srcInfo.Height)
		kv1["an"] = ""
	default:

	}

	switch descInfo.MediaType {
	case NV12:
		kv2["pix_fmt"] = "nv12"
		kv2["vcodec"] = "rawvideo"
		kv2["f"] = "image2"

	case MP4V:
		kv2["vf"] = "pad=ceil(iw/2)*2:ceil(ih/2)*2"
		kv2["pix_fmt"] = "yuv420p"

	default:
	}
	if descInfo.Bitrate != "" {
		kv2["b:v"] = descInfo.Bitrate
	}

	if descInfo.Width*descInfo.Height > 0 {
		kv2["vf"] = fmt.Sprintf("scale=%d:%d", descInfo.Width, descInfo.Height)
	}

	return ffmpeg.Input(srcInfo.Filepath, kv1).Output(descInfo.Filepath, kv2).OverWriteOutput().Run()
}
