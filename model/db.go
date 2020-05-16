package model

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/metaclips/LetsTalk/values"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db *mongo.Database
)

func InitDB() {
	dbHost := os.Getenv("db_host")
	mongoDB, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbHost))
	if err != nil {
		log.Fatalln(err)
	}

	db = mongoDB.Database(values.DatabaseName)

	// Ping mongo database continuosly if up.
	go func(mongoDB *mongo.Client) {
		for {
			if err := mongoDB.Ping(context.TODO(), readpref.Primary()); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 5)
		}
	}(mongoDB)

	values.RoomUsers = make(map[string][]string)
	values.Users = make(map[string]string)

	getContent := func(collection string, content interface{}) {
		result, err := db.Collection(collection).Find(context.TODO(), bson.D{})
		if err != nil {
			log.Fatalln("error while getting collection", err)
		}

		err = result.All(context.TODO(), content)
		if err != nil {
			log.Fatalln("error getting collection results", err)
		}
	}

	var roomChats []Chats
	var users []User

	getContent(values.RoomsCollectionName, &roomChats)
	getContent(values.UsersCollectionName, &users)

	for _, chat := range roomChats {
		values.RoomUsers[chat.RoomID] = chat.RegisteredUsers
	}

	for _, user := range users {
		values.Users[user.Email] = user.Name
	}
}
