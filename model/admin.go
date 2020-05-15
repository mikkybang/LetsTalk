package model

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/metaclips/LetsTalk/values"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *Admin) CheckAdminDetails(password string) error {
	result := db.Collection(values.AdminCollectionName).FindOne(context.TODO(), bson.M{
		"_id": b.StaffDetails.Email,
	})

	if err := result.Decode(&b); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(b.StaffDetails.Password, []byte(password)); err != nil {
		return err
	}
	return nil
}

func UploadUser(user User, r *http.Request) error {
	user.Name = strings.Title(user.Name)
	if names := strings.Split(user.Name, " "); len(names) > 1 {
		var err error
		user.Password, err = bcrypt.GenerateFromPassword([]byte(names[0]), defaultCost)
		if err != nil {
			return err
		}
	}
	id := strings.Split(user.Email, "@")
	if len(id) > 1 {
		user.ID = id[0]
	}

	if user.Class == "student" {
		user.ParentEmail = r.FormValue("parentEmail")
		user.ParentNumber = r.FormValue("parentNumber")
	}

	values.Users[user.Email] = user.Name
	_, err := db.Collection(values.UsersCollectionName).InsertOne(context.TODO(), user)
	return err
}
