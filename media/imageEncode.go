/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: imageEncode.go
 */

package media

import (
	"errors"
	"golang.org/x/image/bmp"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func Encode(w io.Writer, i image.Image, subsampleRatio PixelFormat) error {
	switch subsampleRatio {
	case NV12:
		bs := ImageToNV12Parallel(i, -1)
		_, err := w.Write(bs)
		return err

	case JPEG:
		return jpeg.Encode(w, i, nil)

	case PNG:
		return png.Encode(w, i)

	case BMP:
		return bmp.Encode(w, i)

	default:
		return errors.New("unsupport format")
	}
}
