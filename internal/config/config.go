package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`

	HTTPAddress string `mapstructure:"http_address"`

	DBConn    string `mapstructure:"db_conn"`
	DBMaxIdle int    `mapstructure:"db_max_idle"`
	DBMaxOpen int    `mapstructure:"db_max_open"`
}

func init() {
	// set the default values. make sure to set the default values for all the configs
	// that are being used in the application. the default values should point point to
	// the production setup unless the value is a secret or sensitive information

	// logger configs
	viper.SetDefault("log_level", "error")

	// http server configs
	viper.SetDefault("http_address", "0.0.0.0:8080")

	// database configs
	viper.SetDefault("db_conn", "") // secret value. should be set in the environment
	viper.SetDefault("db_max_idle", 4)
	viper.SetDefault("db_max_open", 4)

	// override default values with environment variables
	viper.AutomaticEnv()
}

func LoadConfig() Config {
	var config Config

	// store values in the config struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to read configs: %v", err)
	}

	return config
}
