package initialize

import (
	"github.com/duyanhitbe/go-ecom/internal/config"
	"github.com/duyanhitbe/go-ecom/pkg/constants"
	"github.com/spf13/viper"
	"os"
)

func LoadConfig() {
	mode := os.Getenv(constants.GoModeEnvKey)
	if mode == "" {
		mode = constants.DevelopmentMode
	}

	viper.AddConfigPath(constants.ConfigDir)
	viper.SetConfigType(constants.ConfigType)
	viper.SetConfigName(mode)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config.Cfg)
	if err != nil {
		panic(err)
	}
}
