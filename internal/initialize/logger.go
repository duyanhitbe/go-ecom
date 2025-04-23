package initialize

import (
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/pkg/logger"
)

func InitLogger() {
	lg := logger.NewZapLogger()

	global.Logger = lg
}
