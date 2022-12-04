package config

import (
	"github.com/spf13/viper"
)

const configPath = "./env"

type AppConfig struct {
	ServerConfig   ServerConfig   `mapstructure:",squash"`
	DatabaseConfig DatabaseConfig `mapstructure:",squash"`
}
type ServerConfig struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	ReadTimeout   int    `mapstructure:"READ_TIMEOUT"`
	WriteTimeout  int    `mapstructure:"WRITE_TIMEOUT"`
}
type DatabaseConfig struct {
	Driver   string `mapstructure:"DB_DRIVER"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	Username string `mapstructure:"DB_USERNAME"`
}

func LoadConfig() (*AppConfig, error) {
	var config AppConfig
	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	return &config, err
}
