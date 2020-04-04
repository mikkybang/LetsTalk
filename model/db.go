package model

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Client

func InitDB() {
	var err error

	db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatalln(err)
	}

	// Ping mongo database if up
	go func() {
		for {
			if err := db.Ping(context.TODO(), readpref.Primary()); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
