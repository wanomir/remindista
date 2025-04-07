package updater

import (
	"sync"

	"github.com/vedomirr/l"
	"github.com/vedomirr/remindista/internal/domain"

	"go.uber.org/zap"
)

type Updater struct {
	telegram     telegramService
	db           repository
	chats        *sync.Map
	outCh        chan domain.Message
	deleteChatCh chan int64
	log          *zap.Logger
}

func NewUpdater(telegram telegramService, db repository) *Updater {
	u := &Updater{
		telegram:     telegram,
		chats:        new(sync.Map),
		deleteChatCh: make(chan int64),
		db:           db,
		log:          l.Logger(),
	}

	go u.deleteInactiveChats()

	return u
}

func (u *Updater) deleteInactiveChats() {
	for chatId := range u.deleteChatCh {
		u.deleteChat(chatId)
	}
}

func (u *Updater) deleteChat(chatId int64) {
	u.chats.Delete(chatId)
}
