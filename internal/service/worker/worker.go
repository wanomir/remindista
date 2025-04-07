package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/vedomirr/l"
	"github.com/vedomirr/remindista/internal/domain"
	r "github.com/vedomirr/remindista/internal/entity/reminder"
	u "github.com/vedomirr/remindista/internal/entity/user"

	"go.uber.org/zap"
)

type telegramService interface {
	ReceiveMessages(ctx context.Context) chan domain.Message
	SendMessage(chatID int64, html string, keyboard domain.Keyboard) error
	SendMessageMarkdownV2(chatID int64, html string, keyboard domain.Keyboard) error
}

type repository interface {
	GetAllUsers(ctx context.Context, limit int, offset int) (users []u.User, err error)
	GetRemindersByUserIdAndTime(ctx context.Context, userId int, userTime time.Time) (rmds []r.Reminder, err error)
	UpdateReminder(ctx context.Context, rmd r.Reminder) (affectd int, err error)
}

type Worker struct {
	telegram telegramService
	db       repository
	log      *zap.Logger
	interval time.Duration
}

func NewWorker(telegram telegramService, repo repository, workInterval time.Duration) *Worker {
	return &Worker{
		telegram: telegram,
		db:       repo,
		interval: workInterval,
		log:      l.Logger(),
	}
}

func (w *Worker) Run(ctx context.Context) {
	t := time.NewTicker(w.interval)

	for {
		select {
		case <-t.C:
			if err := w.processUsers(); err != nil {
				w.log.Error("users pagination error", zap.Error(err))
			}

		case <-ctx.Done():
			w.log.Info("shutting down worker service")
			return
		}
	}
}

func (w *Worker) Stop() {}

func (w *Worker) processUsers() (err error) {
	routinesLimit := make(chan struct{}, 100) // limit the number of goroutines

	users, err := w.db.GetAllUsers(context.Background(), 100, 0)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		routinesLimit <- struct{}{}
		go w.processUser(user, routinesLimit)
	}

	return nil
}

func (w *Worker) processUser(user u.User, limit chan struct{}) {
	defer func() { <-limit }()

	rmds, err := w.db.GetRemindersByUserIdAndTime(context.Background(), user.Id, user.Time())
	if err != nil {
		w.log.Error("error getting reminders for user", zap.Int("user id", user.Id), zap.Error(err))
		return
	}

	for _, rmd := range rmds {
		limit <- struct{}{}
		go w.processReminder(rmd, user, limit)
	}
}

func (w *Worker) processReminder(rmd r.Reminder, user u.User, limit chan struct{}) {
	defer func() { <-limit }()

	if err := w.telegram.SendMessageMarkdownV2(user.ChatId, rmd.StringMdV2(), rmd.Keyboard()); err != nil {
		w.log.Error("failed to send message", zap.Int64("chat id", user.ChatId), zap.Error(err))
	}

	rmd.UpdateNextReminder(user.Time(), user.FloorDuration(), user.CeilDuration())
	w.log.Info("reminder updated", zap.Int("reminder_id", rmd.Id), zap.String("next_reminder", rmd.NextReminderString()))

	if _, err := w.db.UpdateReminder(context.Background(), rmd); err != nil {
		w.log.Error("failed to update reminder", zap.Int("reminder id", rmd.Id), zap.Error(err))
	}
}
