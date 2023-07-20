package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error{
	switch message.Command(){
	case startCommand:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Start)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	resp, err := b.gptAnswer(message.Text)
	if err != nil {
		return fmt.Errorf("ChatCompletion error: %v", err)
	}
	
	msg := tgbotapi.NewMessage(message.Chat.ID, resp.Choices[0].Message.Content)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleAudioMessage(message *tgbotapi.Message) error {
			voice := message.Voice

		   fileID := voice.FileID

			fileUrl, err := b.bot.GetFileDirectURL(fileID)
			if err != nil {
				return err
			}

			transcription, err := b.assembly.Transcription(fileUrl)
			if err != nil {
				return err 
			}

		    audioResponse, err := b.assembly.TranscribedText(transcription)
			if err != nil {
				return err 
			}

			message.Text = audioResponse
			
			return b.handleMessage(message)
}

func(b *Bot) gptAnswer(message string) (response openai.ChatCompletionResponse, err error){
	resp, err := b.gpt.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	return resp, err
}

