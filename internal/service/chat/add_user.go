package chat

import (
	"context"
	"fmt"

	_ "time/tzdata"

	"github.com/vedomirr/remindista/internal/domain"
	u "github.com/vedomirr/remindista/internal/entity/user"

	"go.uber.org/zap"
)

type ChatAddUser struct {
	*Chat
}

func NewChatAddUser(chat *Chat) *ChatAddUser {
	c := &ChatAddUser{chat}

	go c.chat()

	return c
}

func (c *ChatAddUser) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	user := u.NewUser(
		u.WithTelegramId(c.tgId),
		u.WithChatId(c.chatId),
		u.WithIsRunning(true),
	)

	stage := "location"
	c.SendMessage(fmt.Sprintf(domain.ReplySetLocation, user.LocationString()), domain.KbLocations)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch stage {
		case "location":
			if !c.skipped(msg) {
				if err := user.SetLocation(msg); err != nil {
					c.log.Error("failed to set location", zap.Error(err))
					c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingLocation, err).Error(), domain.KbLocations)
					break
				}
			}

			c.SendMessage(fmt.Sprintf(domain.ReplySetWindowFloor, user.WindowFloorString()), domain.KbWindowFloor)
			stage = "floor"

		case "floor":
			if !c.skipped(msg) {
				if err := user.SetWindowFloor(msg); err != nil {
					c.log.Error("failed to set floor time", zap.Error(err))
					c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingTime, err).Error(), nil)
					break
				}
			}

			c.SendMessage(fmt.Sprintf(domain.ReplySetWindowCeil, user.WindowCeilString()), domain.KbWindowCeil)
			stage = "ceil"

		case "ceil":
			if !c.skipped(msg) {
				if err := user.SetWindowCeil(msg); err != nil {
					c.log.Error("failed to set ceil time", zap.Error(err))
					c.SendMessage(domain.ReplyErrorParsingTime, nil)
					break
				}
			}

			if _, err := c.db.CreateUser(context.Background(), user); err != nil {
				c.log.Error("failed to create user", zap.Error(err))
				c.SendMessage(fmt.Errorf(domain.ReplyErrorCreatingUser, err).Error(), nil)
				return
			}

			c.SendMessage(domain.ReplyUserUpdated, nil)
			return

		default:
			c.log.Error("unknown stage", zap.String("stage", stage))
			return
		}
	}
}
