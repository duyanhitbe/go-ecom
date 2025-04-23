package server

import (
	"fmt"
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/internal/routers"
	"log"
	"net/http"
)

type RestfulServer struct {
	addr string
}

func NewRestfulServer(addr string) *RestfulServer {
	return &RestfulServer{
		addr: addr,
	}
}

func (s *RestfulServer) Start() {
	handler := routers.NewRouter()

	srv := &http.Server{
		Addr:    s.addr,
		Handler: handler,
	}

	global.Logger.Info(fmt.Sprintf("Server is running on %s", s.addr))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Fail to listen : %v", err)
	}
}
