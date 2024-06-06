package apiProm

import (
	"bm.parts.server/config"
	"fmt"
	"net/http"
	"net/url"
)

func Get(path string, params map[string]string) (*http.Response, bool) {
	paramsObj := url.Values{}
	for name, value := range params {
		paramsObj.Add(name, value)
	}

	fmt.Printf("Len paramsObj %#v\n", len(paramsObj))

	apiUrl := config.PromUrl + path
	if len(paramsObj) > 0 {
		apiUrl += "?" + paramsObj.Encode()
	}

	fmt.Printf("%#v\n", apiUrl)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, false
	}

	req.Header.Add("Authorization", "Bearer "+config.PromToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, false
	}

	fmt.Printf("Status %#v\n", resp.Status)
	fmt.Printf("Status Code %#v\n", resp.StatusCode)

	return resp, true
}
