package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Config struct {
	Port     uint
	Database struct {
		Host           string
		Port           uint16
		Name           string
		User           string
		Password       string
		MaxConnections int `mapstructure:"max_connections"`
	} `mapstructure:"database"`
	Jwt struct {
		Secret               string
		AccessTokenLifetime  int `mapstructure:"access_token_lifetime"`
		RefreshTokenLifetime int `mapstructure:"refresh_token_lifetime"`
	} `mapstructure:"jwt"`
}

func NewConfig(configPath string) (*Config, error) {
	dir, err := filepath.Abs(filepath.Dir(configPath))
	if err != nil {
		return nil, errors.Wrap(err, "Parse config path")
	}

	configName := filepath.Base(configPath)
	configNameWithoutExt := strings.TrimSuffix(configName, filepath.Ext(configName))

	viper.AddConfigPath(dir)
	viper.SetConfigName(configNameWithoutExt)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Reading config file")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "Parsing config file")
	}

	return &cfg, nil
}
