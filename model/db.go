package model

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/metaclips/FinalYearProject/values"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	db   *mongo.Database
	UUID uuid.UUID
)

func InitDB() {
	mongoDB, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
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

}
