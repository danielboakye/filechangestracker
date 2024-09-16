package config

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Directory      string `validate:"required"`
	CheckFrequency int    `validate:"required,min=1"`
	ReportingAPI   string `validate:"required,url"`
	HTTPPort       string `validate:"required"`
	SocketPath     string `validate:"required"`
}

const FileChangesLogFile = "filechangestracker.log"

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("http_port", "9000")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config.yaml: %w", err)
	}

	config := &Config{
		Directory:      viper.GetString("directory"),
		CheckFrequency: viper.GetInt("check_frequency"),
		ReportingAPI:   viper.GetString("reporting_api"),
		HTTPPort:       viper.GetString("http_port"),
		SocketPath:     viper.GetString("socket_path"),
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return config, nil
}
