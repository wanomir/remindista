package chat

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"
	r "github.com/vedomirr/remindista/internal/entity/reminder"
)

type ChatListReminders struct {
	*Chat
}

func NewChatListReminders(chat *Chat) *ChatListReminders {
	c := &ChatListReminders{chat}

	go c.chat()

	return c
}

func (c *ChatListReminders) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	user, err := c.getUser()
	if err != nil {
		c.SendMessage(domain.ReplyFailedFindUser, nil)
		return
	}

	rmds, err := c.db.GetRemindersByUserId(context.Background(), user.Id)
	if err != nil {
		c.SendMessage(fmt.Errorf(domain.ReplyErrorGettingReminder, err).Error(), nil)
		return
	}

	if len(rmds) == 0 {
		c.SendMessage(domain.ReplyNoReminders, domain.KbAdd)
		return
	}

	c.SendMessage(domain.ReplyListReminders, domain.KbListReminders)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch msg {
		case "all":
			for _, rmd := range rmds {
				c.SendMessage(rmd.StringMdV2(), rmd.Keyboard())
			}
			return

		default:
			rmdsTag := c.rmdsByTag(rmds, msg)
			if len(rmdsTag) == 0 {
				c.SendMessage(domain.ReplyNoRemindersWithTag, domain.KbListReminders)
				break
			}

			for _, rmd := range rmdsTag {
				c.SendMessage(rmd.StringMdV2(), rmd.Keyboard())
			}

			c.SendMessage(domain.ReplyListAnotherTag, domain.KbListReminders)
		}
	}
}

func (c *ChatListReminders) rmdsByTag(rmds []r.Reminder, tag string) []r.Reminder {
	rmdsTag := make([]r.Reminder, 0)

	for _, rmd := range rmds {
		if rmd.TagMatches(tag) {
			rmdsTag = append(rmdsTag, rmd)
		}
	}

	return rmdsTag
}
