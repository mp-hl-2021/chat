package roomrepo

import (
	"github.com/mp-hl-2021/chat/internal/domain"
	"github.com/mp-hl-2021/chat/internal/domain/room"

	"strconv"
	"sync"
)

type Memory struct {
	roomById         map[string]room.Room
	roomsByAccountId map[string]map[string]room.Room
	nextId           uint64
	mu               *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		roomById:         make(map[string]room.Room),
		roomsByAccountId: make(map[string]map[string]room.Room),
		mu:               &sync.Mutex{},
	}
}

func (m *Memory) CreateRoom(creatorId string) (room.Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r := room.Room{
		Id:      strconv.FormatUint(m.nextId, 16),
		Creator: creatorId,
		Members: []string{creatorId},
	}
	m.roomById[r.Id] = r
	accountRooms, ok := m.roomsByAccountId[creatorId]
	if !ok {
		accountRooms = make(map[string]room.Room)
	}
	accountRooms[r.Id] = r
	m.roomsByAccountId[creatorId] = accountRooms
	m.roomById[r.Id] = r
	m.nextId++
	return r, nil
}

func (m *Memory) GetRoomById(actorId, roomId string) (room.Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.roomById[roomId]
	if !ok {
		return r, domain.ErrNotFound
	}
	for _, member := range r.Members {
		if member == actorId {
			return r, nil
		}
	}
	return r, domain.ErrUnauthorized
}

func (m *Memory) UpdateRoom(actorId, roomId string, upd room.UpdateFunc) (room.Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, err := m.GetRoomById(actorId, roomId)
	if err != nil {
		return r, err
	}
	r, err = upd(r)
	if err != nil {
		return r, err
	}
	m.roomById[roomId] = r
	m.roomsByAccountId[actorId][roomId] = r
	return r, nil
}

func (m *Memory) ListRooms(accountId string) ([]room.Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rm, ok := m.roomsByAccountId[accountId]
	if !ok {
		return nil, domain.ErrNotFound
	}
	rr := make([]room.Room, 0, len(rm))
	for _, val := range rm {
		rr = append(rr, val)
	}
	return rr, nil
}
