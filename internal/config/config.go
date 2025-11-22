package config

import (
	"fmt"
	"os"
)

type Config struct {
	SlackToken   string
	SlackChannel string
}

func Load() (*Config, error) {
	config := &Config{
		SlackToken:   os.Getenv("SLACK_TOKEN"),
		SlackChannel: os.Getenv("SLACK_CHANNEL"),
	}

	return config, nil
}

func (c *Config) ValidateSlack() error {
	if c.SlackToken == "" {
		return fmt.Errorf("SLACK_TOKEN is missing")
	}
	return nil
}
