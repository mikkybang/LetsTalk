package model

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/metaclips/LetsTalk/values"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

func InitDB() {
	dbHost := os.Getenv("db_host")
	mongoDB, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbHost))
	if err != nil {
		log.Fatalln(err)
	}

	db = mongoDB.Database(values.DatabaseName)

	// Ping mongo database continuously if down.
	go func(mongoDB *mongo.Client) {
		for {
			if err := mongoDB.Ping(context.TODO(), readpref.Primary()); err != nil {
				log.Fatalln(err)
			}

			time.Sleep(time.Second * 5)
		}
	}(mongoDB)

	values.MapEmailToName = make(map[string]string)

	getCollection := func(collection string, content interface{}) {
		result, err := db.Collection(collection).Find(context.TODO(), bson.D{})
		if err != nil {
			log.Fatalln("error while getting collection", err)
		}

		err = result.All(context.TODO(), content)
		if err != nil {
			log.Fatalln("error getting collection results", err)
		}
	}

	createNewAdminIfNotExist()

	var users []User
	getCollection(values.UsersCollectionName, &users)

	for _, user := range users {
		values.MapEmailToName[user.Email] = user.Name
	}
}

func createNewAdminIfNotExist() {
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), values.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}

	admin := Admin{
		StaffDetails: User{
			Email:    "admin@email.com",
			Name:     "admin",
			Password: password,
		},
		Super: true,
	}

	admin.CreateAdmin()
}
