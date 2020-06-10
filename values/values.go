package values

type MessageType int

const (
	// Files are saved as base64 format to database then if queried,
	TXT MessageType = iota
	INFO
	FILE
)

var (
	// MapEmailToName maps user email to name
	MapEmailToName map[string]string
)
