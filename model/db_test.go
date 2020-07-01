package model

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConfig(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Error("Error loading .env file")
	}
	if os.Getenv("db_host") == "" {
		t.Error("Environment variable not set")
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
