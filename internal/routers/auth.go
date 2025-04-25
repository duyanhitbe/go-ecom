package routers

import (
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/internal/handlers"
	"net/http"
)

func initAuthRouter(mux *http.ServeMux) {
	handler := handlers.NewAuthHandler(global.Repository, global.Token, global.Hash)

	mux.HandleFunc("POST /auth/login/", handler.Login)
}
