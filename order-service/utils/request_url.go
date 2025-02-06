package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func RequestGET(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		req.Header.Set(key, values)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("RequestGET URL: %s\n", url)
	log.Printf("Response Status: %d\n", res.StatusCode)
	log.Printf("Response Body: %s\n", string(body))

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: %s", body)
		return nil, fmt.Errorf("failed to get data, status code: %d", res.StatusCode)
	}

	return body, nil
}

func RequestPOST(url string, headers map[string]string, dataPayload io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, dataPayload)
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		req.Header.Set(key, values)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Error: %s", body)
		return nil, fmt.Errorf("failed to get data, status code: %d", res.StatusCode)
	}

	return body, nil
}
