package model

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/metaclips/LetsTalk/values"

	"github.com/google/uuid"
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
		return values.ErrIncorrectUUID
	}
	return nil
}

func GetUser(key string, user string) (names []string) {
	names = make([]string, 0)
	for email := range values.Users {
		if email == "" || email == user {
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
		if b.UserID == user {
			userExists = true
			break
		}
	}
	if !userExists {
		return nil, values.ErrInvalidUser
	}

	messages.Messages = append(messages.Messages, b)
	_, err = db.Collection(values.RoomsCollectionName).UpdateOne(context.TODO(), bson.M{"_id": b.RoomID},
		bson.M{"$set": bson.M{"messages": messages.Messages}})

	return messages.RegisteredUsers, err
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
	chats.RegisteredUsers = append(chats.RegisteredUsers, b.Email)

	if _, err := db.Collection(values.RoomsCollectionName).InsertOne(context.TODO(), chats); err != nil {
		return "", err
	}

	user := User{Email: b.Email}
	if err := user.updateRoomsJoinedByUsers(chats.RoomID, chats.RoomName); err != nil {
		return "", err
	}

	return chats.RoomID, nil
}

func (b Joined) JoinRoom() ([]string, error) {
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.Email,
	})

	var user User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	var joinRequestLegit bool
	for i, request := range user.JoinRequest {
		if request.RoomID == b.RoomID {
			joinRequestLegit = true
			user.JoinRequest = append(user.JoinRequest[:i], user.JoinRequest[i+1:]...)
			break
		}
	}

	if !joinRequestLegit {
		return nil, values.ErrIllicitJoinRequest
	}

	user.RoomsJoined = append(user.RoomsJoined, RoomsJoined{RoomID: b.RoomID, RoomName: b.RoomName})

	_, err = db.Collection(values.UsersCollectionName).UpdateOne(context.TODO(), bson.M{"_id": b.Email},
		bson.M{"$set": bson.M{"joinRequest": user.JoinRequest, "roomsJoined": user.RoomsJoined}})
	if err != nil {
		return nil, err
	}

	result = db.Collection(values.RoomsCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.RoomID,
	})

	var messages Chats
	if err := result.Decode(&messages); err != nil {
		return nil, err
	}

	message := Message{
		Message: b.Email + " Joined",
		Type:    getContentType(values.INFO),
	}

	messages.RegisteredUsers = append(messages.RegisteredUsers, b.Email)
	messages.Messages = append(messages.Messages, message)

	_, err = db.Collection(values.RoomsCollectionName).UpdateOne(context.TODO(), bson.M{
		"_id": b.RoomID,
	}, bson.M{"$set": bson.M{"registeredUsers": messages.RegisteredUsers, "messages": messages.Messages}})

	return messages.RegisteredUsers, err
}

func (b JoinRequest) RequestUserToJoinRoom(userToJoinEmail string) ([]string, error) {
	var room Chats
	result := db.Collection(values.RoomsCollectionName).FindOne(context.TODO(), bson.M{"_id": b.RoomID})

	if err := result.Decode(&room); err != nil {
		return nil, err
	}

	// Confirm if person making the request is part of the room.
	var requesterLegit bool
	for _, registeredUser := range room.RegisteredUsers {
		if registeredUser == b.RequestingUserID {
			requesterLegit = true
			break
		} else if registeredUser == userToJoinEmail {
			return nil, values.ErrUserExistInRoom
		}
	}

	if !requesterLegit {
		return nil, errors.New("Invalid user made a RequestUsersToJoinRoom request Name: " + b.RequestingUserID)
	}

	result = db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{"_id": userToJoinEmail})
	var user User

	if err := result.Decode(&user); err != nil {
		return nil, err
	}

	// Check if user has already been requested by the room.
	for _, request := range user.JoinRequest {
		if b.RoomID == request.RoomID {
			return nil, values.ErrUserAlreadyRequested
		}
	}
	user.JoinRequest = append(user.JoinRequest, b)

	_, err := db.Collection(values.UsersCollectionName).UpdateOne(context.TODO(), bson.M{"_id": userToJoinEmail},
		bson.M{"$set": bson.M{"joinRequest": user.JoinRequest}})

	if err != nil {
		return nil, err
	}

	message := Message{
		Message: fmt.Sprintf("%s was requested to join the room by %s", userToJoinEmail, b.RequestingUserID),
		Type:    getContentType(values.INFO),
	}
	room.Messages = append(room.Messages, message)

	_, err = db.Collection(values.RoomsCollectionName).UpdateOne(context.TODO(), bson.M{"_id": b.RoomID},
		bson.M{"$set": bson.M{"messages": room.Messages}})
	return room.RegisteredUsers, err
}

func (b User) AddUserToRoom(roomID, roomName string) error {
	b.updateRoomsJoinedByUsers(roomID, roomName)
	var chats Chats
	message := Message{
		Message: b.Email + " Joined",
		Type:    getContentType(values.INFO),
	}

	result := db.Collection(values.RoomsCollectionName).FindOne(context.TODO(), bson.M{
		"_id": roomID,
	})

	if err := result.Decode(&chats); err != nil {
		return err
	}

	chats.Messages = append(chats.Messages, message)
	_, err := db.Collection(values.RoomsCollectionName).UpdateOne(context.TODO(), bson.M{"_id": roomID},
		bson.M{"$set": bson.M{"messages": chats.Messages}})

	return err
}

func (b *User) updateRoomsJoinedByUsers(roomID, roomName string) error {
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.Email,
	})

	if err := result.Decode(&b); err != nil {
		return err
	}

	var roomJoined = RoomsJoined{RoomID: roomID, RoomName: roomName}
	b.RoomsJoined = append(b.RoomsJoined, roomJoined)

	_, err := db.Collection(values.UsersCollectionName).UpdateOne(context.TODO(), bson.M{"_id": b.Email},
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

func GetAllUserRooms(email string) (User, error) {
	var user User
	result := db.Collection(values.UsersCollectionName).FindOne(context.TODO(), bson.M{
		"_id": email,
	})

	if err := result.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func getContentType(contentType values.MessageType) string {
	switch contentType {
	case values.INFO:
		return "info"
	case values.TXT:
		return "txt"
	}

	return ""
}
