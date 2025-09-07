package rooms

type Repository interface {
	GetLiveRooms()
	EndRoom()
}

type repository struct {
}
