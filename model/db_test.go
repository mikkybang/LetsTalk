package model

import (
	"context"
	"testing"
	"time"

	"github.com/metaclips/LetsTalk/values"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDatabase(t *testing.T) {
	if err := values.LoadConfiguration("../config.json"); err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mongo.Connect(ctx, options.Client().ApplyURI(values.Config.DbHost))
	if err != nil {
		t.Error(err)
	}
}
