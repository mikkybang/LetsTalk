package values

type MessageType int

const (
	DatabaseName        = "unilagDatabase"
	AdminCollectionName = "administrators"
	UsersCollectionName = "users"
	RoomsCollectionName = "Rooms"

	AdminCookieName = "Admin"
	UserCookieName  = "User"
	TimeLayout      = "Monday, 02-Jan-06 15:04:05"
	// Files are saved as base64 format to database then if queried,
	TXT MessageType = iota
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
