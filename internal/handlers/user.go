package handlers

import (
	"github.com/duyanhitbe/go-ecom/internal/dto"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	"github.com/duyanhitbe/go-ecom/pkg/utils"
	"net/http"
)

type UserHandler struct {
	repository repositories.Querier
}

func NewUserHandler(repository repositories.Querier) *UserHandler {
	return &UserHandler{repository: repository}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	req, err := dto.NewCreateUserRequest(r)
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusBadRequest, err)
		dto.Write(w, rsp)
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		rsp := dto.NewErrorsResponse(http.StatusBadRequest, errs)
		dto.Write(w, rsp)
		return
	}

	user, err := h.repository.CreateUser(r.Context(), &repositories.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusInternalServerError, err)
		dto.Write(w, rsp)
		return
	}
	rsp := dto.NewCreatedResponse(user)
	dto.Write(w, rsp)
}

func (h *UserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	req, err := dto.NewFindUserRequest(r)
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusBadRequest, err)
		dto.Write(w, rsp)
		return
	}

	page, perPage, offset := utils.GetPaginationMeta(req.Page, req.PerPage)
	users, err := h.repository.FindUser(r.Context(), &repositories.FindUserParams{
		Offset: offset,
		Limit:  perPage,
	})
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusInternalServerError, err)
		dto.Write(w, rsp)
		return
	}
	count, err := h.repository.CountUser(r.Context())
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusInternalServerError, err)
		dto.Write(w, rsp)
		return
	}

	meta := dto.NewMeta(perPage, page, count)
	rsp := dto.NewPaginatedResponse[[]*repositories.User](users, meta)
	dto.Write(w, rsp)
}
