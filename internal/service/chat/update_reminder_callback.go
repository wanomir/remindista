package chat

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"

	"go.uber.org/zap"
)

type ChatUpdateReminderById struct {
	*Chat
	rmdId int
}

func NewChatUpdateReminderById(chat *Chat, rmdId int) *ChatUpdateReminderById {
	c := &ChatUpdateReminderById{Chat: chat, rmdId: rmdId}

	go c.chat()

	return c
}

func (c *ChatUpdateReminderById) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	user, err := c.getUser()
	if err != nil {
		c.SendMessage(domain.ReplyFailedFindUser, nil)
		return
	}

	rmd, err := c.db.GetReminder(context.Background(), c.rmdId)
	if err != nil {
		c.SendMessage(fmt.Errorf(domain.ReplyErrorGettingReminder, err).Error(), nil)
		return
	}

	stage := "text"
	c.SendMessage(fmt.Sprintf(domain.ReplyUpdateReminderText, rmd.TextMdV2()), domain.KbSkip)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch stage {
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
