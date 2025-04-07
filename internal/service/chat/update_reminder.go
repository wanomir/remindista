package chat

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"
	r "github.com/vedomirr/remindista/internal/entity/reminder"

	"go.uber.org/zap"
)

type ChatUpdateReminder struct {
	*Chat
}

func NewChatUpdateReminder(chat *Chat) *ChatUpdateReminder {
	c := &ChatUpdateReminder{chat}

	go c.chat()

	return c
}

func (c *ChatUpdateReminder) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	user, err := c.getUser()
	if err != nil {
		c.SendMessage(domain.ReplyFailedFindUser, nil)
		return
	}

	rmd := r.NewReminder()

	stage := "id"
	c.SendMessage(domain.ReplySendId, domain.KbCancel)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch stage {
		case "id":
			if err := rmd.HexToId(msg); err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingId, err).Error(), nil)
				break
			}

			if rmd, err = c.db.GetReminder(context.Background(), rmd.Id); err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorGettingReminder, err).Error(), nil)
				return

			} else if rmd.Id == 0 {
				c.SendMessage(domain.ReplyNoSuchId, nil)
				break
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyUpdateReminderText, rmd.TextMdV2()), domain.KbSkip)
			stage = "text"

		case "text":
			if !c.skipped(msg) {
				rmd.Text = msg
			}

			tag := rmd.TagMdV2()
			if tag == "" {
				tag = "no tag"
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyUpdateReminderTag, tag), domain.KbSkip)
			stage = "tag"

		case "tag":
			if !c.skipped(msg) {
				if err := rmd.SetTag(msg); err != nil {
					c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingTag, err).Error(), domain.KbSkip)
					break
				}
			}

			if rmd.Prompt == "" {
				rmd.Prompt = domain.ReplyNoPromt
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyUpdateReminderPrompt, rmd.PromptMdV2()), domain.KbSkip)
			stage = "prompt"

		case "prompt":
			if !c.skipped(msg) {
				rmd.Prompt = msg
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyUpdateReminderFrequency, rmd.FreqeuncyString()), domain.KbSkip)
			stage = "frequency"

		case "frequency":
			if !c.skipped(msg) {
				if err := rmd.SetFrequency(msg); err != nil {
					c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingFrequency, err).Error(), nil)
					break
				}
				rmd.UpdateNextReminder(user.Time(), user.FloorDuration(), user.CeilDuration())
			}

			if _, err := c.db.UpdateReminder(context.Background(), rmd); err != nil {
				c.log.Error("failed to update reminder", zap.Error(err))
				c.SendMessage(fmt.Errorf(domain.ReplyErrorUpdatingReminder, err).Error(), nil)
				return
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyReminderSet, rmd.NextReminderString()), nil)
			return
		}
	}
}
