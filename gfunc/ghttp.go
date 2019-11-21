package gfunc

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPPost post请求
func HTTPPost(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	rspbody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return rspbody, nil
}

// HTTPGet get请求
func HTTPGet(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	rspbody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return rspbody, nil
}
