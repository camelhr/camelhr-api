package config

import (
	"log" //nolint:depguard // since functions defined here are called before the structured logger is initialized

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`

	HTTPAddress string `mapstructure:"http_address"`

	DBConn            string `mapstructure:"db_conn"`
	DBMaxOpen         int    `mapstructure:"db_max_open"`
	DBMaxIdle         int    `mapstructure:"db_max_idle"`
	DBMaxIdleConnTime int    `mapstructure:"db_max_idle_conn_time"`
}

const (
	defaultDBMaxOpen         = 4
	defaultDBMaxIdle         = 4
	defaultDBMaxIdleConnTime = 4
)

func init() {
	// set the default values. make sure to set the default values for all the configs
	// that are being used in the application. the default values should point to
	// the production setup unless the value is a secret or sensitive information.

	// logger configs
	viper.SetDefault("log_level", "error")

	// http server configs
	viper.SetDefault("http_address", "0.0.0.0:8080")

	// database configs
	viper.SetDefault("db_conn", "") // secret value. should be set in the environment.
	viper.SetDefault("db_max_open", defaultDBMaxOpen)
	viper.SetDefault("db_max_idle", defaultDBMaxIdle)
	viper.SetDefault("db_max_idle_conn_time", defaultDBMaxIdleConnTime) // in minutes

	// override default values with environment variables.
	viper.AutomaticEnv()
}

// LoadConfig reads the configuration from viper and returns a Config struct.
func LoadConfig() Config {
	var config Config

	// store values in the config struct.
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to read configs: %v", err)
	}

	return config
}
