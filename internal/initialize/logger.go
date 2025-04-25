package initialize

import (
	"github.com/duyanhitbe/go-ecom/internal/global"
	"github.com/duyanhitbe/go-ecom/pkg/logger"
)

func initLogger() {
	lg := logger.NewZapLogger()

	global.Logger = lg
}
