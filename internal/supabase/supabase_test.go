package supabase

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCreateClient(t *testing.T) {
	client := CreateClient("http://localhost", "123")
	if client == nil {
		t.Errorf("Expected client to be created")
	}

	assert.Equal(t, client.BaseURL, "http://localhost")
	assert.Equal(t, client.apiKey, "123")
	assert.Equal(t, client.HTTPClient, &http.Client{
		Timeout: time.Minute,
	})
}

func TestInjectAuthorizationHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost", nil)

	injectAuthorizationHeader(req, "123")

	assert.Equal(t, req.Header.Get("Authorization"), "Bearer 123")
}

func TestSendCustomRequest(t *testing.T) {
	mockClient := new(MockHttpClient)
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: mockClient,
	}

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"test": "foo"}`))

	res := w.Result()

	mockClient.On("Do", mock.Anything).Return(res, nil)

	req, _ := http.NewRequest("GET", "http://localhost", nil)

	var successValue interface{}
	var errorValue interface{}

	result, err := sut.sendCustomRequest(req, &successValue, &errorValue)

	assert.Equal(t, result, false)
	assert.Equal(t, err, nil)
	assert.Equal(t, successValue, map[string]interface{}{"test": "foo"})
	assert.Equal(t, errorValue, interface{}(nil))
}

func TestSendCustomRequestNoContentResponse(t *testing.T) {
	mockClient := new(MockHttpClient)
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: mockClient,
	}

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`{"test": "foo"}`))

	res := w.Result()

	mockClient.On("Do", mock.Anything).Return(res, nil)

	req, _ := http.NewRequest("GET", "http://localhost", nil)

	var successValue interface{}
	var errorValue interface{}

	result, err := sut.sendCustomRequest(req, &successValue, &errorValue)

	assert.Equal(t, result, false)
	assert.Equal(t, err, nil)
	assert.Equal(t, successValue, interface{}(nil))
	assert.Equal(t, errorValue, interface{}(nil))
}

func TestSendCustomeRequestHttpClientError(t *testing.T) {
	mockClient := new(MockHttpClient)
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: mockClient,
	}

	res := httptest.NewRecorder()
	clientError := errors.New("Some test error")

	mockClient.On("Do", mock.Anything).Return(res.Result(), clientError)

	req, _ := http.NewRequest("GET", "http://localhost", nil)

	var successValue interface{}
	var errorValue interface{}

	result, err := sut.sendCustomRequest(req, &successValue, &errorValue)

	assert.Equal(t, result, true)
	assert.Equal(t, err, clientError)
	assert.Equal(t, successValue, interface{}(nil))
	assert.Equal(t, errorValue, interface{}(nil))
}

func TestSendCustomRequestBadRequest(t *testing.T) {
	mockClient := new(MockHttpClient)
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: mockClient,
	}

	res := httptest.NewRecorder()
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`{"test": "foo"}`))

	mockClient.On("Do", mock.Anything).Return(res.Result(), nil)

	req, _ := http.NewRequest("GET", "http://localhost", nil)

	var successValue interface{}
	var errorValue interface{}

	result, err := sut.sendCustomRequest(req, &successValue, &errorValue)

	assert.Equal(t, result, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, errorValue, map[string]interface{}{"test": "foo"})
	assert.Equal(t, successValue, interface{}(nil))
}

func TestSendCustomRequestJsonDecodeError(t *testing.T) {
	mockClient := new(MockHttpClient)
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: mockClient,
	}

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`!`))

	res := w.Result()

	mockClient.On("Do", mock.Anything).Return(res, nil)

	req, _ := http.NewRequest("GET", "http://localhost", nil)

	var successValue interface{}
	var errorValue interface{}

	result, err := sut.sendCustomRequest(req, &successValue, &errorValue)

	assert.Equal(t, result, false)
	assert.Equal(t, err.Error(), "invalid character '!' looking for beginning of value")
	assert.Equal(t, successValue, interface{}(nil))
	assert.Equal(t, errorValue, interface{}(nil))
}

func TestSendCustomRequestBadRequestJsonDecodeError(t *testing.T) {
	mockClient := new(MockHttpClient)
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: mockClient,
	}

	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`!`))

	res := w.Result()

	mockClient.On("Do", mock.Anything).Return(res, nil)

	req, _ := http.NewRequest("GET", "http://localhost", nil)

	var successValue interface{}
	var errorValue interface{}

	result, err := sut.sendCustomRequest(req, &successValue, &errorValue)

	assert.Equal(t, result, false)
	assert.Equal(t, err.Error(), "unknown error, status code: 400")
	assert.Equal(t, successValue, interface{}(nil))
	assert.Equal(t, errorValue, interface{}(nil))
}

func TestNewRequestWithContext(t *testing.T) {
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: &http.Client{},
	}

	req, err := sut.newRequestWithContext(http.MethodGet, "auth/v1/test", nil)

	assert.Equal(t, err, nil)
	assert.Equal(t, req.Method, http.MethodGet)
	assert.Equal(t, req.URL.String(), "http://localhost/auth/v1/test")
	assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, req.Header.Get("Accept"), "application/json")
}

func TestNewRequestWithContextJsonMarshalError(t *testing.T) {
	sut := &SupabaseClient{
		BaseURL:    "http://localhost",
		apiKey:     "123",
		HTTPClient: &http.Client{},
	}

	req, err := sut.newRequestWithContext(http.MethodGet, "auth/v1/test", make(chan int))

	assert.Equal(t, err.Error(), "json: unsupported type: chan int")
	assert.Equal(t, req, (*http.Request)(nil))
}

func TestNewRequestWithContextHttpNewRequestError(t *testing.T) {
	sut := &SupabaseClient{
		BaseURL:    "mysql://example{123",
		apiKey:     "123",
		HTTPClient: &http.Client{},
	}

	req, err := sut.newRequestWithContext(http.MethodGet, "auth/v1/test", nil)

	assert.Equal(t, err.Error(), "parse \"mysql://example{123/auth/v1/test\": invalid character \"{\" in host name")
	assert.Equal(t, req, (*http.Request)(nil))
}
