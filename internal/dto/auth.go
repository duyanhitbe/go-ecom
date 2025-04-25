package dto

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLoginRequest(r *http.Request) (*LoginRequest, error) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (req *LoginRequest) Validate() []Error {
	var errs []Error
	if req.Username == "" {
		errs = append(errs, Error{
			Field:   "username",
			Message: "username is required",
		})
	}
	if req.Password == "" {
		errs = append(errs, Error{
			Field:   "password",
			Message: "password is required",
		})
	}
	return errs
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

func NewLoginResponse(token string, expiresIn float64) *LoginResponse {
	return &LoginResponse{
		AccessToken: token,
		ExpiresIn:   int32(expiresIn),
	}
}
