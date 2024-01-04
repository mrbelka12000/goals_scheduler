package main

import (
	"fmt"
	"os"
	"reflect"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/rs/zerolog"

	"problem_parser_bot/pkg/config"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	cfg, err := config.Get()
	if err != nil {
		log.Fatal().Err(err).Msg("get config")
	}

	fmt.Printf("%+v\n", cfg)

	telegramBot(cfg)
}

func telegramBot(cfg config.Config) {

	//Create bot
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		panic(err)
	}

	//Set update timeout
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Get updates from bot
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//Check if message from user is text
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":
				//Send message
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, from mrbelka12000")
				bot.Send(msg)
			default:
			}
		} else {
			//Send message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}
