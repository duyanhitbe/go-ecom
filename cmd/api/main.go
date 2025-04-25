package main

import (
	"fmt"
	"github.com/duyanhitbe/go-ecom/internal/config"
	"github.com/duyanhitbe/go-ecom/internal/initialize"
	"github.com/duyanhitbe/go-ecom/internal/server"
)

func main() {
	initialize.InitConfig()
	initialize.InitLogger()
	initialize.InitHash()
	initialize.InitDatabase()

	addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
	s := server.NewRestfulServer(addr)
	s.Start()
}
