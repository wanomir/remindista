package updater

import (
	"errors"
	"regexp"

	"github.com/vedomirr/remindista/internal/domain"
)

func (u *Updater) ProcessMessage(m domain.Message) error {
	if u.isValidCmd(m.Text) {
		u.processCmd(m)
		return nil
	}

	if u.isValidCallback(m.Text) {
		u.processCallback(m)
		return nil
	}

	if chat, ok := u.chats.Load(m.ChatId); ok {
		switch v := chat.(type) {
		case chattable:
			v.PassInput(m.Text)
		default:
			return errors.New("error casting to chat interface")
		}
	}

	return nil
}

func (u *Updater) isValidCmd(s string) bool {
	re := regexp.MustCompile(`^\/[a-z_]+$`)
	return re.MatchString(s)
}

func (u *Updater) isValidCallback(s string) bool {
	re := regexp.MustCompile(`^:[a-z_]+ [0-9]+$`)
	return re.MatchString(s)
}
