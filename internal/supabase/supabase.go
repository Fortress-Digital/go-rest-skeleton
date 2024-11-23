package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	postgrest "github.com/nedpals/postgrest-go/pkg"
)

const (
	AuthEndpoint = "auth/v1"
	RestEndpoint = "rest/v1"
)

type ErrorResponse struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

type SupabaseClientInterface interface {
	sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error)
	newRequestWithContext(method string, uri string, data any) (*http.Request, error)
}

type SupabaseClient struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
	DB         *postgrest.Client
}

func CreateClient(baseURL string, supabaseKey string, debug ...bool) *SupabaseClient {
	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s/", baseURL, RestEndpoint))
	if err != nil {
		panic(err)
	}
	client := &SupabaseClient{
		BaseURL: baseURL,
		apiKey:  supabaseKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		DB: postgrest.NewClient(
			*parsedURL,
			postgrest.WithTokenAuth(supabaseKey),
			func(c *postgrest.Client) {
				// debug parameter is only for postgrest-go for now
				if len(debug) > 0 {
					c.Debug = debug[0]
				}
				c.AddHeader("apikey", supabaseKey)
			},
		),
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

		return false, fmt.Errorf("unknown, status code: %d", res.StatusCode)
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
	reqURL := fmt.Sprintf("%s/%s/%s", c.BaseURL, AuthEndpoint, uri)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
