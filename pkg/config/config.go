package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Directory      string `validate:"required"`
	CheckFrequency int    `validate:"required,min=1"`
	ReportingAPI   string `validate:"required,url"`
	HTTPPort       string `validate:"required"`
	SocketPath     string `validate:"required"`
	LogFile        string `validate:"required"`
}

const fileChangesLogFile = "filechangestracker.log"

func LoadConfig(name, path string) (*Config, error) {
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault("http_port", "9000")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config.yaml: %w", err)
	}

	sanitizedDir := strings.ReplaceAll(viper.GetString("directory"), "'", "''")
	validPath := regexp.MustCompile(`^[a-zA-Z0-9/_-]+$`)
	if !validPath.MatchString(sanitizedDir) {
		return nil, fmt.Errorf("invalid directory format")
	}

	cfg := &Config{
		Directory:      sanitizedDir,
		CheckFrequency: viper.GetInt("check_frequency"),
		ReportingAPI:   viper.GetString("reporting_api"),
		HTTPPort:       viper.GetString("http_port"),
		SocketPath:     viper.GetString("socket_path"),
		LogFile:        fileChangesLogFile,
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return cfg, nil
}
