package room

import (
	"github.com/mp-hl-2021/chat/internal/domain"
	"github.com/mp-hl-2021/chat/internal/domain/room"
)

type Room struct {
	Id      string
	Members []string // account ids or some structures later
}

type Interface interface {
	CreateRoom(creatorId string) (Room, error)
	ListRooms(accountId string) ([]Room, error) // todo

	GetRoomById(actorId, roomId string) (Room, error)
	AddMembers(actorId, roomId string, members []string) error
	RemoveMembers(actorId, roomId string, members []string) error
}

type UseCases struct {
	RoomStorage room.Interface
}

func (u *UseCases) CreateRoom(creatorId string) (Room, error) {
	r, err := u.RoomStorage.CreateRoom(creatorId)
	if err != nil {
		return Room{}, err
	}
	return Room{Id: r.Id, Members: r.Members}, nil
}

func (u *UseCases) ListRooms(accountId string) ([]Room, error) {
	rr, err := u.RoomStorage.ListRooms(accountId)
	if err != nil {
		return nil, err
	}
	res := make([]Room, 0, len(rr))
	for _, r := range rr {
		res = append(res, Room{Id: r.Id, Members: r.Members})
	}
	return res, nil
}

func (u *UseCases) GetRoomById(actorId, roomId string) (Room, error) {
	r, err := u.RoomStorage.GetRoomById(actorId, roomId)
	if err != nil {
		return Room{}, err
	}
	for _, m := range r.Members {
		if m == actorId {
			return Room{Id: r.Id, Members: r.Members}, nil
		}
	}
	return Room{}, domain.ErrNotFound // todo: may be "unauthorized"?
}

func (u *UseCases) AddMembers(actorId, roomId string, members []string) error {
	_, err := u.RoomStorage.UpdateRoom(actorId, roomId, func(r room.Room) (room.Room, error) {
		authorized := authorize(actorId, r.Members)
		if !authorized {
			return r, domain.ErrNotFound // todo: return "unauthorized"
		}
		// note: this is not a beautiful way to insert new members
		// may be repository should provide some order for room members.
		newMembers := make([]string, 0, len(members))
		for i := 0; i < len(members); i++ {
			newM := true
			for j := 0; j < len(r.Members); j++ {
				if members[i] == r.Members[j] {
					newM = false
				}
			}
			if newM {
				newMembers = append(newMembers, members[i])
			}
		}
		r.Members = append(r.Members, newMembers...)
		return r, nil
	})
	return err
}

func (u *UseCases) RemoveMembers(actorId, roomId string, members []string) error {
	_, err := u.RoomStorage.UpdateRoom(actorId, roomId, func(r room.Room) (room.Room, error) {
		authorized := authorize(actorId, r.Members)
		if !authorized {
			return r, domain.ErrNotFound // todo: return "unauthorized"
		}
		// note: see AddMembers note
		for i := 0; i < len(members); i++ {
			for j := 0; j < len(r.Members); j++ {
				if members[i] == r.Members[j] {
					r.Members[j] = r.Members[len(r.Members)-1]
					r.Members = r.Members[:len(r.Members)-1]
					break
				}
			}
		}
		return r, nil
	})
	return err
}

func authorize(actorId string, members []string) bool {
	for _, m := range members {
		if m == actorId {
			return true
		}
	}
	return false
}
