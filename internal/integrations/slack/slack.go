package slack

import (
	messengers "CLIMultiChat/internal/integrations"
	"fmt"

	"github.com/slack-go/slack"
)

type Client struct {
	api            *slack.Client
	defaultChannel string
}

func NewClient(token, defaultChannel string) (messengers.Messenger, error) {
	if token == "" {
		return nil, fmt.Errorf("slack token is required")
	}

	api := slack.New(token)

	return &Client{
		api:            api,
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

	_, _, err := c.api.PostMessage(
		channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return fmt.Errorf("failed to send message to Slack: %w", err)
	}

	return nil
}

func (c *Client) GetName() string {
	return "Slack"
}
