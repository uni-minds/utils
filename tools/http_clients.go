/*
 * Copyright (c) 2019-2020
 * Author: LIU Xiangyu
 * File: httpClient.go
 */

package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ConnectData struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}

// var clientPool ClientPool
var pool *ClientPool

func init() {
	pool = NewClientPool(100)
}

func HttpGet(url string) (data ConnectData, bs []byte, err error) {
	client := pool.GetClient()
	defer pool.Recycle(client)

	if resp, err := client.Get(url); err != nil {
		return data, nil, err
	} else if bs, err = io.ReadAll(resp.Body); err != nil {
		return data, nil, err
	} else if err = resp.Body.Close(); err != nil {
		return data, bs, err
	} else {
		json.Unmarshal(bs, &data)
		return data, bs, nil
	}
}

func HttpPost(url string, data interface{}, contentType string) (rdata ConnectData, bs []byte, err error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return rdata, nil, err
	}
	req.Header.Add("content-type", contentType)

	client := pool.GetClient()
	defer pool.Recycle(client)

	resp, err := client.Do(req)
	req.Body.Close()
	if err != nil {
		return rdata, nil, err
	}
	defer resp.Body.Close()

	if bs, err = io.ReadAll(resp.Body); err != nil {
		return rdata, bs, err
	} else {
		if err = json.Unmarshal(bs, &rdata); err != nil {
			fmt.Println("http post: not a standard response")
		}
		return rdata, bs, nil
	}
}
