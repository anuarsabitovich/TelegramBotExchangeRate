package bot

import (
	"TelegramBotExchangeRate/app/exchange"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func connect(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token) //@todo Naming

	if err != nil {
		return nil, fmt.Errorf("connecting to bot: %w", err)
	}

	bot.Debug = true //@todo : Do not set fields

	return bot, nil
}

func getUpdatesChan(b *tgbotapi.BotAPI, offset, to int) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = to                      //@todo : Do not set fields
	updates, err := b.GetUpdatesChan(u) //@todo : err handling. Extract.
	if err != nil {
		return nil, fmt.Errorf("getting updates chan: %w", err)
	}
	return updates, nil
}

func isMessageProcess(text string) bool {
	return strings.HasPrefix(text, "/currency")
}

func Run(token string) error {
	b, err := connect(token)
	if err != nil {
		return fmt.Errorf("run bot: %w", err)
	}

	updates, err := getUpdatesChan(b, 0, 60)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if isMessageProcess(update.Message.Text) {
			currencyOptions := []string{"AUD", "USD", "RUB"}
			var rows []tgbotapi.KeyboardButton
			for _, currency := range currencyOptions {
				button := tgbotapi.NewKeyboardButton(currency)
				rows = append(rows, button)
			}
			replyMarkup := tgbotapi.NewReplyKeyboard(rows)
			replyMarkup.OneTimeKeyboard = true
			response := "Please select a currency:"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyMarkup = replyMarkup
			b.Send(msg)
			continue

		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		fmt.Println(msg.Text)

		msg.ReplyToMessageID = update.Message.MessageID

		currency, err := exchange.GetCurrentRate(msg.Text)
		if err != nil {
			fmt.Println("error currency exchange")
		}
		//@todo : template
		message := fmt.Sprintf("Hello, %s!\n%s %s", update.Message.From.UserName, currency.Title, currency.Current)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, message)
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := b.Send(msg); err != nil {
			fmt.Println(err.Error())
		}
		continue
	}
	return nil
}
