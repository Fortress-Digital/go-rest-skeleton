package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type SupabaseClientInterface interface {
	sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error)
	newRequestWithContext(method string, reqURL string, data any) (*http.Request, error)
}

type SupabaseClient struct {
	BaseURL    string
	apiKey     string
	HTTPClient HttpClientInterface
}

func CreateClient(baseURL string, supabaseKey string) *SupabaseClient {
	client := &SupabaseClient{
		BaseURL: baseURL,
		apiKey:  supabaseKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}

	return client
}

func injectAuthorizationHeader(req *http.Request, value string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", value))
}

func (c *SupabaseClient) sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error) {
	req.Header.Set("apikey", c.apiKey)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return true, err
	}

	defer res.Body.Close()

	statusOK := res.StatusCode >= http.StatusOK && res.StatusCode < 300
	if !statusOK {
		if err = json.NewDecoder(res.Body).Decode(&errorValue); err == nil {
			return true, nil
		}

		return false, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	} else if res.StatusCode != http.StatusNoContent {
		if err = json.NewDecoder(res.Body).Decode(&successValue); err != nil {
			return false, err
		}
	}

	return false, nil
}

func (c *SupabaseClient) newRequestWithContext(method string, uri string, data any) (*http.Request, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	reqURL := fmt.Sprintf("%s/%s", c.BaseURL, uri)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}
