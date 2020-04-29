package model

import (
	"context"
	"errors"
	"net/http"
	"strings"

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
