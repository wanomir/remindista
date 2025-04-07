package chat

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"
	r "github.com/vedomirr/remindista/internal/entity/reminder"

	"go.uber.org/zap"
)

type ChatAddReminder struct {
	*Chat
}

func NewChatAddReminder(chat *Chat) *ChatAddReminder {
	c := &ChatAddReminder{chat}

	go c.chat()

	return c
}

func (c *ChatAddReminder) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	// check if user exists
	user, err := c.getUser()
	if err != nil {
		c.SendMessage(domain.ReplyFailedFindUser, nil)
		return
	}

	// initialize reminder object
	rmd := r.NewReminder(r.WithUserId(user.Id))

	stage := "text"
	c.SendMessage(domain.ReplySetReminderText, domain.KbCancel)

	for msg := range c.inCh {
		if c.isCancel(msg) { // cancel the process
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch stage {
		case "text":
			rmd.Text = msg

			c.SendMessage(domain.ReplySetReminderTag, domain.KbSkip)
			stage = "tag"

		case "tag":
			if !c.skipped(msg) { // do only if this stage was not skipped
				if err := rmd.SetTag(msg); err != nil {
					c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingTag, err).Error(), domain.KbSkip)
					break
				}
			}

			c.SendMessage(domain.ReplySetReminderPrompt, domain.KbSkip)
			stage = "prompt"

		case "prompt":
			if !c.skipped(msg) { // do only if adding prompt was not skipped
				rmd.Prompt = msg
			}

			c.SendMessage(domain.ReplySetReminderFrequency, domain.KbCancel)
			stage = "frequency"

		case "frequency":
			if c.skipped(msg) {
				c.SendMessage(domain.ReplyCannotSkip, nil)
				c.SendMessage(domain.ReplySetReminderFrequency, domain.KbCancel)
				break
			}
			if err := rmd.SetFrequency(msg); err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingFrequency, err).Error(), domain.KbCancel)
				break
			}
			rmd.UpdateNextReminder(user.Time(), user.FloorDuration(), user.CeilDuration())

			if _, err := c.db.CreateReminder(context.Background(), rmd); err != nil {
				c.log.Error("failed to create reminder", zap.Error(err))
				c.SendMessage(fmt.Errorf(domain.ReplyErrorCreatingReminder, err).Error(), nil)
				return
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyReminderSet, rmd.NextReminderString()), nil)
			return

		default:
			c.log.Error("unknown stage", zap.String("stage", stage))
			return
		}
	}
}
