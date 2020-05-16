package values

type MessageType int

const (
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
