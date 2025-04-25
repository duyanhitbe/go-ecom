package initialize

import (
	"github.com/duyanhitbe/go-ecom/internal/config"
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/pkg/token"
	"log"
)

func initToken() {
	jwt, err := token.NewJwt(config.Cfg.Token.Secret)
	if err != nil {
		log.Fatalf("fail to init token: %v", err)
	}
	global.Token = jwt
}
