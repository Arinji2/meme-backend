package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ApiClient struct {
	BaseURL string
	Client  HTTPClient
}

func NewApiClient(baseURL ...string) *ApiClient {
	var url string
	if len(baseURL) > 0 {
		url = baseURL[0]
	} else {
		url = ""
	}
	return &ApiClient{
		BaseURL: url,
		Client:  &http.Client{},
	}
}

func (c *ApiClient) doRequest(req *http.Request, headers map[string]string) (map[string]interface{}, int, error) {
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	var result map[string]interface{}
	var status int
	var err error

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	result = make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("error decoding response: %w", err)
	}

	status = resp.StatusCode
	return result, status, nil
}

func (c *ApiClient) SendRequestWithBody(method, path string, body interface{}, headers map[string]string, contentType string) (result map[string]interface{}, status int, err error) {
	address := fmt.Sprintf("%s%s", c.BaseURL, path)

	var reqBody *bytes.Buffer

	if contentType == "application/x-www-form-urlencoded" {
		formData, ok := body.(map[string]string)
		if !ok {
			return nil, 500, fmt.Errorf("expected body to be map[string]string for form-encoded data")
		}

		formValues := url.Values{}
		for key, value := range formData {
			formValues.Set(key, value)
		}
		reqBody = bytes.NewBufferString(formValues.Encode())
	} else {
		jsonBody, jsonErr := json.Marshal(body)
		if jsonErr != nil {
			return nil, 500, fmt.Errorf("error marshalling json: %w", jsonErr)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		if contentType == "" {
			contentType = "application/json"
		}
	}

	req, reqErr := http.NewRequest(method, address, reqBody)
	if reqErr != nil {
		return nil, 500, fmt.Errorf("error creating request: %w", reqErr)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	req.Header.Set("Content-Type", contentType)

	result, status, err = c.doRequest(req, headers)
	if err != nil {
		return nil, status, fmt.Errorf("error from request doer: %w", err)
	}

	return result, status, nil
}

func (c *ApiClient) SendRequestWithQuery(method, path string, query map[string]string, headers map[string]string) (result map[string]interface{}, status int, err error) {
	queryParams := url.Values{}
	for key, value := range query {
		queryParams.Add(key, value)
	}

	address, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		status = 500
		err = fmt.Errorf("error joining URL paths: %w", err)
		return
	}

	fullURL := fmt.Sprintf("%s?%s", address, queryParams.Encode())
	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		status = 500
		err = fmt.Errorf("error creating request: %w", err)
		return
	}

	result, status, err = c.doRequest(req, headers)
	if err != nil {
		err = fmt.Errorf("error from request doer: %w", err)
		return
	}

	return
}
