/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: imageConvert_test.go
 */

package media

import (
	"golang.org/x/image/bmp"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestImageFromNV12Parallel(t *testing.T) {
	img := ImageGen(1920, 1080, DefaultCopyright)

	fp, _ := os.Create("tmp.bmp")
	bmp.Encode(fp, img)
	fp.Close()
}

func TestImageToNV12Parallel(t *testing.T) {
	imageSrc := "/Users/liuxy/Desktop/1.png"
	fp, _ := os.Open(imageSrc)
	i, _ := png.Decode(fp)
	defer fp.Close()

	bs := ImageToNV12Parallel(i, 4)
	ioutil.WriteFile("out_3584x2240.nv12", bs, 0644)
}
