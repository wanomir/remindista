package updater

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vedomirr/remindista/internal/domain"
	"github.com/vedomirr/remindista/internal/service/chat"

	"go.uber.org/zap"
)

func (u *Updater) processCallback(m domain.Message) {
	callback, rmdId, err := u.parseCallbackParams(m.Text)
	if err != nil {
		u.log.Error("failed to parse callback params", zap.Error(err))
		return
	}

	baseChat := chat.NewChat(m.ChatId, m.TelegramId, u.outCh, u.deleteChatCh, u.db)
	switch callback {
	case domain.CallbackDelete:
		ct := chat.NewChatDeleteReminderById(baseChat, rmdId)
		u.chats.Store(m.ChatId, ct)

	case domain.CallbackUpdate:
		ct := chat.NewChatUpdateReminderById(baseChat, rmdId)
		u.chats.Store(m.ChatId, ct)

	case domain.CallbackIncreaseFrequency:
		u.increaseFrequency(m.TelegramId, m.ChatId, rmdId)

	case domain.CallbackDecreaseFrequency:
		u.decreaseFrequency(m.TelegramId, m.ChatId, rmdId)

	default:
		u.log.Error("unknown callback", zap.String("callback", callback))
	}
}

func (u *Updater) parseCallbackParams(s string) (string, int, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return "", 0, domain.ErrorInvalidCallback
	}

	rmdId, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", domain.ErrorInvalidCallback, err)
	}

	return parts[0], rmdId, nil
}

func (u *Updater) increaseFrequency(tgId, chatId int64, rmdId int) {
	user, err := u.db.GetUserByTelegramId(context.Background(), tgId)
	if err != nil {
		u.log.Error("failed to get user", zap.Int64("telegram_id", tgId), zap.Error(err))
		return
	}

	rmd, err := u.db.GetReminder(context.Background(), rmdId)
	if err != nil {
		u.log.Error("failed to get reminder", zap.Int("reminder_id", rmdId), zap.Error(err))
		return
	}

	if rmd.Frequency <= time.Minute {
		u.outCh <- domain.Message{ChatId: chatId, Text: domain.ReplyMinimumFrequency}
		return
	}

	rmd.Frequency = max(rmd.Frequency/2, time.Minute)
	rmd.UpdateNextReminder(user.Time(), user.FloorDuration(), user.CeilDuration())

	if _, err := u.db.UpdateReminder(context.Background(), rmd); err != nil {
		u.log.Error("failed to update reminder", zap.Int("reminder_id", rmdId), zap.Error(err))
	}

	u.outCh <- domain.Message{ChatId: chatId, Text: fmt.Sprintf(domain.ReplyFrequencyUpdated, rmd.FreqeuncyString())}
}

func (u *Updater) decreaseFrequency(tgId, chatId int64, rmdId int) {
	user, err := u.db.GetUserByTelegramId(context.Background(), tgId)
	if err != nil {
		u.log.Error("failed to get user", zap.Int64("telegram_id", tgId), zap.Error(err))
		return
	}

	rmd, err := u.db.GetReminder(context.Background(), rmdId)
	if err != nil {
		u.log.Error("failed to get reminder", zap.Int("reminder_id", rmdId), zap.Error(err))
		return
	}

	if rmd.Frequency >= time.Hour*24*365 {
		u.outCh <- domain.Message{ChatId: chatId, Text: domain.ReplyMaximumFrequency}
		return
	}

	rmd.Frequency = min(rmd.Frequency*2, time.Hour*24*365)
	rmd.UpdateNextReminder(user.Time(), user.FloorDuration(), user.CeilDuration())

	if _, err := u.db.UpdateReminder(context.Background(), rmd); err != nil {
		u.log.Error("failed to update reminder", zap.Int("reminder_id", rmdId), zap.Error(err))
	}

	u.outCh <- domain.Message{ChatId: chatId, Text: fmt.Sprintf(domain.ReplyFrequencyUpdated, rmd.FreqeuncyString())}
}
