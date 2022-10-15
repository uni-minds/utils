/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: imageConvert.go
 */

package media

import (
	"image"
	"image/color"
	"time"
)

//func YUYV2Image(w int, h int, data []byte) image.Image {
//	img := image.NewRGBA(image.Rect(0, 0, w, h))
//	bitIndex := 0
//	for y := 0; y < h; y++ {
//		for x := 0; x < w; x += 2 {
//			U0 := data[bitIndex+1]
//			V0 := data[bitIndex+3]
//			img.Set(x, y, color.YCbCr{Y: data[bitIndex], Cb: U0, Cr: V0})
//			img.Set(x+1, y, color.YCbCr{Y: data[bitIndex+2], Cb: U0, Cr: V0})
//			bitIndex += 4
//		}
//	}
//	return img
//}

func ImageFromYUYV2Parallel(w int, h int, data []byte, threadMax int) image.Image {
	if threadMax < 1 {
		threadMax = DefaultUseThread
	}
	lineStep := h / threadMax
	bitStart := 0
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y += lineStep {
		threadLock <- 1
		go func(lineIndex, bitIndex int) {
			var Y0, Y1, U0, V0 uint8
			for y := lineIndex; y < lineIndex+lineStep && y < h; y++ {
				for x := 0; x < w; x += 2 {
					Y0 = data[bitIndex]
					Y1 = data[bitIndex+2]
					U0 = data[bitIndex+1]
					V0 = data[bitIndex+3]

					img.Set(x, y, color.YCbCr{Y: Y0, Cb: U0, Cr: V0})
					img.Set(x+1, y, color.YCbCr{Y: Y1, Cb: U0, Cr: V0})
					bitIndex += 4
				}
			}
			<-threadLock
		}(y, bitStart)
		bitStart += w * 2 * lineStep
	}

	for {
		time.Sleep(1)
		if len(threadLock) == 0 {
			return img
		}
	}
}

func ImageFromNV12Parallel(w int, h int, data []byte, threadMax int) image.Image {
	if threadMax < 1 {
		threadMax = DefaultUseThread
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	lineStep := h / threadMax

	for y := 0; y < h; y += lineStep {
		//uvStart := w * h
		threadLock <- 1
		go func(lineIndex int) {
			yIndex := lineIndex * w
			uIndex := (h + lineIndex>>1) * w
			yMax := lineIndex + lineStep
			if yMax > h {
				yMax = h
			}
			var Y0, Y1, Y2, Y3, U, V uint8
			for y := lineIndex; y < yMax; y += 2 {
				// 2 line for each loop
				yIndex += w
				for x := 0; x < w; x += 2 {
					Y0 = data[yIndex]
					Y1 = data[yIndex+1]
					Y2 = data[yIndex+w]
					Y3 = data[yIndex+w+1]
					U = data[uIndex]
					V = data[uIndex+1]

					img.Set(x, y, color.YCbCr{Y: Y0, Cb: U, Cr: V})
					img.Set(x+1, y, color.YCbCr{Y: Y1, Cb: U, Cr: V})
					img.Set(x, y+1, color.YCbCr{Y: Y2, Cb: U, Cr: V})
					img.Set(x+1, y+1, color.YCbCr{Y: Y3, Cb: U, Cr: V})

					uIndex += 2
					yIndex += 2
				}
			}
			<-threadLock
		}(y)
	}

	for {
		time.Sleep(1)
		if len(threadLock) == 0 {
			return img
		}
	}
}

func ImageToYUYVParallel(img image.Image, threadMax int) []byte {
	if threadMax < 1 {
		threadMax = DefaultUseThread
	}

	bitStart := 0

	w := img.Bounds().Size().X
	h := img.Bounds().Size().Y
	data := make([]byte, w*h*2)

	lineStep := h / threadMax

	for y := 0; y < h; y += lineStep {
		threadLock <- 1
		go func(lineIndex, bitIndex int) {
			var r, g, b uint32
			var y1, y2, u1, u2, v1, v2 uint8

			for y := lineIndex; y < lineIndex+lineStep && y < h; y++ {
				for x := 0; x < w; x += 2 {
					r, g, b, _ = img.At(x, y).RGBA()
					y1, u1, v1 = color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

					r, g, b, _ = img.At(x+1, y).RGBA()
					y2, u2, v2 = color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

					data[bitIndex] = y1
					data[bitIndex+1] = u1>>1 + u2>>1
					data[bitIndex+2] = y2
					data[bitIndex+3] = v1>>1 + v2>>1
					bitIndex += 4
				}
			}
			<-threadLock
		}(y, bitStart)
		bitStart += w * 2 * lineStep
	}

	for {
		time.Sleep(1)
		if len(threadLock) == 0 {
			return data
		}
	}
}

func ImageToNV12Parallel(img image.Image, threadMax int) []byte {
	if threadMax < 1 {
		threadMax = DefaultUseThread
	}

	w := img.Bounds().Size().X
	h := img.Bounds().Size().Y
	ustart := w * h
	uline := w / 2
	lineStep := h / threadMax

	data := make([]byte, ustart*3/2)

	for line := 0; line < h; line += lineStep {
		threadLock <- 1
		go func(lineIndex int) {
			yMax := lineIndex + lineStep
			if yMax > h {
				yMax = h
			}

			cp := make([]color.Color, 4)
			for y := lineIndex; y < yMax; y += 2 {
				for x := 0; x < w; x += 2 {
					var cus, cvs uint32
					cp[0] = img.At(x, y)
					cp[1] = img.At(x+1, y)
					cp[2] = img.At(x, y+1)
					cp[3] = img.At(x+1, y+1)
					upleft := y*w + x

					iy := []int{upleft, upleft + 1, upleft + w, upleft + w + 1}
					//fmt.Println(idx)

					iu := ustart + 2*(y/2*uline+x/2)

					for i := 0; i < 4; i++ {
						r, g, b, _ := cp[i].RGBA()
						cy, cu, cv := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))
						data[iy[i]] = cy
						cus += uint32(cu)
						cvs += uint32(cv)
					}

					data[iu] = uint8(cus / 4)
					data[iu+1] = uint8(cvs / 4)
				}
			}

			<-threadLock
		}(line)
	}

	for {
		time.Sleep(1)
		if len(threadLock) == 0 {
			return data
		}
	}
}
