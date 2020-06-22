package values

type MessageType int

const (
	// Files are saved as base64 format to database then if queried,
	TXT MessageType = iota
	INFO
	FILE
)

// All request message types both clients and server
const (
	NewFileUploadMsgType          = "NewFileUpload"
	NewMessageMsgType             = "NewMessage"
	RequestAllMessagesMsgType     = "RequestAllMessages"
	SearchUserMsgType             = "SearchUser"
	WebsocketOpenMsgType          = "WebsocketOpen"
	NewRoomCreatedMsgType         = "NewRoomCreated"
	ExitRoomMsgType               = "ExitRoom"
	RequestUsersToJoinRoomMsgType = "RequestUsersToJoinRoom"
	UserJoinedRoomMsgType         = "UserJoinedRoom"

	UploadFileErrorMsgType   = "UploadFileError" // UploadFileErrorMsgType is sent to client only.
	UploadFileSuccessMsgType = "FileUploadSuccess"
	UploadFileChunkMsgType   = "UploadFileChunk"

	RequestDownloadMsgType     = "RequestDownload"
	DownloadFileChunkMsgType   = "DownloadFileChunk"
	DownloadFileErrorMsgType   = "DownloadFileError"   // DownloadFileErrorMsgType is sent to client only.
	DownloadFileSuccessMsgType = "DownloadFileSuccess" // DownloadFileSuccessMsgType is sent to client only.
)

var (
	// MapEmailToName maps user email to name
	MapEmailToName map[string]string
)
