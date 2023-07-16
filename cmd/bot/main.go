package main

import (
	"TelegramBotExchangeRate/app/bot"
	"TelegramBotExchangeRate/app/config"
	"fmt"
	"log"
)

func main() {
	fmt.Println("start")
	err := bot.Run(config.Load().Token)
	if err != nil {
		fmt.Println(err)
		log.Println("error: bot doesn't run ")
	}
}
