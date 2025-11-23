package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SlackToken       string
	SlackChannel     string
	TelegramBotToken string
	TelegramChatID   string
	DiscordToken     string
	DiscordChannel   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		SlackToken:       os.Getenv("SLACK_TOKEN"),
		SlackChannel:     os.Getenv("SLACK_CHANNEL"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		DiscordChannel:   os.Getenv("DISCORD_CHANNEL"),
	}

	return config, nil
}

func (c *Config) ValidateSlack() error {
	if c.SlackToken == "" {
		return fmt.Errorf("SLACK_TOKEN is missing")
	}
	return nil
}

func (c *Config) ValidateTelegram() error {
	if c.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is missing")
	}
	return nil
}

func (c *Config) ValidateDiscord() error {
	if c.DiscordToken == "" {
		return fmt.Errorf("DISCORD_TOKEN is missing")
	}
	return nil
}
