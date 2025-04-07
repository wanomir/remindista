package chat

import (
	"context"
	"fmt"

	_ "time/tzdata"

	"github.com/vedomirr/remindista/internal/domain"

	"go.uber.org/zap"
)

type ChatUpdateUser struct {
	*Chat
}

func NewChatUpdateUser(chat *Chat) *ChatUpdateUser {
	c := &ChatUpdateUser{chat}

	go c.chat()

	return c
}

func (c *ChatUpdateUser) chat() {
	user, err := c.getUser()
	if err != nil {
		_ = NewChatAddUser(c.Chat)
		return
	}

	defer close(c.inCh)
	defer c.deleteChat()

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

			// TOD: later rewrite windows using duration to evoid bugs
			// below is a temporary fix
			func() {
				if c.skipped(msg) {
					if err := user.SetWindowFloor(user.WindowFloorString()); err != nil {
						c.log.Error("failed to set default floor time", zap.Error(err))
					}
				}
				fmt.Println("user_window_floor", user.WindowFloor)
			}()

			c.SendMessage(fmt.Sprintf(domain.ReplySetWindowCeil, user.WindowCeilString()), domain.KbWindowCeil)
			stage = "ceil"

		case "ceil":
			if !c.skipped(msg) {
				if err := user.SetWindowCeil(msg); err != nil {
					c.log.Error("failed to set ceil time", zap.Error(err))
					c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingTime, err).Error(), nil)
					break
				}
			}

			// TOD: later rewrite windows using duration to evoid bugs
			// below is a temporary fix
			func() {
				if c.skipped(msg) {
					if err := user.SetWindowCeil(user.WindowCeilString()); err != nil {
						c.log.Error("failed to set default floor time", zap.Error(err))
					}
				}
				fmt.Println("user_window_ceil", user.WindowCeil)
			}()

			if _, err := c.db.UpdateUser(context.Background(), user); err != nil {
				c.log.Error("failed to update user", zap.Error(err))
				c.SendMessage(fmt.Errorf(domain.ReplyErrorUpdatingUser, err).Error(), nil)
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
