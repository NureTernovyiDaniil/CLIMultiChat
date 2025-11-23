package main

import (
	"CLIMultiChat/internal/config"
	messengers "CLIMultiChat/internal/integrations"
	"CLIMultiChat/internal/integrations/discord"
	"CLIMultiChat/internal/integrations/slack"
	"CLIMultiChat/internal/integrations/telegram"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfg         *config.Config
	sendAllFlag bool
	slackClient messengers.Messenger
)

func init() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "climessenger",
	Short: "Відправка повідомлень у месенджери",
	Long:  `Програма для відправки повідомлень у Slack, Telegram та Discord.`,
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Відправити повідомлення",
	Long:  `Відправити повідомлення у Slack, Telegram або Discord.`,
}

var slackCmd = &cobra.Command{
	Use:   "slack [канал] [повідомлення]",
	Short: "В Slack",
	Long:  `Відправити повідомлення у Slack. Якщо канал не вказано, використовується стандартний.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var channel, message string

		if len(args) == 1 {
			channel = ""
			message = args[0]
		} else {
			channel = args[0]
			message = args[1]
		}

		if err := cfg.ValidateSlack(); err != nil {
			return fmt.Errorf("slack configuration error: %w", err)
		}

		client, err := slack.NewClient(cfg.SlackToken, cfg.SlackChannel)
		if err != nil {
			return fmt.Errorf("failed to create Slack client: %w", err)
		}

		if err := client.SendMessage(channel, message); err != nil {
			return err
		}

		fmt.Printf("Повідомлення надіслано у Slack\n")
		return nil
	},
}

var telegramCmd = &cobra.Command{
	Use:   "telegram [chat_id] [повідомлення]",
	Short: "В Telegram",
	Long:  `Відправити повідомлення у Telegram. Якщо chat_id не вказано, використовується стандартний.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var chatID, message string

		if len(args) == 1 {
			chatID = ""
			message = args[0]
		} else {
			chatID = args[0]
			message = args[1]
		}

		if err := cfg.ValidateTelegram(); err != nil {
			return fmt.Errorf("telegram configuration error: %w", err)
		}

		client, err := telegram.NewClient(cfg.TelegramBotToken, cfg.TelegramChatID)
		if err != nil {
			return fmt.Errorf("failed to create Telegram client: %w", err)
		}

		if err := client.SendMessage(chatID, message); err != nil {
			return err
		}

		fmt.Printf("Повідомлення надіслано у Telegram\n")
		return nil
	},
}

var discordCmd = &cobra.Command{
	Use:   "discord [канал] [повідомлення]",
	Short: "В Discord",
	Long:  `Відправити повідомлення у Discord. Якщо канал не вказано, використовується стандартний.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var channel, message string

		if len(args) == 1 {
			channel = ""
			message = args[0]
		} else {
			channel = args[0]
			message = args[1]
		}

		if err := cfg.ValidateDiscord(); err != nil {
			return fmt.Errorf("discord configuration error: %w", err)
		}

		client, err := discord.NewClient(cfg.DiscordToken, cfg.DiscordChannel)
		if err != nil {
			return fmt.Errorf("failed to create Discord client: %w", err)
		}

		if err := client.SendMessage(channel, message); err != nil {
			return err
		}

		fmt.Printf("Повідомлення надіслано у Discord\n")
		return nil
	},
}

var allCmd = &cobra.Command{
	Use:   "all [повідомлення]",
	Short: "Усі месенджери",
	Long:  `Відправити одне повідомлення у всі месенджери (Slack, Telegram і Discord).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]
		errors := []error{}

		if cfg.ValidateSlack() == nil {
			client, err := slack.NewClient(cfg.SlackToken, cfg.SlackChannel)
			if err != nil {
				errors = append(errors, fmt.Errorf("slack: %w", err))
			} else {
				if err := client.SendMessage("", message); err != nil {
					errors = append(errors, fmt.Errorf("Slack: %w", err))
				} else {
					fmt.Printf("Повідомлення надіслано у Slack\n")
				}
			}
		} else {
			fmt.Println("Slack не налаштовано, пропущено")
		}

		if cfg.ValidateTelegram() == nil {
			client, err := telegram.NewClient(cfg.TelegramBotToken, cfg.TelegramChatID)
			if err != nil {
				errors = append(errors, fmt.Errorf("telegram: %w", err))
			} else {
				if err := client.SendMessage("", message); err != nil {
					errors = append(errors, fmt.Errorf("Telegram: %w", err))
				} else {
					fmt.Printf("Повідомлення надіслано у Telegram\n")
				}
			}
		} else {
			fmt.Println("Telegram не налаштовано, пропущено")
		}

		if cfg.ValidateDiscord() == nil {
			client, err := discord.NewClient(cfg.DiscordToken, cfg.DiscordChannel)
			if err != nil {
				errors = append(errors, fmt.Errorf("discord: %w", err))
			} else {
				if err := client.SendMessage("", message); err != nil {
					errors = append(errors, fmt.Errorf("Discord: %w", err))
				} else {
					fmt.Printf("Повідомлення надіслано у Discord\n")
				}
			}
		} else {
			fmt.Println("Discord не налаштовано, пропущено")
		}

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "Помилка: %v\n", err)
			}
			return fmt.Errorf("не вдалося надіслати у всі месенджери")
		}

		return nil
	},
}

func main() {
	sendCmd.AddCommand(slackCmd)
	sendCmd.AddCommand(telegramCmd)
	sendCmd.AddCommand(discordCmd)
	sendCmd.AddCommand(allCmd)
	rootCmd.AddCommand(sendCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
