package telegram

import (
	"fmt"
	"strconv"

	"CLIMultiChat/internal/formatter"
	messengers "CLIMultiChat/internal/integrations"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Client struct {
	bot           *tgbotapi.BotAPI
	defaultChatID int64
}

func NewClient(token, defaultChatID string) (messengers.Messenger, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram bot token is required")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	var chatID int64
	if defaultChatID != "" {
		chatID, err = strconv.ParseInt(defaultChatID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid default chat ID: %w", err)
		}
	}

	return &Client{
		bot:           bot,
		defaultChatID: chatID,
	}, nil
}

func (c *Client) SendMessage(chatIDStr, message string) error {
	var chatID int64
	var err error

	if chatIDStr == "" {
		if c.defaultChatID == 0 {
			return fmt.Errorf("chat ID is required")
		}
		chatID = c.defaultChatID
	} else {
		chatID, err = strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid chat ID: %w", err)
		}
	}

	formattedMessage := formatter.ToTelegramMarkdown(message)

	msg := tgbotapi.NewMessage(chatID, formattedMessage)
	msg.ParseMode = "MarkdownV2"

	_, err = c.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Telegram: %w", err)
	}

	return nil
}

func (c *Client) GetName() string {
	return "Telegram"
}
