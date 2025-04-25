package handlers

import (
	"github.com/duyanhitbe/go-ecom/internal/dto"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	"github.com/duyanhitbe/go-ecom/pkg/hash"
	"github.com/duyanhitbe/go-ecom/pkg/utils"
	"net/http"
)

type UserHandler struct {
	repository repositories.Querier
	hash       hash.Hash
}

func NewUserHandler(repository repositories.Querier, hash hash.Hash) *UserHandler {
	return &UserHandler{
		repository: repository,
		hash:       hash,
	}
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

	hashPassword, err := h.hash.Hash(req.Password)
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusInternalServerError, err)
		dto.Write(w, rsp)
		return
	}

	user, err := h.repository.CreateUser(r.Context(), &repositories.CreateUserParams{
		Username: req.Username,
		Password: hashPassword,
	})
	if err != nil {
		rsp := dto.NewErrResponse(http.StatusInternalServerError, err)
		dto.Write(w, rsp)
		return
	}
	usr := dto.NewCreateUserResponse(user)
	rsp := dto.NewCreatedResponse[*dto.CreateUserResponse](usr)
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
	usr := dto.NewFindUserResponse(users)
	rsp := dto.NewPaginatedResponse[[]*dto.FindUserResponse](usr, meta)
	dto.Write(w, rsp)
}
