package config

import (
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/kr0uch/studspace/internal/transport/api/servers"
	"github.com/kr0uch/studspace/pkg/database/postgres"
	"github.com/kr0uch/studspace/pkg/logger"
)

type Config struct {
	postgres.PostgresConfig `yaml:"postgres"`
	servers.ServerConfig    `yaml:"server"`
	logger.LoggerConfig     `yaml:"logger"`
	VK                      servers.VKConfig `yaml:"vk"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".yaml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
