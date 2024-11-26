package supabase

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/url"
	"testing"
)

const SignInUrl = "http://localhost/auth/v1/token?grant_type=password"

type SupabaseClientMock struct {
	mock.Mock
}

func (m *SupabaseClientMock) sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error) {
	args := m.Called(req, successValue, errorValue)

	err, _ := errorValue.(*ErrorResponse)
	err.Code = 400
	err.ErrorCode = "error message"

	if req.URL.String() == SignInUrl {
		p := UserCredentials{}
		json.NewDecoder(req.Body).Decode(&p)

		if p.Email == "invalid_credentials@example.com" {
			err.ErrorCode = "invalid_credentials"
		}
	}

	return args.Bool(0), args.Error(1)
}

func (m *SupabaseClientMock) newRequestWithContext(method string, reqURL string, data interface{}) (*http.Request, error) {
	args := m.Called(method, reqURL, data)
	return args.Get(0).(*http.Request), args.Error(1)
}

func TestNewAuthClient(t *testing.T) {
	authClient := NewAuthClient("http://localhost", "123")

	if authClient == nil {
		t.Errorf("Expected client to be created")
	}

	_, ok := authClient.(*AuthClient)

	assert.Equal(t, ok, true)
}

func TestNewAuthRequestWithContext(t *testing.T) {
	mockClient := new(SupabaseClientMock)
	authClient := &AuthClient{client: mockClient}

	mockClient.
		On("newRequestWithContext", http.MethodGet, "auth/v1/test", nil).
		Return(&http.Request{}, nil)

	req, err := authClient.newAuthRequestWithContext(http.MethodGet, "test", nil)

	assert.Equal(t, req, &http.Request{})
	assert.Equal(t, err, nil)
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name                     string
		newRequestWithContextErr error
		sendCustomRequestRes     bool
		sendCustomRequestErr     error
		expectedUser             any
		expectedSystemErr        any
		expectedErr              error
	}{
		{
			name:                     "Should return user",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedUser:             &User{},
			expectedSystemErr:        nil,
			expectedErr:              nil,
		},
		{
			name:                     "New request with context should return error",
			newRequestWithContextErr: errors.New("new request error"),
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedUser:             nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("new request error"),
		},
		{
			name:                     "Send custom request should return error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     errors.New("send custom request error"),
			expectedUser:             nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("send custom request error"),
		},
		{
			name:                     "Send custom request should return service system error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedUser:             nil,
			expectedSystemErr: &ErrorResponse{
				Code:      400,
				ErrorCode: "error message",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqUrl, _ := url.Parse("http://localhost")
			req := &http.Request{
				Header: map[string][]string{},
				URL:    reqUrl,
			}
			mockClient := new(SupabaseClientMock)
			authClient := &AuthClient{client: mockClient}
			uc := UserCredentials{
				Email:    "test@example.com",
				Password: "password",
			}
			mockClient.
				On("newRequestWithContext", http.MethodPost, "auth/v1/signup", uc).
				Return(req, tt.newRequestWithContextErr)

			mockClient.
				On("sendCustomRequest", req, &User{}, &ErrorResponse{}).
				Return(tt.sendCustomRequestRes, tt.sendCustomRequestErr)

			user, systemErr, err := authClient.SignUp(uc)

			assert.Equal(t, user, tt.expectedUser)
			assert.Equal(t, systemErr, tt.expectedSystemErr)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name                     string
		email                    string
		newRequestWithContextErr error
		sendCustomRequestRes     bool
		sendCustomRequestErr     error
		expectedAuthenticated    any
		expectedSystemErr        any
		expectedErr              error
	}{
		{
			name:                     "Should return authenticated details",
			email:                    "test@example.com",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    &AuthenticatedDetails{},
			expectedSystemErr:        nil,
			expectedErr:              nil,
		},
		{
			name:                     "New request with context should return error",
			email:                    "test@example.com",
			newRequestWithContextErr: errors.New("new request error"),
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("new request error"),
		},
		{
			name:                     "Send custom request should return error",
			email:                    "test@example.com",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     errors.New("send custom request error"),
			expectedAuthenticated:    nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("send custom request error"),
		},
		{
			name:                     "Send custom request should return service system error",
			email:                    "test@example.com",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    nil,
			expectedSystemErr: &ErrorResponse{
				Code:      400,
				ErrorCode: "error message",
			},
			expectedErr: nil,
		},
		{
			name:                     "Send custom request should return invalid_credentials error",
			email:                    "invalid_credentials@example.com",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    nil,
			expectedSystemErr: &ErrorResponse{
				Code:      401,
				ErrorCode: "invalid_credentials",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(SupabaseClientMock)
			authClient := &AuthClient{client: mockClient}
			uc := UserCredentials{
				Email:    tt.email,
				Password: "password",
			}
			reqUrl, _ := url.Parse(SignInUrl)
			reqBody, _ := json.Marshal(uc)
			req := &http.Request{
				URL:    reqUrl,
				Header: map[string][]string{},
				Body:   io.NopCloser(bytes.NewBuffer(reqBody)),
			}

			mockClient.
				On("newRequestWithContext", http.MethodPost, "auth/v1/token?grant_type=password", uc).
				Return(req, tt.newRequestWithContextErr)

			mockClient.
				On("sendCustomRequest", req, &AuthenticatedDetails{}, &ErrorResponse{}).
				Return(tt.sendCustomRequestRes, tt.sendCustomRequestErr)

			authenticated, systemErr, err := authClient.SignIn(uc)

			assert.Equal(t, authenticated, tt.expectedAuthenticated)
			assert.Equal(t, systemErr, tt.expectedSystemErr)
			assert.Equal(t, err, tt.expectedErr)

		})
	}
}

func TestSignOut(t *testing.T) {
	tests := []struct {
		name                     string
		newRequestWithContextErr error
		sendCustomRequestRes     bool
		sendCustomRequestErr     error
		expectedSystemErr        any
		expectedErr              error
	}{
		{
			name:                     "Should return nil error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedSystemErr:        nil,
			expectedErr:              nil,
		},
		{
			name:                     "New request with context should return error",
			newRequestWithContextErr: errors.New("new request error"),
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("new request error"),
		},
		{
			name:                     "Send custom request should return error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     errors.New("send custom request error"),
			expectedSystemErr:        nil,
			expectedErr:              errors.New("send custom request error"),
		},
		{
			name:                     "Send custom request should return service system error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedSystemErr: &ErrorResponse{
				Code:      400,
				ErrorCode: "error message",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		reqUrl, _ := url.Parse("http://localhost")
		req := &http.Request{
			Header: map[string][]string{},
			URL:    reqUrl,
		}
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(SupabaseClientMock)
			authClient := &AuthClient{client: mockClient}

			mockClient.
				On("newRequestWithContext", http.MethodPost, "auth/v1/logout", nil).
				Return(req, tt.newRequestWithContextErr)

			mockClient.
				On("sendCustomRequest", req, nil, &ErrorResponse{}).
				Return(tt.sendCustomRequestRes, tt.sendCustomRequestErr)

			systemErr, err := authClient.SignOut("token")

			expectedHeader := ""
			if tt.newRequestWithContextErr == nil {
				expectedHeader = "Bearer token"
			}

			assert.Equal(t, expectedHeader, req.Header.Get("Authorization"))
			assert.Equal(t, systemErr, tt.expectedSystemErr)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestForgottenPassword(t *testing.T) {
	tests := []struct {
		name                     string
		newRequestWithContextErr error
		sendCustomRequestRes     bool
		sendCustomRequestErr     error
		expectedSystemErr        any
		expectedErr              error
	}{
		{
			name:                     "Should return nil error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedSystemErr:        nil,
			expectedErr:              nil,
		},
		{
			name:                     "New request with context should return error",
			newRequestWithContextErr: errors.New("new request error"),
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("new request error"),
		},
		{
			name:                     "Send custom request should return error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     errors.New("send custom request error"),
			expectedSystemErr:        nil,
			expectedErr:              errors.New("send custom request error"),
		},
		{
			name:                     "Send custom request should return service system error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedSystemErr: &ErrorResponse{
				Code:      400,
				ErrorCode: "error message",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqUrl, _ := url.Parse("http://localhost")
			req := &http.Request{
				Header: map[string][]string{},
				URL:    reqUrl,
			}

			mockClient := new(SupabaseClientMock)
			authClient := &AuthClient{client: mockClient}
			email := "test@example.com"
			contextBody := map[string]string{"email": email}
			mockClient.
				On("newRequestWithContext", http.MethodPost, "auth/v1/recover", contextBody).
				Return(req, tt.newRequestWithContextErr)

			mockClient.
				On("sendCustomRequest", req, nil, &ErrorResponse{}).
				Return(tt.sendCustomRequestRes, tt.sendCustomRequestErr)

			systemErr, err := authClient.ForgottenPassword(email)

			assert.Equal(t, systemErr, tt.expectedSystemErr)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestResetPassword(t *testing.T) {
	tests := []struct {
		name                     string
		newRequestWithContextErr error
		sendCustomRequestRes     bool
		sendCustomRequestErr     error
		expectedSystemErr        any
		expectedErr              error
	}{
		{
			name:                     "Should return nil error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedSystemErr:        nil,
			expectedErr:              nil,
		},
		{
			name:                     "New request with context should return error",
			newRequestWithContextErr: errors.New("new request error"),
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("new request error"),
		},
		{
			name:                     "Send custom request should return error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     errors.New("send custom request error"),
			expectedSystemErr:        nil,
			expectedErr:              errors.New("send custom request error"),
		},
		{
			name:                     "Send custom request should return service system error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedSystemErr: &ErrorResponse{
				Code:      400,
				ErrorCode: "error message",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqUrl, _ := url.Parse("http://localhost")
			req := &http.Request{
				Header: map[string][]string{},
				URL:    reqUrl,
			}
			mockClient := new(SupabaseClientMock)
			authClient := &AuthClient{client: mockClient}
			token := "token"
			password := "password"
			contextBody := map[string]string{"password": password}
			mockClient.
				On("newRequestWithContext", http.MethodPut, "auth/v1/user?type=recovery", contextBody).
				Return(req, tt.newRequestWithContextErr)

			mockClient.
				On("sendCustomRequest", req, nil, &ErrorResponse{}).
				Return(tt.sendCustomRequestRes, tt.sendCustomRequestErr)

			systemErr, err := authClient.ResetPassword(token, password)

			expectedHeader := ""
			if tt.newRequestWithContextErr == nil {
				expectedHeader = "Bearer token"
			}

			assert.Equal(t, expectedHeader, req.Header.Get("Authorization"))
			assert.Equal(t, systemErr, tt.expectedSystemErr)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name                     string
		newRequestWithContextErr error
		sendCustomRequestRes     bool
		sendCustomRequestErr     error
		expectedAuthenticated    any
		expectedSystemErr        any
		expectedErr              error
	}{
		{
			name:                     "Should return authenticated details",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    &AuthenticatedDetails{},
			expectedSystemErr:        nil,
			expectedErr:              nil,
		},
		{
			name:                     "New request with context should return error",
			newRequestWithContextErr: errors.New("new request error"),
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("new request error"),
		},
		{
			name:                     "Send custom request should return error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     false,
			sendCustomRequestErr:     errors.New("send custom request error"),
			expectedAuthenticated:    nil,
			expectedSystemErr:        nil,
			expectedErr:              errors.New("send custom request error"),
		},
		{
			name:                     "Send custom request should return service system error",
			newRequestWithContextErr: nil,
			sendCustomRequestRes:     true,
			sendCustomRequestErr:     nil,
			expectedAuthenticated:    nil,
			expectedSystemErr: &ErrorResponse{
				Code:      400,
				ErrorCode: "error message",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqUrl, _ := url.Parse("http://localhost")
			req := &http.Request{
				Header: map[string][]string{},
				URL:    reqUrl,
			}

			mockClient := new(SupabaseClientMock)
			authClient := &AuthClient{client: mockClient}
			refreshToken := "token"
			contextBody := map[string]string{"refresh_token": refreshToken}

			mockClient.
				On("newRequestWithContext", http.MethodPost, "auth/v1/token?grant_type=refresh_token", contextBody).
				Return(req, tt.newRequestWithContextErr)

			mockClient.
				On("sendCustomRequest", req, &AuthenticatedDetails{}, &ErrorResponse{}).
				Return(tt.sendCustomRequestRes, tt.sendCustomRequestErr)

			authenticated, systemErr, err := authClient.RefreshToken(refreshToken)

			assert.Equal(t, authenticated, tt.expectedAuthenticated)
			assert.Equal(t, systemErr, tt.expectedSystemErr)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}
