package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConf
	Cache  CacheConf
	Logger LoggerConf
}

type CacheConf struct {
	Size int
}
type ServerConf struct {
	Host string
	Port string
}

type LoggerConf struct {
	Level string
	File  string
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}

	return &config, nil
}

func (s *ServerConf) DSN() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
