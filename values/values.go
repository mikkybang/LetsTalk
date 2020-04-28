package values

const (
	DatabaseName        = "unilagDatabase"
	AdminCollectionName = "administrators"
	UsersCollectionName = "users"
	RoomsCollectionName = "Rooms"

	AdminCookieName = "Admin"
	UserCookieName  = "User"
)

var (
	RoomUsers map[string][]string
	Users     map[string]string
)
