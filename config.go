package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	// Port to use
	Port string `json:"port"`

	// Certificate (pubkey) file to use
	CertFile string `json:"cert_file"`

	// Key (privatekey) file to use
	KeyFile string `json:"key_file"`

	// Output CSV file
	OutputFile string `json:"output_file"`

	// Requests to allow per second
	RateLimit float64 `json:"rate_limit"`

	// Setting to change where tollbooth checks for IPs.
	// It is CRITICAL that you set this properly, otherwise you will blacklist a cloudflare node.
	CFMode bool `json:"cf_mode"`

	// mTLS cert location for Cloudflare. (PEM ONLY)
	// Automatically enables mTLS when provided.
	MTLSFile string `json:"mtls_file"`
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
