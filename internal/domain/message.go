package domain

type Item struct{ Key, Val string }

type Keyboard [][]Item

type Message struct {
	ChatId     int64
	TelegramId int64
	UserName   string
	Text       string
	Keyboard
}
