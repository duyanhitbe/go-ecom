package routers

import "net/http"

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	initUserRouter(mux)
	initAuthRouter(mux)

	return mux
}
