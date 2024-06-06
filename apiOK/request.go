package apiOK

import (
	"bm.parts.server/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Get(url string) (*http.Response, bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}

	return resp, true
}

func Post(url string, params map[string]interface{}) *http.Response {
	body, _ := json.Marshal(params)

	client := &http.Client{
		Timeout: time.Minute * 10,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp
}

func PostFull(url string, params db.Full) *http.Response {
	body, _ := json.Marshal(params)

	client := &http.Client{
		Timeout: time.Minute * 10,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return resp
}
