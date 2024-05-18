package config

import (
	"crypto/rand"
	"encoding/base64"
	"log" //nolint:depguard // since functions defined here are called before the structured logger is initialized

	"github.com/spf13/viper"
)

type Config struct {
	AppSecret string `mapstructure:"app_secret"`

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
	// that are being used in the application. the default values should be suitable for production
	// to avoid any runtime issues.

	// secure configs
	// app secret is used to sign the jwt tokens.
	// it is set to a random value by default. it must be set in the environment variable for production.
	viper.SetDefault("app_secret", generateDefaultRandomAppSecret())

	// logger configs
	viper.SetDefault("log_level", "info")

	// http server configs
	// 0.0.0.0 is a non-routable meta-address used to listen on all available network interfaces
	// so it can be accessed via any IP address that the machine has.
	viper.SetDefault("http_address", "0.0.0.0:8080")

	// database configs
	viper.SetDefault("db_conn", "") // secret value. must be set in the environment.
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

// generateDefaultRandomAppSecret generates a random app secret.
func generateDefaultRandomAppSecret() string {
	const length = 32
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("failed to generate default random app secret: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
