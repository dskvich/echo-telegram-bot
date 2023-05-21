package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, update tgbotapi.Update) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("failed to create Telegram bot: %w", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch {
		case strings.HasPrefix(update.Message.Text, "/userid"):
			handleUserIDCommand(ctx, bot, update)
		default:
			handleEchoMessage(ctx, bot, update)
		}
	}
}

func handleUserIDCommand(_ context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	response := fmt.Sprintf("Your user ID is %d", update.Message.From.ID)
	sendTelegramMessage(bot, update.Message.Chat.ID, response, update.Message.MessageID)
}

func sendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, text string, replyToMessageID int) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = replyToMessageID
	if _, err := bot.Send(msg); err != nil {
		log.Fatalf("failed to send message via Telegram bot: %w", err)
	}
}

func handleEchoMessage(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	sendTelegramMessage(bot, update.Message.Chat.ID, update.Message.Text, update.Message.MessageID)
}
