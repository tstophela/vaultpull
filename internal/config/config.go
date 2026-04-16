package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	VaultPath  string
	OutputFile string
	Namespace  string
}

// Load reads configuration from environment variables and an optional .env file.
func Load(envFile string) (*Config, error) {
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	cfg := &Config{
		VaultAddr:  getEnv("VAULT_ADDR", "http://127.0.0.1:8200"),
		VaultToken: os.Getenv("VAULT_TOKEN"),
		VaultPath:  os.Getenv("VAULT_PATH"),
		OutputFile: getEnv("VAULTPULL_OUTPUT", ".env"),
		Namespace:  os.Getenv("VAULT_NAMESPACE"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required but not set")
	}
	if c.VaultPath == "" {
		return errors.New("VAULT_PATH is required but not set")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
