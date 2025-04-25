package routers

import (
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/internal/handlers"
	"net/http"
)

func initUserRouter(mux *http.ServeMux) {
	handler := handlers.NewUserHandler(global.Repository, global.Hash)

	mux.HandleFunc("GET /users/", handler.FindUser)
	mux.HandleFunc("POST /users/", handler.CreateUser)
}
