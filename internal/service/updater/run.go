package updater

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"

	"go.uber.org/zap"
)

func (u *Updater) Run(ctx context.Context) {
	incoming := u.telegram.ReceiveMessages(ctx)
	defer close(incoming)

	// callbacks u.telegram.ReceiveCallbacks(ctx)

	u.outCh = make(chan domain.Message)
	defer close(u.outCh)

	// send out messages that are ready
	go func(outgoing chan domain.Message) {
		for msg := range outgoing {
			u.sendOut(msg)
		}
	}(u.outCh)

	// process incoming messages in goroutines
	limit := make(chan struct{}, 100)
	for {
		select {
		case msg := <-incoming:
			limit <- struct{}{}
			go u.processIncoming(msg, limit)

		case <-ctx.Done():
			u.log.Info("stopped telegram updater")
			return
		}
	}
}

func (u *Updater) Stop() {}

func (u *Updater) sendOut(message domain.Message) {
	if err := u.telegram.SendMessageMarkdownV2(message.ChatId, message.Text, message.Keyboard); err != nil {
		u.log.Error(fmt.Sprintf("failed to send message [%s] %s (id: %v, chatId: %v)", message.UserName, message.Text, message.TelegramId, message.ChatId), zap.Error(err))
		return
	}
}

func (u *Updater) processIncoming(message domain.Message, limit chan struct{}) {
	defer func() { <-limit }()

	u.log.Info(fmt.Sprintf("[%s] %s (id: %v, chatId: %v)", message.UserName, message.Text, message.TelegramId, message.ChatId))

	if err := u.ProcessMessage(message); err != nil {
		u.log.Error(fmt.Sprintf("failed to process message [%s] %s (id: %v, chatId: %v)", message.UserName, message.Text, message.TelegramId, message.ChatId), zap.Error(err))
		return
	}
}
