package model

import (
	"github.com/gorilla/websocket"
)

type CookieDetail struct {
	Email      string
	Collection string
	CookieName string
	Path       string
	Data       map[string]interface{}
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

// FileType save files separately and make sure they are distinct
type FileType struct {
	Downloaded bool   `bson:"downloaded" json:"downloaded"`
	Sha256     string `bson:"_id" json:"sha256"`
}

type WSMessage struct {
	Data []byte
	User string
}

type Subscription struct {
	Conn *Connection
	User string
}

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	WS *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// hub maintains the set of active connections and broadcasts messages to the
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
