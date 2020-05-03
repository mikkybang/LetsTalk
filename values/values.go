package values

const (
	DatabaseName        = "unilagDatabase"
	AdminCollectionName = "administrators"
	UsersCollectionName = "users"
	RoomsCollectionName = "Rooms"

	AdminCookieName = "Admin"
	UserCookieName  = "User"
	// Files are saved as base64 format to database then if queried,
	TXT = iota
	INFO
	MP3
	EXE
	MP4
	WAV
	JPG
	PNG
)

var (
	RoomUsers map[string][]string
	Users     map[string]string
)
