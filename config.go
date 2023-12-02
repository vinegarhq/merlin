package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Port       string  `json:"port"`
	CertFile   string  `json:"cert_file"`
	KeyFile    string  `json:"key_file"`
	OutputFile string  `json:"output_file"`
	RateLimit  float64 `json:"rate_limit"` // per second
}

func LoadConfig(name string) (Config, error) {
	cf, err := os.ReadFile(name)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(cf, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.OutputFile == "" {
		return Config{}, errors.New("no output file given")
	}

	if cfg.Port == "" {
		return Config{}, errors.New("no port given")
	}

	return cfg, nil
}
