package telegram

import (
	"context"
	"regexp"

	"github.com/vedomirr/l"
	"github.com/vedomirr/remindista/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
	log *zap.Logger
}

func NewTelegram(token string) (t *Telegram, err error) {
	t = &Telegram{bot: nil, log: l.Logger()}
	if t.bot, err = tgbotapi.NewBotAPI(token); err != nil {
		return nil, err
	}

	t.log.Info("authorized telegram service", zap.String("account", t.bot.Self.UserName))
	return t, nil
}

func (t *Telegram) ReceiveMessages(ctx context.Context) chan domain.Message {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)
	messages := make(chan domain.Message)

	go func(updates tgbotapi.UpdatesChannel, messages chan domain.Message, ctx context.Context) {
		defer t.bot.StopReceivingUpdates()

		for {
			select {
			case update := <-updates:
				if update.Message != nil {
					messages <- mapMessage(update.Message)
				} else if update.CallbackQuery != nil {
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
					if _, err := t.bot.Request(callback); err != nil {
						t.log.Error("failed to send callback", zap.Error(err))
					}

					messages <- mapCallback(update.CallbackQuery)
				}

			case <-ctx.Done():
				t.log.Info("telegram shutdown")
				return
			}
		}
	}(updates, messages, ctx)

	return messages
}

func (t *Telegram) SendMessage(chatID int64, text string, keyboard domain.Keyboard) error {
	msg := tgbotapi.NewMessage(chatID, text)

	if keyboard != nil {
		msg.ReplyMarkup = inlineKeyboard(keyboard)
	}

	if _, err := t.bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func (t *Telegram) SendMessageMarkdownV2(chatID int64, html string, keyboard domain.Keyboard) error {
	msg := tgbotapi.NewMessage(chatID, html)

	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if keyboard != nil {
		msg.ReplyMarkup = inlineKeyboard(keyboard)
	}

	if _, err := t.bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func (t *Telegram) SendMessageHTML(chatID int64, html string, keyboard domain.Keyboard) error {
	msg := tgbotapi.NewMessage(chatID, html)

	msg.ParseMode = tgbotapi.ModeHTML

	if keyboard != nil {
		msg.ReplyMarkup = inlineKeyboard(keyboard)
	}

	if _, err := t.bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func mapMessage(m *tgbotapi.Message) domain.Message {
	return domain.Message{
		ChatId:     m.Chat.ID,
		TelegramId: m.From.ID,
		UserName:   m.From.UserName,
		Text:       m.Text,
	}
}

func mapCallback(c *tgbotapi.CallbackQuery) domain.Message {
	return domain.Message{
		ChatId:     c.Message.Chat.ID,
		TelegramId: c.From.ID,
		UserName:   c.From.UserName,
		Text:       c.Data,
	}
}

func inlineKeyboard(keyboardValues domain.Keyboard) (keyboard tgbotapi.InlineKeyboardMarkup) {
	for _, row := range keyboardValues {
		keyboardRow := tgbotapi.NewInlineKeyboardRow()
		for _, value := range row {
			if isLink(value.Val) {
				keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonURL(value.Key, value.Val))
				continue
			}
			keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(value.Key, value.Val))
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, keyboardRow)
	}
	return keyboard
}

func isLink(s string) bool {
	re := regexp.MustCompile(`^https?:\/\/.+$`)
	return re.MatchString(s)
}
