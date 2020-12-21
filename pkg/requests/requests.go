package requests

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// DoGetRequest makes a GET Request to a specific URL
func DoGetRequest(url string) (*string, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(res.Body)

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}
	strResult := string(data)
	return &strResult, nil
}

// DoPostRequestWithBody makes a POST Request with body to a specific URL
func DoPostRequestWithBody (URL string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}