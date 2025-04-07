package chat

import (
	"context"
	"errors"
	"strings"

	"github.com/vedomirr/l"
	"github.com/vedomirr/remindista/internal/domain"
	u "github.com/vedomirr/remindista/internal/entity/user"

	"go.uber.org/zap"
)

type Chat struct {
	chatId, tgId int64

	inCh     chan string
	outCh    chan domain.Message
	deleteCh chan int64

	db repository

	log *zap.Logger
}

func NewChat(chatId, tgId int64, outCh chan domain.Message, deleteCh chan int64, db repository) *Chat {
	return &Chat{
		chatId: chatId,
		tgId:   tgId,

		inCh:     make(chan string),
		outCh:    outCh,
		deleteCh: deleteCh,

		db: db,

		log: l.Logger(),
	}
}

func (c *Chat) SendMessage(text string, keyboard domain.Keyboard) {
	c.outCh <- domain.Message{UserName: "Remindista", ChatId: c.chatId, Text: text, Keyboard: keyboard}
}

func (c *Chat) PassInput(input string) {
	c.inCh <- input
}

//nolint:golint,unused
func (c *Chat) chat() {
	defer close(c.inCh)

	c.SendMessage("Started new chat!", nil)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			return
		}

		c.SendMessage(msg, nil)
	}

	c.SendMessage("That's enough!", nil)
	c.deleteChat()
}

func (c *Chat) deleteChat() {
	c.deleteCh <- c.chatId
}

func (c *Chat) skipped(s string) bool {
	return strings.ToLower(strings.TrimSpace(s)) == "skip"
}

func (c *Chat) getUser() (user u.User, err error) {
	if user, err = c.db.GetUserByTelegramId(context.Background(), c.tgId); err != nil || user.TelegramId == 0 {
		if err != nil {
			c.log.Error("failed to get user", zap.Int64("chat id", c.chatId), zap.Int64("telegram id", c.tgId), zap.Error(err))
		}
		return user, errors.New("failed to get user")
	}

	// update chat id associated with this user if needed
	if user.ChatId != c.chatId {
		user.ChatId = c.chatId
		if _, err = c.db.UpdateUser(context.Background(), user); err != nil {
			return user, err
		}
	}

	return user, nil
}

func (c *Chat) isCancel(msg string) bool {
	return strings.ToLower(msg) == "cancel"
}
