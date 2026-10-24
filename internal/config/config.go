package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
