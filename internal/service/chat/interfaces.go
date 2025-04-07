package chat

import (
	"context"

	r "github.com/vedomirr/remindista/internal/entity/reminder"
	u "github.com/vedomirr/remindista/internal/entity/user"
)

type repository interface {
	repoUsers
	repoReminders
}

type repoUsers interface {
	CreateUser(ctx context.Context, user u.User) (id int, err error)
	GetUser(ctx context.Context, id int) (user u.User, err error)
	GetUserByTelegramId(ctx context.Context, telegramId int64) (user u.User, err error)
	UpdateUser(ctx context.Context, user u.User) (affected int, err error)
	DeleteUser(ctx context.Context, id int, telegramId int64) (affected int, err error)
}

type repoReminders interface {
	CreateReminder(ctx context.Context, rmd r.Reminder) (id int, err error)
	GetReminder(ctx context.Context, id int) (rmd r.Reminder, err error)
	GetRemindersByUserId(ctx context.Context, userId int) (rmds []r.Reminder, err error)
	UpdateReminder(ctx context.Context, rmd r.Reminder) (affected int, err error)
	DeleteReminder(ctx context.Context, id int) (affected int, err error)
	DeleteRemindersByTag(ctx context.Context, userId int, tag string) (affected int, err error)
	DeleteRemindersByUserId(ctx context.Context, userId int) (affected int, err error)
}
