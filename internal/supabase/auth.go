package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserCredentials struct {
	Email    string
	Password string
	Data     interface{}
}

type User struct {
	ID                 string                    `json:"id"`
	Aud                string                    `json:"aud"`
	Role               string                    `json:"role"`
	Email              string                    `json:"email"`
	InvitedAt          time.Time                 `json:"invited_at"`
	ConfirmedAt        time.Time                 `json:"confirmed_at"`
	ConfirmationSentAt time.Time                 `json:"confirmation_sent_at"`
	AppMetadata        struct{ provider string } `json:"app_metadata"`
	UserMetadata       map[string]interface{}    `json:"user_metadata"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
}

type AuthenticatedDetails struct {
	AccessToken          string `json:"access_token"`
	TokenType            string `json:"token_type"`
	ExpiresIn            int    `json:"expires_in"`
	RefreshToken         string `json:"refresh_token"`
	User                 User   `json:"user"`
	ProviderToken        string `json:"provider_token"`
	ProviderRefreshToken string `json:"provider_refresh_token"`
}

type Auth struct {
	client *Client
}

func CreateAuth(baseURL string, supabaseKey string, debug ...bool) *Auth {
	client := CreateClient(baseURL, supabaseKey, debug...)
	return &Auth{client: client}
}

func (a *Auth) newRequestWithContext(method string, uri string, data any) (*http.Request, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	reqURL := fmt.Sprintf("%s/%s/%s", a.client.BaseURL, AuthEndpoint, uri)

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (a *Auth) SignUp(credentials UserCredentials) (*User, *ErrorResponse, error) {
	req, err := a.newRequestWithContext(http.MethodPost, "signup", credentials)
	if err != nil {
		return nil, nil, err
	}

	var res = User{}
	var errRes = ErrorResponse{}

	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)
	if err != nil {
		return nil, nil, err
	}

	if hasCustomError {
		return nil, &errRes, nil
	}

	return &res, nil, nil
}

func (a *Auth) SignIn(credentials UserCredentials) (*AuthenticatedDetails, *ErrorResponse, error) {
	req, err := a.newRequestWithContext(http.MethodPost, "token?grant_type=password", credentials)
	if err != nil {
		return nil, nil, err
	}

	res := AuthenticatedDetails{}
	errRes := ErrorResponse{}
	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)
	if err != nil {
		return nil, nil, err
	}

	if hasCustomError {
		if errRes.ErrorCode == "invalid_credentials" {
			errRes.Code = http.StatusUnauthorized
		}

		return nil, &errRes, nil
	}

	return &res, nil, nil
}

func (a *Auth) SignOut(userToken string) (*ErrorResponse, error) {
	req, err := a.newRequestWithContext(http.MethodPost, "logout", nil)
	if err != nil {
		return nil, err
	}

	injectAuthorizationHeader(req, userToken)

	errRes := ErrorResponse{}
	hasCustomError, err := a.client.sendCustomRequest(req, nil, &errRes)

	if err != nil {
		return nil, err
	}

	if hasCustomError {
		return &errRes, nil

	}

	return nil, nil
}

func (a *Auth) ForgottenPassword(email string) (*ErrorResponse, error) {
	reqBody := map[string]string{"email": email}
	req, err := a.newRequestWithContext(http.MethodPost, "recover", reqBody)
	if err != nil {
		return nil, err
	}

	errRes := ErrorResponse{}
	hasCustomError, err := a.client.sendCustomRequest(req, nil, &errRes)

	if err != nil {
		return nil, err
	}

	if hasCustomError {
		return &errRes, nil
	}

	return nil, nil
}

func (a *Auth) ResetPassword(userToken string, password string) (*ErrorResponse, error) {
	reqBody := map[string]string{"password": password}
	req, err := a.newRequestWithContext(http.MethodPut, "user?type=recovery", reqBody)
	if err != nil {
		return nil, err
	}

	injectAuthorizationHeader(req, userToken)

	res := AuthenticatedDetails{}
	errRes := ErrorResponse{}
	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)

	if err != nil {
		return nil, err
	}

	if hasCustomError {
		return &errRes, nil
	}

	return nil, nil
}

func (a *Auth) RefreshToken(refreshToken string) (*AuthenticatedDetails, *ErrorResponse, error) {
	reqBody := map[string]string{"refresh_token": refreshToken}
	req, err := a.newRequestWithContext(http.MethodPost, "token?grant_type=refresh_token", reqBody)
	if err != nil {
		return nil, nil, err
	}

	res := AuthenticatedDetails{}
	errRes := ErrorResponse{}
	hasCustomError, err := a.client.sendCustomRequest(req, &res, &errRes)

	if err != nil {
		return nil, nil, err
	}

	if hasCustomError {
		return nil, &errRes, nil
	}

	return &res, nil, nil
}
