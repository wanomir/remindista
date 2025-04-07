package updater

import (
	"context"
	"fmt"

	"github.com/vedomirr/remindista/internal/domain"
	"github.com/vedomirr/remindista/internal/service/chat"

	"go.uber.org/zap"
)

func (u *Updater) processCmd(m domain.Message) {
	baseChat := chat.NewChat(m.ChatId, m.TelegramId, u.outCh, u.deleteChatCh, u.db)

	switch m.Text {
	case domain.CmdStart:
		// check if user exists
		if ok, err := u.userExists(m.TelegramId); err != nil {
			u.log.Error("failed to chek user existance", zap.Int64("telegram_id", m.TelegramId), zap.Error(err))
			break
			// if user not found, create a new one
		} else if !ok {
			u.outCh <- domain.Message{ChatId: m.ChatId, Text: fmt.Sprintf(domain.ReplyCreateNewUser, m.UserName)}
			ct := chat.NewChatAddUser(baseChat)
			u.chats.Store(m.ChatId, ct)
			break
		}
		// in case user exists, greet him
		u.outCh <- domain.Message{ChatId: m.ChatId, Text: fmt.Sprintf(domain.ReplyStart, m.UserName)}
		u.deleteChat(m.ChatId) // delete any existing chats, just in case

	case domain.CmdHelp:
		u.outCh <- domain.Message{ChatId: m.ChatId, Text: domain.ReplyHelp}
		u.deleteChat(m.ChatId) // delete any existing chats, just in case

	case domain.CmdUpdateUser:
		ct := chat.NewChatUpdateUser(baseChat)
		u.chats.Store(m.ChatId, ct)

	case domain.CmdAdd:
		ct := chat.NewChatAddReminder(baseChat)
		u.chats.Store(m.ChatId, ct)

	case domain.CmdList:
		ct := chat.NewChatListReminders(baseChat)
		u.chats.Store(m.ChatId, ct)

	case domain.CmdDelete:
		ct := chat.NewChatDeleteReminder(baseChat)
		u.chats.Store(m.ChatId, ct)

	case domain.CmdUpdate:
		ct := chat.NewChatUpdateReminder(baseChat)
		u.chats.Store(m.ChatId, ct)

	default:
		u.outCh <- domain.Message{UserName: "Remindista", ChatId: m.ChatId, Text: domain.ReplyUnkonwCommand}
	}
}

func (u *Updater) userExists(tgId int64) (bool, error) {
	user, err := u.db.GetUserByTelegramId(context.Background(), tgId)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	if user.TelegramId == 0 {
		return false, nil
	}

	return true, nil
}
