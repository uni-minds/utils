/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: formatYUV.go
 */

/*
 * Import from https://github.com/shethchintan7/yuv
 * File: formatYUV.go
 */

package media

import (
	"image"
	"image/color"
)

type YUV struct {
	Y, UV       []uint8
	YStride     int
	CStride     int
	PixelFormat PixelFormat
	Rect        image.Rectangle
}

func (p *YUV) ColorModel() color.Model {
	return color.YCbCrModel
}

func (p *YUV) Bounds() image.Rectangle {
	return p.Rect
}

func (p *YUV) At(x, y int) color.Color {
	return p.YUVAt(x, y)
}

func (p *YUV) YUVAt(x, y int) color.YCbCr {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.YCbCr{}
	}
	yi := p.YOffset(x, y)
	ci := p.COffset(x, y)
	return color.YCbCr{
		p.Y[yi],
		p.UV[ci],
		p.UV[ci+1],
	}
}

// YOffset returns the index of the first element of Y that corresponds to
// the pixel at (x, y).
func (p *YUV) YOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.YStride + (x - p.Rect.Min.X)
}

// COffset returns the index of the first element of Cb or Cr that corresponds
// to the pixel at (x, y).
func (p *YUV) COffset(x, y int) int {
	switch p.PixelFormat {
	case NV12:
		return 2 * ((y/2-p.Rect.Min.Y/2)*p.CStride + (x/2 - p.Rect.Min.X/2))
	}
	// Default to 4:4:4 subsampling.
	return (y-p.Rect.Min.Y)*p.CStride + (x - p.Rect.Min.X)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *YUV) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &YUV{
			PixelFormat: p.PixelFormat,
		}
	}
	yi := p.YOffset(r.Min.X, r.Min.Y)
	ci := p.COffset(r.Min.X, r.Min.Y)
	return &YUV{
		Y:           p.Y[yi:],
		UV:          p.UV[ci:],
		PixelFormat: p.PixelFormat,
		YStride:     p.YStride,
		CStride:     p.CStride,
		Rect:        r,
	}
}

func (p *YUV) Opaque() bool {
	return true
}

func yuvSize(r image.Rectangle, subsampleRatio PixelFormat) (w, h, cw, ch int) {
	w, h = r.Dx(), r.Dy()
	switch subsampleRatio {
	case NV12:
		cw = (r.Max.X+1)/2 - r.Min.X/2
		ch = (r.Max.Y+1)/2 - r.Min.Y/2
	default:
		// Default to 4:4:4 subsampling.
		cw = w
		ch = h
	}
	return
}

// NewYCbCr returns a new YCbCr image with the given bounds and subsample ratio.
func NewYUV(r image.Rectangle, subsampleRatio PixelFormat) *YUV {
	w, h, cw, ch := yuvSize(r, subsampleRatio)
	i0 := w*h + 0*cw*ch
	//i1 := w*h + 1*cw*ch
	i2 := w*h + 2*cw*ch
	b := make([]byte, i2)
	return &YUV{
		Y:  b[:i0:i0],
		UV: b[i0:i2:i2],
		//V:              b[i1:i2:i2],
		PixelFormat: subsampleRatio,
		YStride:     w,
		CStride:     cw,
		Rect:        r,
	}
}
