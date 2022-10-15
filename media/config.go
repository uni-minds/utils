/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: config.go
 */

package media

const MaxUseThread = 20
const DefaultUseThread = 8
const DefaultCopyright = "©2022 Beihang University & Uni-Ledger Co. Ltd."

type PixelFormat int

const (
	AUTO    PixelFormat = iota
	NV12    PixelFormat = iota
	YUYV422 PixelFormat = iota
	MP4V    PixelFormat = iota
	IMGS    PixelFormat = iota
	JPEG    PixelFormat = iota
	BMP     PixelFormat = iota
	PNG     PixelFormat = iota
)
