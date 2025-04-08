package tg

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"
)

type Sender struct {
	log zerolog.Logger

	client      *tbot.Client
	serviceName string
	chatID      string
}

func NewSender(telegramToken, telegramChatID, serviceName string, log zerolog.Logger) *Sender {
	return &Sender{
		client:      tbot.New(telegramToken).Client(),
		serviceName: serviceName,
		chatID:      telegramChatID,
		log:         log,
	}
}

func (s *Sender) SendError(err error) {
	if err != nil {
		_, sendErr := s.client.SendMessage(s.chatID, fmt.Sprintf("Error: %v\nServiceName: %s", err, s.serviceName))
		if sendErr != nil {
			s.log.Error().Err(sendErr).Msgf("Error sending message")
		}
	}
}
