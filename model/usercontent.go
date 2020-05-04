package model

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/metaclips/FinalYearProject/values"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func (b User) CreateUserLogin(password string, w http.ResponseWriter) error {
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.Email,
	})
	err := result.Decode(&b)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(b.Password, []byte(password)); err != nil {
		return err
	}

	err = CookieDetail{
		Email:      b.Email,
		Collection: values.UsersCollectionName,
		CookieName: values.UserCookieName,
		Path:       "/",
		Data: map[string]interface{}{
			"Email": b.Email,
		}}.CreateCookie(w)

	return err
}

func (b User) ValidateUser(email, uniqueID string) error {
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": email,
	})
	err := result.Decode(&b)
	if err != nil {
		return err
	}

	if b.UUID != uniqueID {
		return errors.New("Incorrect UUID")
	}
	return nil
}

func GetUser(key string) (names []string) {
	names = make([]string, 0)
	for email := range values.Users {
		if email == "" {
			continue
		}
		if strings.Contains(email, key) {
			names = append(names, email)
		}
	}

	return
}

func (b Message) SaveMessageContent() ([]string, error) {
	var messages Chats
	result := db.Collection(values.RoomsCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.RoomID,
	})
	err := result.Decode(&messages)
	if err != nil {
		return nil, err
	}
	var userExists bool
	// todo: checking all users really isn't required.
	// Check if user is registered to the room
	for _, user := range messages.RegisteredUsers {
		if b.User == user {
			userExists = true
		}
	}
	if !userExists {
		return nil, errors.New("Invalid user")
	}

	messages.Messages = append(messages.Messages, b)
	_, err = db.Collection(values.RoomsCollectionName).UpdateOne(context.TODO(), bson.M{
		"_id": b.RoomID,
	}, messages)

	return messages.RegisteredUsers, err
}

func (b Joined) JoinOrExitRoom() ([]string, error) {
	var messages Chats
	var broadcastToUsers []string
	result := db.Collection(values.RoomsCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.RoomID,
	})
	if err := result.Err(); err != nil {
		return nil, err
	}

	if err := result.Decode(&messages); err != nil {
		return nil, err
	}
	broadcastToUsers = messages.RegisteredUsers
	if b.Joined {
		messages.RegisteredUsers = append(messages.RegisteredUsers, b.Email)
	} else {
		users := make([]string, 0)
		for _, user := range messages.RegisteredUsers {
			if user == b.Email {
				continue
			}
			users = append(users, user)
		}
		messages.RegisteredUsers = users
	}

	_, err := db.Collection(values.RoomsCollectionName).UpdateOne(context.TODO(), bson.M{
		"_id": b.RoomID,
	}, messages)

	return broadcastToUsers, err
}

func (b NewRoomRequest) CreateNewRoom() (string, error) {
	var chats Chats
	message := Message{
		Message: b.Email + " Joined",
		Type:    getContentType(values.INFO),
	}

	chats.Messages = append(chats.Messages, message)
	chats.RoomID = uuid.New().String()
	chats.RoomName = b.RoomName
	_, err := db.Collection(values.RoomsCollectionName).InsertOne(context.TODO(), chats)
	if err != nil {
		return "", err
	}
	user := User{Email: b.Email}
	if err := user.AddUserToRoom(chats.RoomID, chats.RoomName); err != nil {
		return "", err
	}

	return chats.RoomID, nil
}

func (b User) AddUserToRoom(roomID, roomName string) error {
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.Email,
	})
	if err := result.Decode(&b); err != nil {
		return err
	}

	var roomJoined = RoomsJoined{RoomID: roomID, RoomName: roomName}
	b.RoomsJoined = append(b.RoomsJoined, roomJoined)

	_, err := db.Collection(values.UsersCollectionName).UpdateOne(context.TODO(), map[string]interface{}{"_id": b.Email},
		bson.M{"$set": bson.M{"roomsJoined": b.RoomsJoined}})

	return err
}

func GetAllMessageInRoom(roomID string) ([]Message, string, error) {
	result := db.Collection(values.RoomsCollectionName).FindOne(context.TODO(), bson.M{"_id": roomID})

	var chat Chats
	if err := result.Decode(&chat); err != nil {
		return nil, "", err
	}
	return chat.Messages, chat.RoomName, nil
}

func GetAllUserRooms(email string) ([]RoomsJoined, error) {
	var user User
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": email,
	})

	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	return user.RoomsJoined, nil
}

func getContentType(contentType int) string {
	switch contentType {
	case values.INFO:
		return "info"
	case values.TXT:
		return "txt"
	}

	return ""
}
