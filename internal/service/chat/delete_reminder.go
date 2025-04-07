package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/vedomirr/remindista/internal/domain"
	r "github.com/vedomirr/remindista/internal/entity/reminder"

	"go.uber.org/zap"
)

type ChatDeleterReminder struct {
	*Chat
}

func NewChatDeleteReminder(chat *Chat) *ChatDeleterReminder {
	c := &ChatDeleterReminder{chat}

	go c.chat()

	return c
}

func (c *ChatDeleterReminder) chat() {
	defer close(c.inCh)
	defer c.deleteChat()

	user, err := c.getUser()
	if err != nil {
		c.SendMessage(domain.ReplyFailedFindUser, nil)
		return
	}

	rmd := r.NewReminder()

	stage := "mode"
	c.SendMessage(domain.ReplySetMode, domain.KbSetMode)

	for msg := range c.inCh {
		if c.isCancel(msg) {
			c.SendMessage(domain.ReplyCancel, nil)
			return
		}

		switch stage {
		case "mode":

			switch msg {
			case "id":
				c.SendMessage(domain.ReplySendId, domain.KbCancel)
				stage = "id"
			case "tag":
				c.SendMessage(domain.ReplySendTag, domain.KbCancel)
				stage = "tag"
			case "all":
				c.SendMessage(domain.ReplyConfirmDeletingAll, domain.KbYesNo)
				stage = "all"
			case "stop":
				stage = "stop"
			default:
				c.SendMessage(domain.ReplyUnknown, nil)
				c.SendMessage(domain.ReplyModes, domain.KbSetMode)
			}

		case "id":
			if err := rmd.HexToId(msg); err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorParsingId, err).Error(), nil)
				break
			}

			if nDeleted, err := c.db.DeleteReminder(context.Background(), rmd.Id); err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorDeletingReminder, err).Error(), nil)
				break
			} else if nDeleted == 0 {
				c.SendMessage(domain.ReplyNoSuchId, nil)
				break
			}

			c.SendMessage(domain.ReplyDone, nil)
			return

		case "tag":
			if err := rmd.SetTag(msg); err != nil {
				c.SendMessage(domain.ReplyErrorParsingTag, domain.KbSkip)
				break
			}

			nDeleted, err := c.db.DeleteRemindersByTag(context.Background(), user.Id, rmd.Tag)
			if err != nil {
				c.SendMessage(fmt.Errorf(domain.ReplyErrorDeletingReminder, err).Error(), nil)
				break
			} else if nDeleted == 0 {
				c.SendMessage(domain.ReplyNoSuchTag, domain.KbCancel)
				break
			}

			c.SendMessage(fmt.Sprintf(domain.ReplyDeletedMultiple, nDeleted), nil)
			return

		case "all":
			msg := strings.ToLower(msg)

			switch msg {
			case "yes":
				nDeleted, err := c.db.DeleteRemindersByUserId(context.Background(), user.Id)
				if err != nil {
					c.SendMessage(fmt.Errorf(domain.ReplyErrorDeletingReminder, err).Error(), nil)
					break
				} else if nDeleted == 0 {
					c.SendMessage(domain.ReplyNoReminders, nil)
					break
				}

				c.SendMessage(fmt.Sprintf(domain.ReplyDeletedMultiple, nDeleted), nil)
				return

			case "no":
				c.SendMessage(domain.ReplyCancel, nil)
				c.SendMessage(domain.ReplyModes, domain.KbSetMode)
				stage = "mode"

			default:
				c.SendMessage(domain.ReplyYesNo, domain.KbYesNo)
			}

		default:
			c.log.Error("unknown stage", zap.String("stage", stage))
			return
		}
	}
}
