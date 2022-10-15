package tools

import (
	"testing"
)

func TestHttpGet(t *testing.T) {
	_, bs, err := HttpGet("http://www.baidu.com")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(string(bs))
	}
}
