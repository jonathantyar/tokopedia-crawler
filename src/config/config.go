package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	if !strings.HasSuffix(os.Args[0], ".test") {
		viper.SetConfigFile(".env")
		viper.AddConfigPath(".")
		_ = viper.ReadInConfig()
	} else {
		viper.SetConfigFile("../.test.env")
		viper.AddConfigPath(".")
		_ = viper.ReadInConfig()
	}
}
