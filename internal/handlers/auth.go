package handlers

import (
	"context"
	"github.com/duyanhitbe/go-ecom/internal/dto"
	"github.com/duyanhitbe/go-ecom/internal/repositories"
	"github.com/duyanhitbe/go-ecom/pkg/constants"
	"github.com/duyanhitbe/go-ecom/pkg/hash"
	"github.com/duyanhitbe/go-ecom/pkg/token"
	"github.com/duyanhitbe/go-ecom/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AuthHandler struct {
	repository repositories.Querier
	token      token.Token
	hash       hash.Hash
}

func NewAuthHandler(repository repositories.Querier, token token.Token, hash hash.Hash) *AuthHandler {
	return &AuthHandler{
		repository: repository,
		token:      token,
		hash:       hash,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, err := dto.NewLoginRequest(r)
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

	user, errRsp := getExistUser(r.Context(), h.repository, req.Username)
	if errRsp != nil {
		dto.Write(w, errRsp)
		return
	}

	errRsp = verifyUserPassword(h.hash, user.Password, req.Password)
	if errRsp != nil {
		dto.Write(w, errRsp)
		return
	}

	expiresIn := constants.TokenDuration
	tk, errRsp := signToken(h.token, user.ID, expiresIn)
	if errRsp != nil {
		dto.Write(w, errRsp)
		return
	}

	rsp := dto.NewLoginResponse(tk, expiresIn.Seconds())
	result := dto.NewOKResponse[*dto.LoginResponse](rsp)
	dto.Write(w, result)
}

func getExistUser(ctx context.Context, repository repositories.Querier, username string) (*repositories.User, *dto.Response[*any]) {
	user, err := repository.FindOneUserByUsername(ctx, username)
	if err != nil {
		if utils.IsErrNoRows(err) {
			return nil, dto.NewErrWithFieldResponse(http.StatusBadRequest, "username", constants.ErrUserNotFound)

		}

		return nil, dto.NewErrResponse(http.StatusInternalServerError, err)
	}
	return user, nil
}

func verifyUserPassword(hash hash.Hash, userPassword string, reqPassword string) *dto.Response[*any] {
	ok, err := hash.Verify(userPassword, reqPassword)
	if err != nil {
		return dto.NewErrResponse(http.StatusInternalServerError, err)
	}
	if !ok {
		return dto.NewErrResponse(http.StatusBadRequest, constants.ErrInvalidUsernamePassword)
	}
	return nil
}

func signToken(tkn token.Token, userID uuid.UUID, expiresIn time.Duration) (string, *dto.Response[*any]) {
	claims := token.NewClaims(userID.String(), expiresIn)
	tk, err := tkn.Sign(claims)
	if err != nil {
		return "", dto.NewErrResponse(http.StatusInternalServerError, err)
	}
	return tk, nil
}
