package messagerepo

import (
	"github.com/mp-hl-2021/chat/internal/domain/message"

	"strconv"
	"sync"
	"time"
)

type Memory struct {
	messagesByRoom map[string][]message.Message
	nextId         uint64
	mu             *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		messagesByRoom: make(map[string][]message.Message),
		mu:             &sync.Mutex{},
	}
}

func (m *Memory) CreateMessage(creatorId, roomId string, text string, createdAt time.Time) (message.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg := message.Message{
		Id:        strconv.FormatUint(m.nextId, 16),
		Author:    creatorId,
		Room:      roomId,
		CreatedAt: createdAt,
		Text:      text,
	}
	msgs, ok := m.messagesByRoom[roomId]
	if !ok {
		m.messagesByRoom[roomId] = make([]message.Message, 0, 1)
	}
	m.messagesByRoom[roomId] = append(msgs, msg)
	m.nextId++
	return msg, nil
}

func (m *Memory) ListMessages(actorId, roomId string) ([]message.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msgs, ok := m.messagesByRoom[roomId]
	if !ok {
		return []message.Message{}, nil
	}
	return msgs, nil
}
