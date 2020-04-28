package model

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/metaclips/FinalYearProject/values"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db          *mongo.Database
	UUID        uuid.UUID
	defaultCost = 10
)

func InitDB() {
	os.Setenv("db_host", "mongodb://localhost:27017")
	dbHost := os.Getenv("db_host")
	mongoDB, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbHost))
	if err != nil {
		log.Fatalln(err)
	}

	db = mongoDB.Database(values.DatabaseName)

	// Ping mongo database if up
	go func(mongoDB *mongo.Client) {
		for {
			if err := mongoDB.Ping(context.TODO(), readpref.Primary()); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 5)
		}
	}(mongoDB)

	UUID, err = uuid.NewUUID()
	if err != nil {
		log.Fatalln("could not initiate uuid, err: ", err)
	}

	values.RoomUsers = make(map[string][]string)
	values.Users = make(map[string]string)

	// Why I love Generics :(
	result, err := db.Collection(values.RoomsCollectionName).Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalln("error while getting all room names ", err)
	}
	var roomChats []Chats
	err = result.Decode(&roomChats)
	// todo: since nothing has been added to the database....
	// if err != nil {
	// 	log.Fatalln("error converting room users interface ", err)
	// }

	for _, chat := range roomChats {
		values.RoomUsers[chat.RoomID] = chat.RegisteredUsers
	}

	result, err = db.Collection(values.UsersCollectionName).Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalln("error while getting all room names ", err)
	}

	var users []User
	err = result.Decode(&users)
	// if err != nil {
	// 	log.Fatalln("error converting room users interface ", err)
	// }

	for _, user := range users {
		values.Users[user.ID] = user.Name
	}
}
