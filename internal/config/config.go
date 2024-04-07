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
	// logger configs
	viper.SetDefault("log_level", "error")

	// http server configs
	viper.SetDefault("http_address", "0.0.0.0:8080")

	// database configs
	viper.SetDefault("db_conn", "host=localhost port=5433 dbname=camelhr_api_dev sslmode=disable user=postgres")
	viper.SetDefault("db_max_idle", 4)
	viper.SetDefault("db_max_open", 4)
	viper.SetDefault("db_max_read_queue", 10)
	viper.SetDefault("db_max_write_queue", 10)

	// override default values with environment variables
	viper.AutomaticEnv()
}

func LoadConfig() Config {
	var config Config

	// store the values in the config struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to read configs: %v", err)
	}

	return config
}
