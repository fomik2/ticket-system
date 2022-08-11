package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		HTTP      `yaml:"http"`
		Templates `yaml:"templates"`
		Files     `yaml:"files"`
		CSS       `yaml:"css"`
	}
	// HTTP -.
	HTTP struct {
		Port string `yaml:"port"`
	}

	// Templates -.
	Templates struct {
		Basedir string `yaml:"basedir"`
		Index   string `yaml:"index"`
		Editor  string `yaml:"editor"`
	}

	// Files -.
	Files struct {
		Tasks   string `yaml:"tasks"`
		Counter string `yaml:"counter"`
	}

	CSS struct {
		Path string `yaml:"path"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
