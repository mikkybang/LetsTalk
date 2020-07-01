package model

import (
	"time"

	"github.com/gorilla/websocket"
)

type Configuration struct {
	DbHost string
	Port   string
}

type CookieDetail struct {
	Email      string
	Collection string
	CookieName string
	Path       string
	Data       CookieData
}

type CookieData struct {
	ExitTime time.Time
	UUID     string
	Email    string
	Super    bool
}

type User struct {
	Email string `bson:"_id" json:"email"`
	Name  string `bson:"name" json:"name"`
	DOB   string `bson:"age" json:"age"`
	Class string `bson:"class" json:"class"`
	// ID should either users matric or email stripping @....
	ID           string        `bson:"userID" json:"userID"`
	ParentEmail  string        `bson:"parentEmail" json:"parentEmail"`
	ParentNumber string        `bson:"parentNumber" json:"parentNumber"`
	Password     []byte        `bson:"password" json:"password"`
	Faculty      string        `bson:"faculty" json:"faculty"`
	UUID         string        `bson:"loginUUID" json:"uuid"`
	RoomsJoined  []RoomsJoined `bson:"roomsJoined" json:"roomsJoined"`
	JoinRequest  []JoinRequest `bson:"joinRequest" json:"joinRequest"`
}

type RoomsJoined struct {
	RoomID   string `bson:"rooomID" json:"roomID"`
	RoomName string `bson:"rooomName" json:"roomName"`
}

type JoinRequest struct {
	RoomID             string   `bson:"_id" json:"roomID"`
	RoomName           string   `bson:"roomName" json:"roomName"`
	RequestingUserName string   `bson:"requestingUserName" json:"requestingUserName"`
	RequestingUserID   string   `bson:"requestingUserID" json:"requestingUserID"`
	Users              []string `bson:"-" json:"users"`
}

type Admin struct {
	StaffDetails User `bson:",inline"`
	Super        bool `bson:"super" json:"super"`
}

type Room struct {
	RoomID          string    `bson:"_id" json:"email"`
	RoomName        string    `bson:"roomName" json:"roomName"`
	RegisteredUsers []string  `bson:"registeredUsers"`
	Messages        []Message `bson:"messages" json:"messages"`
}

type Message struct {
	RoomID      string `bson:"-" json:"roomID"`
	Message     string `bson:"message" json:"message"`
	FileSize    string `bson:"fileSize,omitempty" json:"fileSize,omitempty"`
	FileHash    string `bson:"fileHash,omitempty" json:"fileHash,omitempty"`
	UserID      string `bson:"userID" json:"userID"`
	Name        string `bson:"name" json:"name"`
	Index       int    `bson:"index" json:"index"`
	Time        string `bson:"time" json:"time"`
	Type        string `bson:"type" json:"type"`
	MessageType string `bson:"-" json:"msgType"`
}

type Joined struct {
	RoomID      string `json:"roomID"`
	RoomName    string `json:"roomName"`
	Email       string `json:"email"`
	Joined      bool   `json:"joined"`
	MessageType string `bson:"-" json:"msgType"`
}

type NewRoomRequest struct {
	Email       string `json:"email"`
	RoomName    string `json:"roomName"`
	MessageType string `bson:"-" json:"msgType"`
}

// File save files making sure they are distinct
type File struct {
	MsgType        string `bson:"-" json:"msgType,omitempty"`
	UniqueFileHash string `bson:"_id" json:"fileHash"`
	FileName       string `bson:"fileName" json:"fileName"`
	User           string `bson:"userID" json:"userID"`
	FileSize       string `bson:"fileSize" json:"fileSize"`
	FileType       string `bson:"fileType" json:"fileType"`
	Chunks         int    `bson:"chunks,omitempty" json:"chunks"`
}

type FileChunks struct {
	MsgType            string `bson:"-" json:"msgType,omitempty"`
	FileName           string `bson:"-" json:"fileName,omitempty"`
	UniqueFileHash     string `bson:"_id" json:"fileHash,omitempty"`
	CompressedFileHash string `bson:"compressedFileHash" json:"compressedFileHash,omitempty"`
	FileBinary         string `bson:"fileChunk" json:"fileChunk,omitempty"`
	ChunkIndex         int    `bson:"chunkIndex" json:"chunkIndex"`
}

type WSMessage struct {
	Data []byte
	User string
}

type Subscription struct {
	Conn *Connection
	User string
}

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	WS *websocket.Conn

	Send chan []byte
}

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Users map[string]map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan WSMessage

	// Register requests from the connections.
	Register chan Subscription

	// Unregister requests from connections.
	UnRegister chan Subscription
}
