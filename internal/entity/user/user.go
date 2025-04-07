package user

import (
	"errors"
	"fmt"
	"log"
	"time"

	_ "time/tzdata"
)

const (
	defaultLocation    = "Europe/Moscow"
	defaultWindowFloor = "8:00"
	defaultWindowCeil  = "22:00"
)

type User struct {
	Id          int
	TelegramId  int64
	ChatId      int64
	IsRunning   bool
	Location    *time.Location
	WindowFloor time.Time
	WindowCeil  time.Time
}

type UserOption func(*User)

func NewUser(opts ...UserOption) (u User) {
	u.loadDefaultWindowFloor()
	u.loadDefaultWindowCeil()

	for _, opt := range opts {
		opt(&u)
	}

	if u.Location == nil {
		u.loadDefaultLocation()
	}

	return u
}

func WithTelegramId(telegramId int64) UserOption {
	return func(u *User) {
		u.TelegramId = telegramId
	}
}

func WithChatId(chatId int64) UserOption {
	return func(u *User) {
		u.ChatId = chatId
	}
}

func WithIsRunning(isRunning bool) UserOption {
	return func(u *User) {
		u.IsRunning = isRunning
	}
}

func (u *User) LocationString() string {
	if u.Location == nil {
		return "Unknonw Location"
	}
	return u.Location.String()
}

func (u *User) WindowFloorString() string {
	return u.WindowFloor.Format("15:04")
}

func (u *User) WindowCeilString() string {
	return u.WindowCeil.Format("15:04")
}

func (u *User) SetLocation(location string) (err error) {
	if u.Location, err = time.LoadLocation(location); err != nil {
		return err
	}

	return nil
}

func (u *User) SetWindowFloor(time string) error {
	floor, err := u.parseTime(time)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	u.WindowFloor = floor

	return nil
}

func (u *User) SetWindowCeil(time string) error {
	ceil, err := u.parseTime(time)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	if !u.WindowFloor.Before(ceil) {
		return errors.New("floor time must be before ceil")
	}

	u.WindowCeil = ceil

	return nil
}

func (u *User) Time() time.Time {
	if u.Location == nil {
		u.loadDefaultLocation()
	}

	return time.Now().In(u.Location)
}

func (u *User) FloorDuration() time.Duration {
	day := 24 * time.Hour
	floorTrunc := u.WindowFloor.Truncate(day)
	return u.WindowFloor.Sub(floorTrunc)
}

func (u *User) CeilDuration() time.Duration {
	day := 24 * time.Hour
	ceilTrunc := u.WindowCeil.Truncate(day)
	return u.WindowCeil.Sub(ceilTrunc)
}

func (u *User) loadDefaultLocation() {
	loc, err := time.LoadLocation(defaultLocation)
	if err != nil {
		log.Fatal("failed to load default location")
	}
	u.Location = loc
}

func (u *User) loadDefaultWindowFloor() {
	floor, err := u.parseTime(defaultWindowFloor)
	if err != nil {
		log.Fatal("failed to parse default floor time")
	}
	u.WindowFloor = floor
}

func (u *User) loadDefaultWindowCeil() {
	ceil, err := u.parseTime(defaultWindowCeil)
	if err != nil {
		log.Fatal("failed to parse default ceil time")
	}

	if !u.WindowFloor.Before(ceil) {
		log.Fatal("default floor time must be before ceil")
	}

	u.WindowCeil = ceil
}

func (u *User) parseTime(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}
