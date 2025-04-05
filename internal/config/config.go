package config

import (
	"strings"
	"url_shortener/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port    string `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	viper.SetEnvPrefix("URL_SHORTENER") // Префикс для ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Warn("config.yaml not found, loading from environment variables")
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		logger.Log.Fatal("Failed to parse configuration: ", err)
	}

	if AppConfig.Database.User == "" || AppConfig.Database.Password == "" {
		logger.Log.Fatal("Missing required storage configuration (user/password)")
	}

	logger.Log.Info("Configuration loaded successfully")
}
