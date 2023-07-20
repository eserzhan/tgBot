package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/sashabaranov/go-openai"

	"github.com/eserzhan/tgBot/pkg/assembly"
	"github.com/eserzhan/tgBot/pkg/config"
	"github.com/eserzhan/tgBot/pkg/telegram"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	api.Debug = true
	clientGpt := openai.NewClient(cfg.GptToken)

	assembly, err := assembly.NewClient(cfg.AssemblyToken)
	if err != nil {
		log.Fatal(err)
	}
	
	bot := telegram.NewBot(api, clientGpt, assembly, cfg.Messages)

	bot.Run() 

}