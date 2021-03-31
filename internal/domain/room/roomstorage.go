package room

type Room struct {
	Id      string
	Creator string
	Members []string
}

type Interface interface {
	CreateRoom(creatorId string) (Room, error)
	GetRoomById(actorId, roomId string) (Room, error)
	UpdateRoom(actorId, roomId string, upd UpdateFunc) (Room, error)
	ListRooms(accountId string) ([]Room, error)
}

type UpdateFunc func(r Room) (Room, error)
