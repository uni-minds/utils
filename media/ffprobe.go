/*
 * Copyright (c) 2019-2022
 * Author: LIU Xiangyu
 * File: ffprobe.go
 */

package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type FfprobeResult struct {
	Streams []struct {
		Index              int    `json:"index"`
		CodecName          string `json:"codec_name"`
		CodecLongName      string `json:"codec_long_name"`
		Profile            string `json:"profile"`
		CodecType          string `json:"codec_type"`
		CodecTagString     string `json:"codec_tag_string"`
		CodecTag           string `json:"codec_tag"`
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		CodedWidth         int    `json:"coded_width"`
		CodedHeight        int    `json:"coded_height"`
		ClosedCaptions     int    `json:"closed_captions"`
		FilmGrain          int    `json:"film_grain"`
		HasBFrames         int    `json:"has_b_frames"`
		SampleAspectRatio  string `json:"sample_aspect_ratio"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
		PixFmt             string `json:"pix_fmt"`
		Level              int    `json:"level"`
		ColorRange         string `json:"color_range"`
		ColorSpace         string `json:"color_space"`
		ColorTransfer      string `json:"color_transfer"`
		ColorPrimaries     string `json:"color_primaries"`
		ChromaLocation     string `json:"chroma_location"`
		FieldOrder         string `json:"field_order"`
		Refs               int    `json:"refs"`
		IsAvc              string `json:"is_avc"`
		NalLengthSize      string `json:"nal_length_size"`
		Id                 string `json:"id"`
		RFrameRate         string `json:"r_frame_rate"`
		AvgFrameRate       string `json:"avg_frame_rate"`
		TimeBase           string `json:"time_base"`
		StartPts           int    `json:"start_pts"`
		StartTime          string `json:"start_time"`
		DurationTs         int    `json:"duration_ts"`
		Duration           string `json:"duration"`
		BitRate            string `json:"bit_rate"`
		BitsPerRawSample   string `json:"bits_per_raw_sample"`
		NbFrames           string `json:"nb_frames"`
		ExtradataSize      int    `json:"extradata_size"`
		Disposition        struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
			Captions        int `json:"captions"`
			Descriptions    int `json:"descriptions"`
			Metadata        int `json:"metadata"`
			Dependent       int `json:"dependent"`
			StillImage      int `json:"still_image"`
		} `json:"disposition"`
		Tags struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
			VendorId    string `json:"vendor_id"`
		} `json:"tags"`
	} `json:"streams"`
	Format struct {
		Filename string `json:"filename"`
		Duration string `json:"duration"`
		Size     string `json:"size"`
		BitRate  string `json:"bit_rate"`
	} `json:"format"`
}

func GetDuration(fpath string) (dur time.Duration, err error) {
	rst, err := FfprobeMedia(fpath)
	if err != nil {
		return 0, err
	}

	dur, err = time.ParseDuration(rst.Format.Duration + "s")
	return dur, err
}

func GetFps(fpath string) (fps float64, err error) {
	rst, err := FfprobeMedia(fpath)
	if err != nil {
		return 0, err
	}

	if len(rst.Streams) == 0 {
		return 0, errors.New("no stream")
	}

	fpsStr := rst.Streams[0].RFrameRate
	d := strings.Split(fpsStr, "/")
	if len(d) > 1 {
		if f1, err := strconv.ParseFloat(d[0], 10); err != nil {
			return 0, err
		} else if f2, err := strconv.ParseFloat(d[1], 10); err != nil {
			return 0, err
		} else {
			return f1 / f2, nil
		}
	} else {
		return strconv.ParseFloat(d[0], 10)
	}
}

func FfprobeMedia(fpath string) (rst FfprobeResult, err error) {
	if _, err = os.Stat(fpath); err != nil {
		return rst, err
	}
	command := fmt.Sprintf("ffprobe -select_streams v -show_entries format=duration,size,bit_rate,filename -show_streams -v quiet -of json -i \"%s\"", fpath)
	c := exec.Command("sh", "-c", command)
	bs, err := c.Output()
	if err != nil {
		return rst, err
	}

	err = json.Unmarshal(bs, &rst)
	return rst, err
}
