/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: imageGenerator.go
 */

package media

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"time"
)

var font *truetype.Font
var FontFilePath = []string{"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf", "/Library/Fonts/Arial Unicode.ttf"}

func ImageGen(width, height int, str string) image.Image {
	i := image.NewAlpha(image.Rect(0, 0, width, height))
	size := 20

	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			i.Set(x, y, color.White)
		}
		for y := height - size; y < height; y++ {
			i.Set(x, y, color.White)

		}
	}

	for x := width - size; x < width; x++ {
		for y := 0; y < size; y++ {
			i.Set(x, y, color.White)
		}
		for y := height - size; y < height; y++ {
			i.Set(x, y, color.White)
		}
	}

	if font == nil {
		for _, fpFont := range FontFilePath {
			if _, err := os.Stat(fpFont); err != nil {
				continue
			}
			fByte, _ := ioutil.ReadFile(fpFont)
			font, _ = freetype.ParseFont(fByte)
			break
		}
	}

	fontContext := freetype.NewContext()
	fontContext.SetDPI(72)
	fontContext.SetFont(font)
	fontContext.SetFontSize(40)
	fontContext.SetClip(i.Bounds())
	fontContext.SetDst(i)
	fontContext.SetSrc(image.White)

	//pt := freetype.Pt(width>>5, height>>1+int(fontContext.PointToFixed(40)>>8))

	pt := freetype.Pt(size+10, 50)
	_, _ = fontContext.DrawString(str, pt)
	pt = freetype.Pt(width-300, 50)
	_, _ = fontContext.DrawString("Emulator Mode", pt)
	pt = freetype.Pt(size+10, height-10)
	_, _ = fontContext.DrawString(time.Now().String(), pt)

	return i
}
