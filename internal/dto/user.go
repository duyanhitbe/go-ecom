package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	"github.com/google/uuid"
	"net/http"
	"time"
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

type CreateUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCreateUserResponse(user *repositories.User) *CreateUserResponse {
	return &CreateUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
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
		if pageNum < 1 {
			return nil, errors.New("page must be greater than 0")
		}
		req.Page = &pageNum
	}

	if perPage := query.Get("per_page"); perPage != "" {
		perPageNum := int32(0)
		_, err := fmt.Sscanf(perPage, "%d", &perPageNum)
		if err != nil {
			return nil, err
		}
		if perPageNum < 1 {
			return nil, errors.New("per_page must be greater than 0")
		}
		req.PerPage = &perPageNum
	}

	return &req, nil
}

type FindUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFindUserResponse(users []*repositories.User) []*FindUserResponse {
	var u []*FindUserResponse

	for _, user := range users {
		u = append(u, &FindUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return u
}
