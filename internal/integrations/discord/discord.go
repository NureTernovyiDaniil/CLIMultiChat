package discord

import (
	messengers "CLIMultiChat/internal/integrations"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	session        *discordgo.Session
	defaultChannel string
}

func NewClient(token, defaultChannel string) (messengers.Messenger, error) {
	if token == "" {
		return nil, fmt.Errorf("discord token is required")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	return &Client{
		session:        session,
		defaultChannel: defaultChannel,
	}, nil
}

func (c *Client) SendMessage(channel, message string) error {
	if channel == "" {
		if c.defaultChannel == "" {
			return fmt.Errorf("channel is required")
		}
		channel = c.defaultChannel
	}

	_, err := c.session.ChannelMessageSend(channel, message)
	if err != nil {
		return fmt.Errorf("failed to send message to Discord: %w", err)
	}

	return nil
}

func (c *Client) GetName() string {
	return "Discord"
}
