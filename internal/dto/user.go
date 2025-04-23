package dto

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewCreateUserRequest(r *http.Request) (*CreateUserRequest, error) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (req *CreateUserRequest) Validate() []Error {
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

type FindUserRequest struct {
	Page    *int32 `json:"page"`
	PerPage *int32 `json:"per_page"`
}

func NewFindUserRequest(r *http.Request) (*FindUserRequest, error) {
	query := r.URL.Query()

	var req FindUserRequest

	if page := query.Get("page"); page != "" {
		pageNum := int32(0)
		_, err := fmt.Sscanf(page, "%d", &pageNum)
		if err != nil {
			return nil, err
		}
		req.Page = &pageNum
	}

	if perPage := query.Get("per_page"); perPage != "" {
		perPageNum := int32(0)
		_, err := fmt.Sscanf(perPage, "%d", &perPageNum)
		if err != nil {
			return nil, err
		}
		req.PerPage = &perPageNum
	}

	return &req, nil
}
