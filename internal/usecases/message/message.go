package message

import (
	"github.com/mp-hl-2021/chat/internal/domain/message"

	"time"
)

type Message struct {
	Text      string
	Author    string // account id
	Room      string // room id
	CreatedAt time.Time
}

type Interface interface {
	CreateMessage(creatorId, roomId string, text string) error
	ListMessages(actorId, roomId string) ([]Message, error)
}

type UseCases struct {
	MessageStorage message.Interface
}

func (u *UseCases) CreateMessage(creatorId, roomId string, text string) error {
	t := time.Now()
	// todo: check whether room exists
	_, err := u.MessageStorage.CreateMessage(creatorId, roomId, text, t)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCases) ListMessages(actorId, roomId string) ([]Message, error) {
	mm, err := u.MessageStorage.ListMessages(actorId, roomId)
	if err != nil {
		return nil, err
	}
	res := make([]Message, 0, len(mm))
	for _, m := range mm {
		res = append(res, Message{
			Text:      m.Text,
			Author:    m.Author,
			Room:      m.Room,
			CreatedAt: m.CreatedAt,
		})
	}
	return res, nil
}
