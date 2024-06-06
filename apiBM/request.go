package apiBM

import (
	"bm.parts.server/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var RateRemaining = 10000
var RateResetTime = 0

func Get(path string, params map[string]string, lang string) (*http.Response, bool) {
	if RateRemaining < 10 && RateResetTime > int(time.Now().UTC().Unix()) {
		fmt.Println("Rate Remaining ", RateRemaining, ", Waiting for ", RateResetTime-int(time.Now().UTC().Unix()), "seconds")
		time.Sleep(time.Duration(RateResetTime-int(time.Now().UTC().Unix())) * time.Second)
	}

	uri := config.BmUrl + path

	if len(params) > 0 {
		paramsObj := url.Values{}
		for name, value := range params {
			paramsObj.Add(name, value)
		}

		if len(paramsObj) > 0 {
			uri += "?" + paramsObj.Encode()
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", config.BmUrl+path, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Add("Authorization", config.BmToken)
	req.Header.Add("User-Agent", config.BmUserAgent)
	req.Header.Add("Accept-Language", lang)

	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}

	RateRemaining, _ = strconv.Atoi(resp.Header["X-Ratelimit-Remaining"][0])
	RateResetTime, _ = strconv.Atoi(resp.Header["X-Ratelimit-Reset"][0])

	return resp, true
}

func Post(path string, params map[string]interface{}) *http.Response {
	if RateRemaining < 10 && RateResetTime > int(time.Now().UTC().Unix()) {
		fmt.Println("Rate Remaining ", RateRemaining, ", Waiting for ", RateResetTime-int(time.Now().UTC().Unix()), "seconds")
		time.Sleep(time.Duration(RateResetTime-int(time.Now().UTC().Unix())) * time.Second)
	}

	body, _ := json.Marshal(params)

	client := &http.Client{}
	req, err := http.NewRequest("POST", config.BmUrl+path, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Add("Authorization", config.BmToken)
	req.Header.Add("User-Agent", config.BmUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	RateRemaining, _ = strconv.Atoi(resp.Header["X-Ratelimit-Remaining"][0])
	RateResetTime, _ = strconv.Atoi(resp.Header["X-Ratelimit-Reset"][0])

	return resp
}
