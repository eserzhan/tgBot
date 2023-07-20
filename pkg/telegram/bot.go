package telegram

import (
	"log"

	"github.com/eserzhan/tgBot/pkg/assembly"
	"github.com/eserzhan/tgBot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"github.com/eserzhan/tgBot/pkg/config"
)

const (
	startCommand = "start"

	//startMessage = "Привет! Чтобы начать сохранять ссылки в своем Pocket аккаунте, для начала тебе необходимо дать мне на это доступ. Для этого переходи по ссылке:\n%s"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	gpt *openai.Client
	assembly *assembly.Client
	messages config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, gpt *openai.Client, assembly *assembly.Client, messages config.Messages) *Bot {
	return &Bot{bot: bot, gpt: gpt, assembly: assembly, messages: messages}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // If we got a message
			continue 
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				// b.errorHandler(update.Message.Chat.ID, err)
				continue
			}
			continue
		}

	    if update.Message.Voice != nil {
			if err := b.handleAudioMessage(update.Message); err != nil {
				logger.Error(err)
				continue
			}
		}else if update.Message.Text != "" {
			if err := b.handleMessage(update.Message); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}