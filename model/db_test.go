package model

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConfig(t *testing.T) {
	err := LoadConfiguration()
	if err != nil {
		t.Error("Could not Load Config file", err)
	}
}

func TestDatabase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("db_host")))
	if err != nil {
		t.Errorf(err.Error())
	}
}
