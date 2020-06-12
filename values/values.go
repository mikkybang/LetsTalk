package values

type MessageType int

const (
	// Files are saved as base64 format to database then if queried,
	TXT MessageType = iota
	INFO
	FILE
)

// All request message types both clients and server
var (
	NewFileUploadMsgType          = "NewFileUpload"
	NewMessageMsgType             = "NewMessage"
	RequestAllMessagesMsgType     = "RequestAllMessages"
	SearchUserMsgType             = "SearchUser"
	WebsocketOpenMsgType          = "WebsocketOpen"
	NewRoomCreatedMsgType         = "NewRoomCreated"
	ExitRoomMsgType               = "ExitRoom"
	RequestUsersToJoinRoomMsgType = "RequestUsersToJoinRoom"
	UserJoinedRoomMsgType         = "UserJoinedRoom"

	UploadFileErrorMsgType = "UploadFileError"
)

var (
	// MapEmailToName maps user email to name
	MapEmailToName map[string]string
)
