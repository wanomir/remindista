package chat

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"
)

type ChatDeleteReminderById struct {
	*Chat
	rmdId int
}

func NewChatDeleteReminderById(chat *Chat, rmdId int) *ChatDeleteReminderById {
	c := &ChatDeleteReminderById{chat, rmdId}

	go c.chat()

	return c
}

func (c *ChatDeleteReminderById) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	c.SendMessage(domain.ReplyeConfirmDelete, domain.KbYesNo)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch msg {
		case "yes":
			if nDeleted, err := c.db.DeleteReminder(context.Background(), c.rmdId); err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorDeletingReminder, err).Error(), nil)
				break

			} else if nDeleted == 0 {
				c.SendMessage(domain.ReplyNoSuchId, nil)
				break
			}

			c.SendMessage(domain.ReplyDone, nil)
			return

		case "no":
			c.SendMessage(domain.ReplyCancel, nil)
			return

		default:
			c.SendMessage(domain.ReplyeConfirmDelete, domain.KbYesNo)
		}
	}
}
