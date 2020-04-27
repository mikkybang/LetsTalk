package model

// Files are saved as base64 format to database then if queried,
const (
	TXT = iota
	MP3
	EXE
	MP4
	WAV
	JPG
	PNG
)

type CookieDetail struct {
	Email      string
	Collection string
	CookieName string
	Path       string
	Data       map[string]interface{}
}

type User struct {
	Email        string `bson:"_id" json:"email"`
	Name         string `bson:"name" json:"name"`
	DOB          string `bson:"age" json:"age"`
	Class        string `bson:"class" json:"class"`
	ParentEmail  string `bson:"parentEmail" json:"parentEmail"`
	ParentNumber string `bson:"parentNumber" json:"parentNumber"`
	Password     []byte `bson:"password" json:"password"`
	Faculty      string `bson:"faculty" json:"faculty"`
	UUID         string `bson:"uuid" json:"uuid"`
}

type Admin struct {
	StaffDetails User `bson:",inline"`
	Super        bool `bson:"super" json:"super"`
}

type Chats struct {
	UserEmail string  `bson:"_id" json:"email"`
	Posts     []posts `bson:"posts" json:"posts"`
}

type posts struct {
	UserEmail string     `bson:"_id" json:"email"`
	Messages  []messages `bson:"messages" json:"messages"`
}

type messages struct {
	Post  string `bson:"post" json:"post"`
	Index int    `bson:"index" json:"index"`
	Type  int    `bson:"type" json:"type"`
}

// FileType save files separately and make sure they are distinct
type FileType struct {
	Downloaded bool   `bson:"downloaded" json:"downloaded"`
	Sha256     string `bson:"_id" json:"sha256"`
}
