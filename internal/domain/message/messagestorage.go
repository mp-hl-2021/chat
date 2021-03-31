package message

import "time"

type Message struct {
	Id        string
	Author    string
	Room      string
	CreatedAt time.Time

	Text string
}

type Interface interface {
	CreateMessage(creatorId, roomId string, text string, createdAt time.Time) (Message, error)
	ListMessages(actorId, roomId string) ([]Message, error)
}
