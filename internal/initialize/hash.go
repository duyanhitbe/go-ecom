package initialize

import (
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/pkg/hash"
)

func initHash() {
	h := hash.NewBcrypt()
	global.Hash = h
}
