/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: imageProcessor.go
 */

package media

import (
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"os"
	"sync"
)

var threadLock chan int

func init() {
	threadLock = make(chan int, MaxUseThread)
}

type FrameProcessor struct {
	img       *image.RGBA
	data      []byte
	width     int
	height    int
	uvStart   int
	ch        chan int
	converted bool
	imgType   string
	//idxMap1 []int
	idxMap []int
}

func (f *FrameProcessor) SetDataNV12(w, h int, data []byte) {
	if f.width != w || f.height != h {
		f.width = w
		f.height = h
		f.uvStart = f.width * f.height
		f.GenerateIndexMapNV12(w, h)
	}

	f.data = data
	f.converted = false
	f.imgType = "nv12"
}

func (f *FrameProcessor) SetDataYUYV(w, h int, data []byte) {
	f.width = w
	f.height = h
	f.data = data
	f.converted = false
	f.imgType = "yuyv"
}

func (f *FrameProcessor) GenerateIndexMapNV12(width, height int) {
	index := 0
	idxMap := make([]int, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			iU := f.uvStart + x/2*2 + y/2*width
			idxMap[index] = iU
			index++
		}
	}
	f.idxMap = idxMap
	//fmt.Println("map data generated:", width, height, index)
}

func (f *FrameProcessor) GetImage() *image.RGBA {
	return f.GetImageParallel()

	//if !f.converted {
	//	f.img = image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	//	idx := 0
	//
	//	for y := 0; y < f.height; y++ {
	//		for x := 0; x < f.width; x++ {
	//			iU := f.idxMap[idx]
	//			f.img.Set(x, y, color.YCbCr{Y: f.data[idx], Cb: f.data[iU], Cr: f.data[iU+1]})
	//			idx++
	//		}
	//	}
	//
	//	f.converted = true
	//}
	//return f.img
}

func (f *FrameProcessor) GetImageParallel() *image.RGBA {
	if !f.converted {
		var wg sync.WaitGroup
		core := 2
		step := f.height / core

		f.img = image.NewRGBA(image.Rect(0, 0, f.width, f.height))
		for y := 0; y < f.height; y += step {
			wg.Add(1)
			go func(line, count int) {
				f.CalcLines(line, count)
				wg.Done()
			}(y, step)
		}

		wg.Wait()
		f.converted = true
	}
	return f.img
}

func (f *FrameProcessor) CalcLines(lineStart, lineCount int) {
	index := lineStart * f.width
	ymax := lineStart + lineCount
	if ymax > f.height {
		ymax = f.height
	}

	for y := lineStart; y < ymax; y++ {
		for x := 0; x < f.width; x++ {
			iU := f.idxMap[index]
			f.img.Set(x, y, color.YCbCr{Y: f.data[index], Cb: f.data[iU], Cr: f.data[iU+1]})
			index++
		}
	}
}

func (f *FrameProcessor) SaveBMP(filename string) error {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer fp.Close()
	if err != nil {
		return err
	}
	img := f.GetImage()
	if err = bmp.Encode(fp, img); err != nil {
		return err
	} else {

		return nil
	}
}

func (f *FrameProcessor) GetYUYV() (w, h int, data []byte) {
	img := f.GetImage()

	w = img.Bounds().Size().X
	h = img.Bounds().Size().Y
	data = make([]byte, w*h*2)

	bitIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w-1; x += 2 {
			r, g, b, _ := img.At(x, y).RGBA()
			y1, u1, v1 := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			r, g, b, _ = img.At(x+1, y).RGBA()
			y2, u2, v2 := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			data[bitIndex] = y1
			data[bitIndex+1] = u1>>1 + u2>>1
			data[bitIndex+2] = y2
			data[bitIndex+3] = v1>>1 + v2>>1
			bitIndex += 4
		}
	}
	return w, h, data
}
